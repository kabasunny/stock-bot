package container

import (
	"context"
	"fmt"
	"log/slog"
	"stock-bot/domain/event"
	"stock-bot/domain/repository"
	"stock-bot/domain/service"
	"stock-bot/internal/app"
	"stock-bot/internal/config"
	"stock-bot/internal/eventprocessing"
	"stock-bot/internal/infrastructure/client"
	"stock-bot/internal/infrastructure/errors"
	repository_impl "stock-bot/internal/infrastructure/repository"
	"stock-bot/internal/infrastructure/scheduler"
	"stock-bot/internal/tradeservice"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Container はDIコンテナ
type Container struct {
	config *config.Config
	logger *slog.Logger
	db     *gorm.DB

	// Infrastructure Layer
	tachibanaClient      *client.TachibanaClientImpl
	unifiedClient        *client.TachibanaUnifiedClient
	unifiedClientAdapter *client.TachibanaUnifiedClientAdapter
	eventClient          client.EventClient
	appSession           *client.Session

	// Repository Layer
	orderRepo    repository.OrderRepository
	masterRepo   repository.MasterRepository
	positionRepo repository.PositionRepository
	strategyRepo repository.StrategyRepository

	// Use Case Layer
	balanceUseCase   app.BalanceUseCase
	orderUseCase     app.OrderUseCase
	masterUseCase    app.MasterUseCase
	positionUseCase  app.PositionUseCase
	priceUseCase     app.PriceUseCase
	executionUseCase app.ExecutionUseCase
	strategyUseCase  app.StrategyUseCase

	// Domain Service Layer
	tradeService    service.TradeService
	strategyService service.StrategyService

	// Domain Event System
	eventPublisher event.EventPublisher
	unitOfWork     repository.UnitOfWork

	// Infrastructure Services
	masterScheduler     *scheduler.MasterDataScheduler
	orderEventProcessor *tradeservice.OrderEventProcessor
	errorHandler        *errors.HTTPErrorHandler
}

// ContainerConfig はコンテナの設定
type ContainerConfig struct {
	SkipSync    bool
	NoDB        bool
	NoTachibana bool
}

// NewContainer は新しいDIコンテナを作成する
func NewContainer(cfg *config.Config, logger *slog.Logger, containerCfg *ContainerConfig) (*Container, error) {
	container := &Container{
		config: cfg,
		logger: logger,
	}

	// エラーハンドラーの初期化
	container.errorHandler = errors.NewHTTPErrorHandler(logger)

	// データベース接続の初期化
	if !containerCfg.NoDB {
		if err := container.initDatabase(); err != nil {
			return nil, fmt.Errorf("failed to initialize database: %w", err)
		}
		if err := container.initRepositories(); err != nil {
			return nil, fmt.Errorf("failed to initialize repositories: %w", err)
		}
		if err := container.initDomainEventSystem(); err != nil {
			return nil, fmt.Errorf("failed to initialize domain event system: %w", err)
		}
	}

	// Tachibanaクライアントの初期化
	if !containerCfg.NoTachibana {
		if err := container.initTachibanaClients(); err != nil {
			return nil, fmt.Errorf("failed to initialize Tachibana clients: %w", err)
		}
		if err := container.initUseCases(); err != nil {
			return nil, fmt.Errorf("failed to initialize use cases: %w", err)
		}
		if err := container.initDomainServices(); err != nil {
			return nil, fmt.Errorf("failed to initialize domain services: %w", err)
		}
	}

	// インフラサービスの初期化
	if !containerCfg.NoDB && !containerCfg.NoTachibana {
		if err := container.initInfraServices(containerCfg.SkipSync); err != nil {
			return nil, fmt.Errorf("failed to initialize infrastructure services: %w", err)
		}
	}

	// エラーハンドラーの初期化
	container.errorHandler = errors.NewHTTPErrorHandler(logger)

	return container, nil
}

// initDatabase はデータベース接続を初期化する
func (c *Container) initDatabase() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		c.config.DBHost, c.config.DBUser, c.config.DBPassword, c.config.DBName, c.config.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	c.db = db
	c.logger.Info("database connection established")
	return nil
}

// initRepositories はリポジトリ層を初期化する
func (c *Container) initRepositories() error {
	c.orderRepo = repository_impl.NewOrderRepository(c.db)
	c.masterRepo = repository_impl.NewMasterRepository(c.db)
	c.positionRepo = repository_impl.NewPositionRepository(c.db, c.orderRepo)
	c.strategyRepo = repository_impl.NewStrategyRepository(c.db)
	return nil
}

