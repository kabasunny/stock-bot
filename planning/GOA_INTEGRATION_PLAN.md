# Goaフレームワーク導入手順書 (改訂版)

## 1. 目的
本ドキュメントは、`stock-bot`プロジェクトにGoaフレームワークを導入し、**Go内部のエージェント**がGo製のAPIラッパーを呼び出すための、堅牢なHTTP APIを構築する手順を定めることを目的とする。このAPI中心のアプローチは、将来的にエージェントをRust製マイクロサービスとして分離・移行する際の基盤となる。

## 2. 前提条件
- Go言語の開発環境がセットアップ済みであること。
- プロジェクトのルートディレクトリ (`C:\Users\kabas\project\stock-bot`) で作業を行うこと。

## 3. 導入手順

### ステップ1: `goa`ツールのインストール
まず、Goaのコードジェネレータをインストールします。これにより、ターミナルで`goa`コマンドが使えるようになります。

```shell
go install goa.design/goa/v3/cmd/goa@v3
```

### ステップ2: API設計ファイルの作成
次に、APIの設計を記述するためのファイルを作成します。

1.  プロジェクトルートに`design`ディレクトリを作成します。
2.  `design/design.go`というファイルを新規作成し、以下の内容を貼り付けます。

    これは、**Go内部のエージェント**からの注文リクエストを受け付けるための`/order`エンドポイントの設計です。`domain/model/order.go`の構造に合わせてpayloadを定義しています。

    ```go
    package design

    import (
        . "goa.design/goa/v3/dsl"
    )

    // API全体の定義
    var _ = API("stockbot", func() {
        Title("Stock Bot Service")
        Description("Service for placing and managing stock orders")
        Server("stockbot", func() {
            Host("localhost", func() {
                // ポートは後ほど設定ファイルから読み込む
                URI("http://localhost:8080")
            })
        })
    })

    // 注文サービス(Order)の定義
    var _ = Service("order", func() {
        Description("The order service handles placing stock orders.")

        // POST /order
        Method("create", func() {
            Description("Create a new stock order.")

            // リクエストのペイロード(JSONボディ)
            Payload(func() {
                Attribute("symbol", String, "銘柄コード (例: 7203)")
                Attribute("trade_type", String, "売買区分 (BUY/SELL)", func() {
                    Enum("BUY", "SELL")
                })
                Attribute("order_type", String, "注文種別 (MARKET/LIMITなど)", func() {
                    Enum("MARKET", "LIMIT", "STOP", "STOP_LIMIT")
                })
                Attribute("quantity", UInt64, "発注数量")
                Attribute("price", Float64, "発注価格 (LIMIT注文の場合)", func() {
                    Default(0)
                })
                Attribute("is_margin", Boolean, "信用取引かどうか", func() {
                    Default(false)
                })
                Required("symbol", "trade_type", "order_type", "quantity")
            })

            // レスポンス
            Result(func() {
                Description("ID of the created order")
                Attribute("order_id", String, "受付済み注文ID")
                Required("order_id")
            })

            // HTTPプロトコルとのマッピング
            HTTP(func() {
                POST("/order")
                Response(StatusCreated)
            })
        })
    })
    ```

### ステップ3: コード生成の実行
設計ファイルができたので、`goa`コマンドを使ってAPIサーバーの雛形コードを自動生成します。

```shell
goa gen stock-bot/design
```

このコマンドを実行すると、プロジェクトに`gen`ディレクトリが作成され、その中にサーバー関連のファイルが生成されます。

**重要**: 生成された`gen`ディレクトリはGitの管理対象から除外することを推奨します。`.gitignore`ファイルに以下の1行を追加してください。

```
/gen
```

### ステップ4: 依存関係の整理
Goaが使用するライブラリをプロジェクトの依存関係に追加します。

```shell
go mod tidy
```

### ステップ5: サービスロジックの実装
GoaはAPIのインターフェースを生成しましたが、実際の処理内容は私たちが実装する必要があります。

