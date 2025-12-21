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
	"stock-bot/internal/agent"
	"stock-bot/internal/app"
	"stock-bot/internal/config"
	"stock-bot/internal/handler/web"
	"stock-bot/internal/infrastructure/client"
	repository_impl "stock-bot/internal/infrastructure/repository"
	"sync"
	"syscall"
	"time"

	_ "stock-bot/internal/logger" // loggerパッケージをインポートし、slog.Default()を初期化

	balance "stock-bot/gen/balance"
	balancesvr "stock-bot/gen/http/balance/server"
	mastersvr "stock-bot/gen/http/master/server" // New import
	ordersvr "stock-bot/gen/http/order/server"
	positionsvr "stock-bot/gen/http/position/server" // New import
	mastergen "stock-bot/gen/master"                 // New import
	order "stock-bot/gen/order"
	positiongen "stock-bot/gen/position" // New import

	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/http/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	request_auth "stock-bot/internal/infrastructure/client/dto/auth/request"
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
	flag.Parse()

	// 1. ロガーのセットアップ
	goaLogger := &GoaSlogger{slog.Default()}

	// 2. 設定ファイルの読み込み
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		slog.Default().Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	// 3. データベース接続
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Default().Error("failed to connect database", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Default().Info("database connection established")

	// 4. Usecaseなどの依存関係を初期化
	// 4-1. 証券会社APIクライアントを初期化
	tachibanaClient := client.NewTachibanaClient(cfg)

	slog.Default().Info("logging in to Tachibana API...")
	loginReq := request_auth.ReqLogin{
		UserId:   cfg.TachibanaUserID,
		Password: cfg.TachibanaPassword,
	}
	appSession, err := tachibanaClient.LoginWithPost(context.Background(), loginReq)
	if err != nil {
		slog.Default().Error("failed to login", slog.Any("error", err))
		os.Exit(1)
	}
	if appSession.ResultCode != "0" {
		slog.Default().Error("failed to login: API returned error", slog.String("code", appSession.ResultCode), slog.String("text", appSession.ResultText))
		os.Exit(1)
	}
	slog.Default().Info("login successful")

	// 4-2. リポジトリを初期化
	orderRepo := repository_impl.NewOrderRepository(db)
	masterRepo := repository_impl.NewMasterRepository(db)

	// 4-3. ユースケースを初期化
	orderUsecase := app.NewOrderUseCaseImpl(tachibanaClient, orderRepo)
	balanceUsecase := app.NewBalanceUseCaseImpl(tachibanaClient)
	positionUsecase := app.NewPositionUseCaseImpl(tachibanaClient)
	masterUsecase := app.NewMasterUseCaseImpl(tachibanaClient, masterRepo)

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
	goaTradeService := agent.NewGoaTradeService(
		tachibanaClient, // tachibanaClient は BalanceClient インターフェースを実装
		tachibanaClient, // tachibanaClient は OrderClient インターフェースを実装
		appSession,
		slog.Default(),
	)

	// 5. Goaサービスの実装を初期化
	orderSvc := web.NewOrderService(orderUsecase, slog.Default(), appSession)
	balanceSvc := web.NewBalanceService(balanceUsecase, slog.Default(), appSession)
	positionSvc := web.NewPositionService(positionUsecase, slog.Default(), appSession)
	masterSvc := web.NewMasterService(masterUsecase, slog.Default(), appSession)

	// 6. GoaのエンドポイントとHTTPハンドラを構築
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	orderEndpoints := order.NewEndpoints(orderSvc)
	balanceEndpoints := balance.NewEndpoints(balanceSvc)
	positionEndpoints := positiongen.NewEndpoints(positionSvc)
	masterEndpoints := mastergen.NewEndpoints(masterSvc)

	mux := goahttp.NewMuxer()

	server := ordersvr.New(orderEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
	balanceserver := balancesvr.New(balanceEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
	positionserver := positionsvr.New(positionEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
	masterserver := mastersvr.New(masterEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)

	ordersvr.Mount(mux, server)
	balancesvr.Mount(mux, balanceserver)
	positionsvr.Mount(mux, positionserver)
	mastersvr.Mount(mux, masterserver)

	fs := http.FileServer(http.Dir("./gen/http/openapi"))
	mux.Handle("GET", "/swagger/", http.HandlerFunc(http.StripPrefix("/swagger/", fs).ServeHTTP))

	// 7. HTTPサーバーとエージェントの起動
	addr := fmt.Sprintf("http://localhost:%d", cfg.HTTPPort)
	u, err := url.Parse(addr)
	if err != nil {
		slog.Default().Error("invalid URL", slog.String("address", addr), slog.Any("error", err))
		os.Exit(1)
	}

	// 7-1. エージェントの初期化と起動
	agentConfigPath := "agent_config.yaml" // TODO: コマンドライン引数で渡せるようにする
	stockAgent, err := agent.NewAgent(agentConfigPath, goaTradeService)
	if err != nil {
		slog.Default().Error("failed to create agent", "config", agentConfigPath, slog.Any("error", err))
		os.Exit(1)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		stockAgent.Start()
	}()

	// 7-2. HTTPサーバーの起動
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

	// エージェントを停止
	stockAgent.Stop()

	// サーバーを停止
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Default().Error("failed to shutdown server", slog.Any("error", err))
	}

	wg.Wait()
	slog.Default().Info("shutdown complete")
}
