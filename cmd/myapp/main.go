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
	"stock-bot/internal/config"
	"stock-bot/internal/handler/web"
	"stock-bot/internal/infrastructure/container"
	"sync"
	"syscall"
	"time"

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

	// 2. ロガーのセットアップ
	goaLogger := &GoaSlogger{slog.Default()}

	// 3. 設定ファイルの読み込み
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

	// 4. DIコンテナの初期化
	containerCfg := &container.ContainerConfig{
		SkipSync:    *skipSync,
		NoDB:        *noDB,
		NoTachibana: *noTachibana,
	}

	diContainer, err := container.NewContainer(cfg, slog.Default(), containerCfg)
	if err != nil {
		slog.Default().Error("failed to initialize DI container", slog.Any("error", err))
		os.Exit(1)
	}
	defer diContainer.Shutdown()

	// 5. Goaサービスの実装を初期化
	var balanceSvc balance.Service
	var positionSvc positiongen.Service
	var priceSvc pricegen.Service
	var tradeSvc tradegen.Service
	var orderSvc order.Service
	var masterSvc mastergen.Service

	if !*noTachibana {
		balanceSvc = web.NewBalanceService(diContainer.GetBalanceUseCase(), diContainer.GetLogger(), diContainer.GetAppSession())
		positionSvc = web.NewPositionService(diContainer.GetPositionUseCase(), diContainer.GetLogger(), diContainer.GetAppSession())
		priceSvc = web.NewPriceService(diContainer.GetPriceUseCase(), diContainer.GetLogger(), diContainer.GetAppSession())
		tradeSvc = web.NewTradeService(diContainer.GetTradeService(), diContainer.GetLogger(), diContainer.GetAppSession())
	}
	if !*noDB && !*noTachibana {
		orderSvc = web.NewOrderService(diContainer.GetOrderUseCase(), diContainer.GetLogger(), diContainer.GetAppSession())
		masterSvc = web.NewMasterService(diContainer.GetMasterUseCase(), diContainer.GetLogger(), diContainer.GetAppSession())
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
		if !*noDB && !*noTachibana {
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

	// 8. HTTPサーバーの起動
	srv := &http.Server{
		Addr:    u.Host,
		Handler: diContainer.GetErrorHandler().Middleware(middleware.Log(goaLogger)(mux)),
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

	// 9. Graceful Shutdownの設定
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
