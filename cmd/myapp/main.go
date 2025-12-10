package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"stock-bot/domain/model"
	"stock-bot/domain/repository"
	"stock-bot/internal/app"
	"stock-bot/internal/config"
	"stock-bot/internal/handler/web"
	"stock-bot/internal/infrastructure/client"
	"sync"
	"syscall"
	"time"

	_ "stock-bot/internal/logger" // loggerパッケージをインポートし、slog.Default()を初期化

	ordersvr "stock-bot/gen/http/order/server"
	order "stock-bot/gen/order"

	goahttp "goa.design/goa/v3/http"
	"goa.design/goa/v3/http/middleware"

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

// dummyOrderRepo は repository.OrderRepository のダミー実装です。
// 将来的には infrastructure/repository の実装に置き換えます。
type dummyOrderRepo struct{}

func (r *dummyOrderRepo) Save(ctx context.Context, order *model.Order) error {
	slog.Default().Info("dummyOrderRepo: Save called", slog.Any("order", order))
	return nil // 常に成功
}
func (r *dummyOrderRepo) FindByID(ctx context.Context, orderID string) (*model.Order, error) {
	return nil, nil
}
func (r *dummyOrderRepo) FindByStatus(ctx context.Context, status model.OrderStatus) ([]*model.Order, error) {
	return nil, nil
}
func NewDummyOrderRepository() repository.OrderRepository {
	return &dummyOrderRepo{}
}

func main() {
	// 1. ロガーのセットアップ
	goaLogger := &GoaSlogger{slog.Default()}

	// 2. 設定ファイルの読み込み
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		slog.Default().Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	// 3. Usecaseなどの依存関係を初期化
	// 3-1. 証券会社APIクライアントを初期化
	// NewTachibanaClient は OrderClient インターフェースなどを満たす実装を返す
	tachibanaClient := client.NewTachibanaClient(cfg)

	slog.Default().Info("logging in to Tachibana API...")
	loginReq := request_auth.ReqLogin{
		UserId:   cfg.TachibanaUserID,
		Password: cfg.TachibanaPassword,
	}
	_, err = tachibanaClient.Login(context.Background(), loginReq)
	if err != nil {
		slog.Default().Error("failed to login", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Default().Info("login successful")

	// 3-2. リポジトリを初期化 (今回はダミー)
	orderRepo := NewDummyOrderRepository()

	// 3-3. ユースケースを初期化
	// tachibanaClient は OrderClient インターフェースを満たしているので直接渡せる
	orderUsecase := app.NewOrderUseCaseImpl(tachibanaClient, orderRepo, cfg.TachibanaPassword)

	// 4. Goaサービスの実装を初期化
	orderSvc := web.NewOrderService(orderUsecase, slog.Default())

	// 5. GoaのエンドポイントとHTTPハンドラを構築
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	orderEndpoints := order.NewEndpoints(orderSvc)

	mux := goahttp.NewMuxer()

	server := ordersvr.New(orderEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)

	ordersvr.Mount(mux, server)

	// OpenAPI仕様配信用に静的ファイルサーバーをマウント
	fs := http.FileServer(http.Dir("./gen/http/openapi"))
	mux.Handle("GET", "/swagger/", http.HandlerFunc(http.StripPrefix("/swagger/", fs).ServeHTTP))

	// 6. HTTPサーバーの起動
	addr := fmt.Sprintf("http://localhost:%d", cfg.HTTPPort)
	u, err := url.Parse(addr)
	if err != nil {
		slog.Default().Error("invalid URL", slog.String("address", addr), slog.Any("error", err))
		os.Exit(1)
	}

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

	// 7. Graceful Shutdownの設定
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
	case sig := <-c:
		slog.Default().Info(fmt.Sprintf("received signal %s, shutting down", sig))
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Default().Error("failed to shutdown server", slog.Any("error", err))
	}

	wg.Wait()
	slog.Default().Info("shutdown complete")
}
