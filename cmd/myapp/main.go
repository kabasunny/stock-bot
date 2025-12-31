package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"stock-bot/internal/app"
	"stock-bot/internal/config"
	"stock-bot/internal/handler/web"
	"stock-bot/internal/infrastructure/client"
	repository_impl "stock-bot/internal/infrastructure/repository"
	"stock-bot/internal/tradeservice"
	"sync"
	"syscall"
	"time"

	"stock-bot/domain/repository"
	_ "stock-bot/internal/logger" // loggerパッケージをインポートし、slog.Default()を初期化

	balance "stock-bot/gen/balance"
	balancesvr "stock-bot/gen/http/balance/server"
	mastersvr "stock-bot/gen/http/master/server"
	ordersvr "stock-bot/gen/http/order/server"
	positionsvr "stock-bot/gen/http/position/server"
	pricesvr "stock-bot/gen/http/price/server"
	tradesvr "stock-bot/gen/http/trade/server"
	mastergen "stock-bot/gen/master"
	order "stock-bot/gen/order"
	positiongen "stock-bot/gen/position"
	pricegen "stock-bot/gen/price"
	tradegen "stock-bot/gen/trade"

	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/http/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GoaSlogger は *slog.Logger を goa.design/goa/v3/http/middleware.Logger インターフェースに適合させるためのラッパーです。
type GoaSlogger struct {
	logger *slog.Logger
}

// Log は middleware.Logger インターフェースの Log メソッドを実装します。
func (l *GoaSlogger) Log(keyvals ...interface{}) error {
	l.logger.Info("HTTP Request", keyvals...)
	return nil
}

