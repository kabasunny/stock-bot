package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"stock-bot/domain/model"
	"stock-bot/internal/app"
	"stock-bot/internal/config"
	infra "stock-bot/internal/infrastructure/client"
	repo "stock-bot/internal/infrastructure/repository"
	_ "stock-bot/internal/logger" // カスタムロガーのinit()を呼び出すために必要

	// Goa-generated files
	balance "stock-bot/gen/balance"
	order "stock-bot/gen/order"
	position "stock-bot/gen/position"

	balancesvr "stock-bot/gen/http/balance/server"
	ordersvr "stock-bot/gen/http/order/server"
	positionsvr "stock-bot/gen/http/position/server"
	goahttp "goa.design/goa/v3/http"
	goahttpmiddleware "goa.design/goa/v3/http/middleware"

	// Service implementations
	balanceservice "stock-bot/internal/interface/balance"
	orderservice "stock-bot/internal/interface/order"
	positionservice "stock-bot/internal/interface/position"
)

func main() {
	// 1. 設定の読み込み
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// 2. データベースへの接続
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Tokyo",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		os.Exit(1)
	}

	// 3. データベースのマイグレーション
	if err := db.AutoMigrate(&model.Order{}, &model.Execution{}, &model.Position{}); err != nil {
		slog.Error("failed to auto migrate database", "error", err)
		os.Exit(1)
	}
	slog.Info("database migration completed")

	// 4. 依存関係の解決 (DI)
	// API Client
	tachibanaClient := infra.NewTachibanaClient(cfg)
	if err := tachibanaClient.Login(context.Background()); err != nil {
		slog.Error("failed to login tachibana client", "error", err)
		os.Exit(1)
	}
	slog.Info("tachibana client logged in")

	// Repository
	orderRepo := repo.NewOrderRepository(db)
	positionRepo := repo.NewPositionRepository(db)
	masterRepo := repo.NewMasterRepository(db)

	// UseCase
	orderUC := app.NewOrderUseCase(orderRepo)
	positionUC := app.NewPositionUseCase(positionRepo)
	balanceUC := app.NewBalanceUseCase(tachibanaClient)

	// Service
	balanceSvc := balanceservice.NewBalanceService(balanceUC)
	orderSvc := orderservice.NewOrderService(orderUC)
	positionSvc := positionservice.NewPositionService(positionUC)

	// 5. Goaエンドポイントの作成
	balanceEndpoints := balance.NewEndpoints(balanceSvc)
	orderEndpoints := order.NewEndpoints(orderSvc)
	positionEndpoints := position.NewEndpoints(positionSvc)

	// 6. HTTPサーバーのセットアップ
	mux := goahttp.NewMuxer()

	balanceServer := balancesvr.New(balanceEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
	orderServer := ordersvr.New(orderEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)
	positionServer := positionsvr.New(positionEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, nil, nil)

	balancesvr.Mount(mux, balanceServer)
	ordersvr.Mount(mux, orderServer)
	positionsvr.Mount(mux, positionServer)

	var handler http.Handler = mux
	handler = goahttpmiddleware.Log(slog.Default())(handler)
	handler = goahttpmiddleware.RequestID()(handler)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: handler,
	}

	// 7. サーバーの起動とGraceful Shutdown
	errCh := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errCh <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		slog.Info("HTTP server starting", "port", cfg.HTTPPort)
		errCh <- httpServer.ListenAndServe()
	}()

	slog.Warn("server exiting", "reason", <-errCh)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	httpServer.Shutdown(ctx)
	slog.Info("server gracefully shut down")
}