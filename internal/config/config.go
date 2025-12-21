// internal/config/config.go
package config

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config はアプリケーションの設定を保持する構造体
type Config struct {
	TachibanaBaseURL  string `env:"TACHIBANA_BASE_URL"` // 立花証券APIベースURL
	TachibanaUserID   string `env:"TACHIBANA_USER_ID"`  // 立花証券ユーザーID
	TachibanaPassword string `env:"TACHIBANA_PASSWORD"` // 立花証券パスワード
	EventRid          string `env:"EVENT_RID"`          // EVENT I/F p_rid
	EventBoardNo      string `env:"EVENT_BOARD_NO"`     // EVENT I/F p_board_no
	EventNo           string `env:"EVENT_NO"`           // EVENT I/F p_e_no
	EventEvtCmd       string `env:"EVENT_EVT_CMD"`      // EVENT I/F p_evt_cmd
	DBHost            string `env:"DB_HOST"`            // データベースホスト名
	DBPort            int    `env:"DB_PORT"`            // データベースポート番号
	DBUser            string `env:"DB_USER"`            // データベースユーザー名
	DBPassword        string `env:"DB_PASSWORD"`        // データベースパスワード
	DBName            string `env:"DB_NAME"`            // データベース名
	LogLevel          string `env:"LOG_LEVEL"`          // ログレベル (debug, info, warn, error など)
	HTTPPort          int    `env:"HTTP_PORT"`          // HTTPサーバーポート番号
	WatchedStocks     []string
}

// LoadConfig は .env ファイルと環境変数から設定を読み込み、Config 構造体を返す
func LoadConfig(envPath string) (*Config, error) {
	// .env ファイルの読み込み (存在する場合)
	if envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			fmt.Printf("Error loading .env file: %v\n", err)
			// .env ファイルが読み込めなくても、環境変数から設定を読み込むので続行
		}
	}

	// 必須の設定項目 (エラーチェック)
	baseURLStr := os.Getenv("TACHIBANA_BASE_URL")
	if baseURLStr == "" {
		return nil, fmt.Errorf("TACHIBANA_BASE_URL is required")
	}
	_, err := url.Parse(baseURLStr) // URLとして有効かチェック
	if err != nil {
		return nil, fmt.Errorf("invalid TACHIBANA_BASE_URL: %w", err)
	}

	userID := os.Getenv("TACHIBANA_USER_ID")
	if userID == "" {
		return nil, fmt.Errorf("TACHIBANA_USER_ID is required")
	}

	password := os.Getenv("TACHIBANA_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("TACHIBANA_PASSWORD is required")
	}

	// オプションの設定項目 (デフォルト値)
	dbPort := GetInt("DB_PORT", 5432)
	httpPort := GetInt("HTTP_PORT", 8080)
	logLevel := GetString("LOG_LEVEL", "info")

	// 監視銘柄リストの読み込み
	watchedStocks, err := loadWatchedStocks("watched_stocks.csv")
	if err != nil {
		return nil, fmt.Errorf("failed to load watched stocks: %w", err)
	}

	return &Config{
		TachibanaBaseURL:  baseURLStr, // stringのまま
		TachibanaUserID:   userID,
		TachibanaPassword: password,
		EventRid:          GetString("EVENT_RID", ""),      // デフォルト値は空文字列
		EventBoardNo:      GetString("EVENT_BOARD_NO", ""), // デフォルト値は空文字列
		EventNo:           GetString("EVENT_NO", ""),       // デフォルト値は空文字列
		EventEvtCmd:       GetString("EVENT_EVT_CMD", ""),  // デフォルト値は空文字列
		DBHost:            GetString("DB_HOST", ""),        // デフォルト値は空文字列
		DBPort:            dbPort,
		DBUser:            GetString("DB_USER", ""),     // デフォルト値は空文字列
		DBPassword:        GetString("DB_PASSWORD", ""), // デフォルト値は空文字列
		DBName:            GetString("DB_NAME", ""),     // デフォルト値は空文字列
		LogLevel:          logLevel,
		HTTPPort:          httpPort,
		WatchedStocks:     watchedStocks,
	}, nil
}

// loadWatchedStocks はCSVファイルから監視銘柄のリストを読み込む
func loadWatchedStocks(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		// ファイルが存在しない場合はエラーとせず、空のリストを返す
		if os.IsNotExist(err) {
			fmt.Printf("Warning: %s not found, loading empty list of watched stocks\n", filePath)
			return []string{}, nil
		}
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// ヘッダー行を読み飛ばす
	if _, err := reader.Read(); err != nil {
		if err == io.EOF {
			return []string{}, nil // 空のファイル
		}
		return nil, err
	}

	var stocks []string
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if len(record) > 0 && record[0] != "" {
			stocks = append(stocks, record[0])
		}
	}

	fmt.Printf("Loaded %d watched stocks from %s\n", len(stocks), filePath)
	return stocks, nil
}


// GetInt は、環境変数から整数値を取得するヘルパー関数
func GetInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Printf("Warning: Invalid integer value for %s: %v. Using default value: %d\n", key, err, defaultValue)
		return defaultValue // エラーの場合はデフォルト値を返す
	}
	return value
}

// GetString は、環境変数から文字列値を取得するヘルパー関数
func GetString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