// initTachibanaClients はTachibanaクライアントを初期化する
func (c *Container) initTachibanaClients() error {
	// 基本クライアントの初期化
	c.tachibanaClient = client.NewTachibanaClient(c.config)
	c.eventClient = client.NewEventClient(c.logger)

	// 統合クライアントの初期化
	c.unifiedClient = client.NewTachibanaUnifiedClient(
		c.tachibanaClient, // AuthClient
		c.tachibanaClient, // BalanceClient
		c.tachibanaClient, // OrderClient
		c.tachibanaClient, // PriceInfoClient
		c.tachibanaClient, // MasterDataClient
		c.eventClient,     // EventClient
		c.config.TachibanaUserID,
		c.config.TachibanaPassword,
		c.config.TachibanaSecondPassword,
		c.logger,
	)

	// 統合クライアントアダプターの初期化
	c.unifiedClientAdapter = client.NewTachibanaUnifiedClientAdapter(c.unifiedClient)

	// ログイン処理
	c.logger.Info("logging in to Tachibana API via unified client...")
	session, err := c.unifiedClient.GetSession(context.Background())
	if err != nil {
		return fmt.Errorf("failed to login via unified client: %w", err)
	}
	c.appSession = session
	c.logger.Info("login successful via unified client")

	return nil
}

// initUseCases はユースケース層を初期化する
func (c *Container) initUseCases() error {
	c.balanceUseCase = app.NewBalanceUseCaseImpl(c.tachibanaClient)
	c.positionUseCase = app.NewPositionUseCaseImpl(c.tachibanaClient)
	c.priceUseCase = app.NewPriceUseCaseImpl(c.tachibanaClient, c.appSession)

	if c.orderRepo != nil {
		c.orderUseCase = app.NewOrderUseCaseImpl(c.tachibanaClient, c.orderRepo)
	}
	if c.masterRepo != nil {
		c.masterUseCase = app.NewMasterUseCaseImpl(c.tachibanaClient, c.masterRepo)
	}
	if c.orderRepo != nil && c.positionRepo != nil {
		c.executionUseCase = app.NewExecutionUseCaseImpl(c.orderRepo, c.positionRepo)
	}
	if c.strategyRepo != nil {
		c.strategyUseCase = app.NewStrategyUseCaseImpl(c.strategyRepo, c.strategyService, c.unitOfWork, c.logger)
	}

	return nil
}

// ドメインサービス層を初期化する
func (c *Container) initDomainServices() error {
	if c.unifiedClientAdapter != nil && c.orderRepo != nil && c.masterRepo != nil {
		tradeService := tradeservice.NewGoaTradeService(
			c.unifiedClientAdapter, // BalanceClient
			c.unifiedClientAdapter, // OrderClient
			c.unifiedClientAdapter, // PriceInfoClient
			c.orderRepo,
			c.masterRepo,
			c.appSession,
			c.logger,
		)

		// セッション回復用の統合クライアントを設定
		tradeService.SetUnifiedClient(c.unifiedClient)
		c.tradeService = tradeService
	}
	if c.strategyRepo != nil && c.eventPublisher != nil {
		c.strategyService = tradeservice.NewStrategyService(c.strategyRepo, c.eventPublisher, c.logger)
	}
	return nil
}

// initInfraServices はインフラサービスを初期化する
func (c *Container) initInfraServices(skipSync bool) error {
	// マスターデータの初期同期
	if !skipSync && c.masterUseCase != nil {
		c.logger.Info("Starting initial master data synchronization...")
		err := c.masterUseCase.DownloadAndStoreMasterData(context.Background(), c.appSession)
		if err != nil {
			return fmt.Errorf("failed to download and store master data on startup: %w", err)
		}
		c.logger.Info("Initial master data synchronization completed successfully.")
	}

	// マスターデータスケジューラーの初期化
	if c.masterUseCase != nil {
		c.masterScheduler = scheduler.NewMasterDataScheduler(
			c.masterUseCase,
			c.appSession,
			c.logger,
		)
		c.masterScheduler.Start()
	}

	// 注文イベント処理器の初期化
	if c.orderRepo != nil && c.positionRepo != nil {
		c.orderEventProcessor = tradeservice.NewOrderEventProcessor(
			c.orderRepo,
			c.positionRepo,
			c.eventClient,
			c.appSession,
			c.logger,
		)

		if err := c.orderEventProcessor.Start(context.Background(), []string{}); err != nil {
			c.logger.Error("failed to start order event processor", slog.Any("error", err))
			// エラーでも続行（WebSocketは必須ではない）
		} else {
			c.logger.Info("order event processor started successfully")
		}
	}

	return nil
}