1.  `internal/interface/web/order_service.go` というファイルを新規作成します。
2.  以下の内容を貼り付けます。これは、APIが受け取ったリクエストを、既存の`OrderUsecase`に渡すための「接着剤」となるコードです。

    ```go
    package web

    import (
    	"context"
    	"log"
    	"stock-bot/internal/app"
    	ordersvr "stock-bot/gen/order"
    )
    
    // order.Serviceインターフェースを実装する構造体
    type OrderService struct {
        usecase app.OrderUseCase
        logger  *log.Logger 
    }

    // コンストラクタ
    func NewOrderService(usecase app.OrderUseCase, logger *log.Logger) ordersvr.Service {
        return &OrderService{
            usecase: usecase,
            logger:  logger,
        }
    }

    // Createメソッドの実装 (Goaが生成したインターフェースを満たす)
    func (s *OrderService) Create(ctx context.Context, p *ordersvr.CreatePayload) (res *ordersvr.OrderResult, err error) {
        s.logger.Printf("order.create method called with payload: %+v", p)

        // TODO: ここでUsecaseを呼び出すロジックを実装する
        // 例:
        // domainOrder, err := convertToDomainModel(p) // payloadからドメインモデルへ変換
        // if err != nil {
        //     return nil, err
        // }
        //
        // createdOrder, err := s.usecase.ExecuteOrder(ctx, *domainOrder)
        // if err != nil {
        //     return nil, err // Goaが適切なエラーレスポンスに変換してくれる
        // }
        //
        // res = &ordersvr.OrderResult{OrderID: createdOrder.OrderID}

        // 現時点ではダミーのレスポンスを返す
        dummyOrderID := "order-12345"
        res = &ordersvr.OrderResult{OrderID: &dummyOrderID}
        s.logger.Printf("order.create method successfully processed.")

        return res, nil
    }
    ```
    *注: Usecase呼び出し部分は、`convertToDomainModel`のような変換処理と合わせて後ほど実装します。*

### ステップ6: `main.go`の作成とサーバー起動
アプリケーションのエントリーポイントである`cmd/myapp/main.go`を作成し、Goaが生成したサーバーを起動するようにします。

`cmd/myapp/main.go`に以下のように記述してください。

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os"
    "os/signal"
    "stock-bot/internal/config" // configパッケージをインポート
import "stock-bot/internal/handler/web"
    "sync"
    "syscall"
    "time"

    ordersvr "stock-bot/gen/order/server"
    goahttp "goa.design/goa/v3/http"
    "goa.design/goa/v3/http/middleware"
)

func main() {
    // 1. ロガーのセットアップ
    logger := log.New(os.Stderr, "[stockbot] ", log.Ltime)

    // 2. 設定ファイルの読み込み
    cfg, err := config.LoadConfig(".env")
    if err != nil {
        logger.Fatalf("failed to load config: %s", err)
    }

    // 3. Usecaseなどの依存関係を初期化（今回はダミー）
    // orderUsecase := ...

    // 4. Goaサービスの実装を初期化
    orderSvc := web.NewOrderService(nil, logger) // Usecaseはまだなのでnilを渡す

    // 5. GoaのエンドポイントとHTTPハンドラを構築
    var wg sync.WaitGroup
    ctx, cancel := context.WithCancel(context.Background())

    orderEndpoints := ordersvr.NewEndpoints(orderSvc)
    
    errorHandler := func(ctx context.Context, w http.ResponseWriter, err error) {
		goahttp.ErrorHandler(ctx, w, err)
	}

    server := ordersvr.New(orderEndpoints, nil, goahttp.RequestDecoder, goahttp.ResponseEncoder, errorHandler, nil)
    mux := goahttp.NewMuxer()
    ordersvr.Mount(mux, server)

    // OpenAPI仕様配信用に静的ファイルサーバーをマウント
    fs := http.FileServer(http.Dir("./gen/http/openapi"))
    mux.Handle("/swagger/", http.StripPrefix("/swagger/", fs))

    // 6. HTTPサーバーの起動
    addr := fmt.Sprintf("http://localhost:%d", cfg.HTTPPort)
    u, err := url.Parse(addr)
    if err != nil {
        logger.Fatalf("invalid URL %#v: %s", addr, err)
    }

    srv := &http.Server{
        Addr:    u.Host,
        Handler: middleware.Log(logger)(mux), // リクエストログ用ミドルウェア
    }

    (*wg).Add(1)
    go func() {
        defer (*wg).Done()
        logger.Printf("HTTP server listening on %q", u.Host)
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            logger.Printf("server error: %s", err)
            cancel()
        }
    }()

    // 7. Graceful Shutdownの設定
    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

    select {
    case <-ctx.Done():
    case sig := <-c:
        logger.Printf("received signal %s, shutting down", sig)
    }
    
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer shutdownCancel()
    if err := srv.Shutdown(shutdownCtx); err != nil {
        logger.Printf("failed to shutdown server: %s", err)
    }

    (*wg).Wait()
    logger.Println("shutdown complete")
}