func main() {
	// 1. コマンドラインフラグの定義
	skipSync := flag.Bool("skip-sync", false, "Skip initial master data synchronization on startup")
	noDB := flag.Bool("no-db", false, "Disable database connection and related features")
	noTachibana := flag.Bool("no-tachibana", false, "Disable Tachibana API client initialization and login")
	flag.Parse()

	// 1. ロガーのセットアップ
	goaLogger := &GoaSlogger{slog.Default()}

	// 2. 設定ファイルの読み込み
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		if !*noTachibana { // Only error out if Tachibana is needed
			slog.Default().Error("failed to load config", slog.Any("error", err))
			os.Exit(1)
		} else {
			slog.Default().Warn("failed to load .env, but --no-tachibana is set. Using default http port.", slog.Any("error", err))
			// Initialize cfg with default values to prevent nil pointer dereference
			cfg = &config.Config{
				HTTPPort: 8080, // Default value from config.go
			}
		}
	}

	// 3. データベース接続と依存コンポーネントの初期化
	var db *gorm.DB
	var orderRepo repository.OrderRepository
	var masterRepo repository.MasterRepository
	var orderUsecase app.OrderUseCase
	var masterUsecase app.MasterUseCase
	var orderSvc order.Service
	var masterSvc mastergen.Service
	var goaTradeService *tradeservice.GoaTradeService // Declare executionUseCase

	if !*noDB {
		// 3a. データベース接続
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
			cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			slog.Default().Error("failed to connect database", slog.Any("error", err))
			os.Exit(1)
		}
		slog.Default().Info("database connection established")

		// 3b. DB依存リポジトリを初期化
		orderRepo = repository_impl.NewOrderRepository(db)
		masterRepo = repository_impl.NewMasterRepository(db)
		// positionRepo は不要（エージェント専用）
	} else {
		slog.Default().Warn("database connection is disabled due to --no-db flag")
	}

	// 4. Usecaseなどの依存関係を初期化
	var tachibanaClient *client.TachibanaClientImpl
	var unifiedClient *client.TachibanaUnifiedClient
	var unifiedClientAdapter *client.TachibanaUnifiedClientAdapter
	var appSession *client.Session

	// 4-Z. イベントクライアントの初期化は不要（エージェント専用）

	if !*noTachibana {
		// 4-1. 証券会社APIクライアントを初期化
		tachibanaClient = client.NewTachibanaClient(cfg)

		// 4-1b. 統合クライアントを初期化
		unifiedClient = client.NewTachibanaUnifiedClient(
			tachibanaClient, // AuthClient
			tachibanaClient, // BalanceClient
			tachibanaClient, // OrderClient
			tachibanaClient, // PriceInfoClient
			tachibanaClient, // MasterDataClient
			nil,             // EventClient (エージェント専用のため不要)
			cfg.TachibanaUserID,
			cfg.TachibanaPassword,
			cfg.TachibanaSecondPassword,
			slog.Default(),
		)

		// 4-1c. 統合クライアントアダプターを初期化
		unifiedClientAdapter = client.NewTachibanaUnifiedClientAdapter(unifiedClient)

		slog.Default().Info("logging in to Tachibana API via unified client...")

		// 統合クライアント経由でログイン（自動認証）
		appSession, err = unifiedClient.GetSession(context.Background())
		if err != nil {
			slog.Default().Error("failed to login via unified client", slog.Any("error", err))
			os.Exit(1)
		}

		slog.Default().Info("login successful via unified client")
	} else {
		slog.Default().Warn("Tachibana API client and login are disabled due to --no-tachibana flag.")
	}

	// 4-2. ユースケースを初期化 (DB依存/非依存)
	var balanceUsecase app.BalanceUseCase
	var positionUsecase app.PositionUseCase
	var priceUsecase app.PriceUseCase

	if !*noTachibana {
		balanceUsecase = app.NewBalanceUseCaseImpl(tachibanaClient)
		positionUsecase = app.NewPositionUseCaseImpl(tachibanaClient)
		priceUsecase = app.NewPriceUseCaseImpl(tachibanaClient, appSession)
	}

	if !*noDB && !*noTachibana {
		orderUsecase = app.NewOrderUseCaseImpl(tachibanaClient, orderRepo)
		masterUsecase = app.NewMasterUseCaseImpl(tachibanaClient, masterRepo)
		// executionUseCase は不要（エージェント専用） // Initialize executionUseCase

		if !*skipSync {
			slog.Default().Info("Starting initial master data synchronization...")
			err = masterUsecase.DownloadAndStoreMasterData(context.Background(), appSession)
			if err != nil {
				slog.Default().Error("failed to download and store master data on startup", slog.Any("error", err))
				os.Exit(1)
			}
			slog.Default().Info("Initial master data synchronization completed successfully.")
		} else {
			slog.Default().Info("Skipping initial master data synchronization.")
		}

		// 4-Y. エージェント用トレードサービスの初期化
		goaTradeService = tradeservice.NewGoaTradeService(
			unifiedClientAdapter, // TachibanaUnifiedClientAdapter implements BalanceClient
			unifiedClientAdapter, // TachibanaUnifiedClientAdapter implements OrderClient
			unifiedClientAdapter, // TachibanaUnifiedClientAdapter implements PriceInfoClient
			orderRepo,
			appSession,
			slog.Default(),
		)
	}

	// 5. Goaサービスの実装を初期化
	var balanceSvc balance.Service
	var positionSvc positiongen.Service
	var priceSvc pricegen.Service
	var tradeSvc tradegen.Service

	if !*noTachibana {
		balanceSvc = web.NewBalanceService(balanceUsecase, slog.Default(), appSession)
		positionSvc = web.NewPositionService(positionUsecase, slog.Default(), appSession)
		priceSvc = web.NewPriceService(priceUsecase, slog.Default(), appSession)
		tradeSvc = web.NewTradeService(goaTradeService, slog.Default(), appSession)
	}
	if !*noDB && !*noTachibana {
		orderSvc = web.NewOrderService(orderUsecase, slog.Default(), appSession)
		masterSvc = web.NewMasterService(masterUsecase, slog.Default(), appSession)
	}

	// 6. GoaのエンドポイントとHTTPハンドラを構築
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	mux := goahttp.NewMuxer()

	if !*noTachibana {
		balanceEndpoints := balance.NewEndpoints(balanceSvc)
		positionEndpoints := positiongen.NewEndpoints(positionSvc)
		priceEndpoints := pricegen.NewEndpoints(priceSvc)
		tradeEndpoints := tradegen.NewEndpoints(tradeSvc)

		balancesvr.Mount(mux, balancesvr.New(balanceEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil))
		positionsvr.Mount(mux, positionsvr.New(positionEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil))
		pricesvr.Mount(mux, pricesvr.New(priceEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil))
		tradesvr.Mount(mux, tradesvr.New(tradeEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil))

		// DB依存かつTachibana API依存のエンドポイント
		if !*noDB {
			orderEndpoints := order.NewEndpoints(orderSvc)
			masterEndpoints := mastergen.NewEndpoints(masterSvc)
			ordersvr.Mount(mux, ordersvr.New(orderEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil))
			mastersvr.Mount(mux, mastersvr.New(masterEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil))
		}
	} else {
		slog.Default().Warn("Skipping mounting of all API endpoints due to --no-tachibana flag.")
	}

	fs := http.FileServer(http.Dir("./gen/http/openapi"))
	mux.Handle("GET", "/swagger/", http.HandlerFunc(http.StripPrefix("/swagger/", fs).ServeHTTP))

	// 7. HTTPサーバーとエージェントの起動
	addr := fmt.Sprintf("http://localhost:%d", cfg.HTTPPort)
	u, err := url.Parse(addr)
	if err != nil {
		slog.Default().Error("invalid URL", slog.String("address", addr), slog.Any("error", err))
		os.Exit(1)
	}

	// 7. HTTPサーバーの起動
	srv := &http.Server{
		Addr:    u.Host,
		Handler: middleware.Log(goaLogger)(mux),
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		slog.Default().Info(fmt.Sprintf("HTTP server listening on %q", u.Host))
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			slog.Default().Error("server error", slog.Any("error", err))
			cancel()
		}
	}()

	// 8. Graceful Shutdownの設定
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case sig := <-c:
		slog.Default().Info(fmt.Sprintf("received signal %s, shutting down", sig))
	}

	// サーバーを停止
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Default().Error("failed to shutdown server", slog.Any("error", err))
	}

	wg.Wait()
	slog.Default().Info("shutdown complete")
}