// Shutdown はコンテナのリソースをクリーンアップする
func (c *Container) Shutdown() {
	if c.masterScheduler != nil {
		c.masterScheduler.Stop()
	}
	if c.orderEventProcessor != nil {
		c.orderEventProcessor.Stop()
	}
	if c.db != nil {
		if sqlDB, err := c.db.DB(); err == nil {
			sqlDB.Close()
		}
	}
}

// Getters for accessing dependencies

func (c *Container) GetConfig() *config.Config {
	return c.config
}

func (c *Container) GetLogger() *slog.Logger {
	return c.logger
}

func (c *Container) GetDB() *gorm.DB {
	return c.db
}

func (c *Container) GetTachibanaClient() *client.TachibanaClientImpl {
	return c.tachibanaClient
}

func (c *Container) GetUnifiedClient() *client.TachibanaUnifiedClient {
	return c.unifiedClient
}

func (c *Container) GetUnifiedClientAdapter() *client.TachibanaUnifiedClientAdapter {
	return c.unifiedClientAdapter
}

func (c *Container) GetEventClient() client.EventClient {
	return c.eventClient
}

func (c *Container) GetAppSession() *client.Session {
	return c.appSession
}

func (c *Container) GetOrderRepo() repository.OrderRepository {
	return c.orderRepo
}

func (c *Container) GetMasterRepo() repository.MasterRepository {
	return c.masterRepo
}

func (c *Container) GetPositionRepo() repository.PositionRepository {
	return c.positionRepo
}

func (c *Container) GetBalanceUseCase() app.BalanceUseCase {
	return c.balanceUseCase
}

func (c *Container) GetOrderUseCase() app.OrderUseCase {
	return c.orderUseCase
}

func (c *Container) GetMasterUseCase() app.MasterUseCase {
	return c.masterUseCase
}

func (c *Container) GetPositionUseCase() app.PositionUseCase {
	return c.positionUseCase
}

func (c *Container) GetPriceUseCase() app.PriceUseCase {
	return c.priceUseCase
}

func (c *Container) GetExecutionUseCase() app.ExecutionUseCase {
	return c.executionUseCase
}

func (c *Container) GetTradeService() service.TradeService {
	return c.tradeService
}

func (c *Container) GetMasterScheduler() *scheduler.MasterDataScheduler {
	return c.masterScheduler
}

func (c *Container) GetOrderEventProcessor() *tradeservice.OrderEventProcessor {
	return c.orderEventProcessor
}
func (c *Container) GetErrorHandler() *errors.HTTPErrorHandler {
	return c.errorHandler
}

// initDomainEventSystem はドメインイベントシステムを初期化する
func (c *Container) initDomainEventSystem() error {
	// イベントパブリッシャーの初期化
	c.eventPublisher = event.NewInMemoryEventPublisher(c.logger)

	// UnitOfWorkの初期化
	c.unitOfWork = repository_impl.NewUnitOfWork(c.db, c.eventPublisher, c.logger)

	// ドメインイベントハンドラーの登録
	orderEventHandler := eventprocessing.NewOrderEventHandler(c.orderRepo, c.logger)
	positionEventHandler := eventprocessing.NewPositionEventHandler(c.positionRepo, c.logger)
	riskEventHandler := eventprocessing.NewRiskEventHandler(c.logger)

	// イベントハンドラーを登録
	c.eventPublisher.Subscribe("OrderPlaced", orderEventHandler)
	c.eventPublisher.Subscribe("OrderExecuted", orderEventHandler)
	c.eventPublisher.Subscribe("OrderCancelled", orderEventHandler)
	c.eventPublisher.Subscribe("PositionOpened", positionEventHandler)
	c.eventPublisher.Subscribe("PositionClosed", positionEventHandler)
	c.eventPublisher.Subscribe("RiskLimitExceeded", riskEventHandler)

	c.logger.Info("domain event system initialized",
		slog.Int("registered_event_types", len(c.eventPublisher.GetAllEventTypes())))

	return nil
}

func (c *Container) GetEventPublisher() event.EventPublisher {
	return c.eventPublisher
}

func (c *Container) GetUnitOfWork() repository.UnitOfWork {
	return c.unitOfWork
}
func (c *Container) GetStrategyRepo() repository.StrategyRepository {
	return c.strategyRepo
}

func (c *Container) GetStrategyUseCase() app.StrategyUseCase {
	return c.strategyUseCase
}

func (c *Container) GetStrategyService() service.StrategyService {
	return c.strategyService
}
