package infrastructure

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql" // MySQL ドライバ
)

// DBConfig はデータベース接続の設定を保持します
type DBConfig struct {
	User      string
	Password  string
	Host      string
	Port      string
	DBName    string
	Charset   string
	ParseTime string
	Loc       string
}

// LoadDBConfig は環境変数からデータベース設定を読み込みます
func LoadDBConfig() (*DBConfig, error) {
	if err := godotenv.Load("./api/config/.env"); err != nil {
		return nil, fmt.Errorf("Error loading .env file: %v", err)
	}

	return &DBConfig{
		User:      os.Getenv("MYSQL_USER"),
		Password:  os.Getenv("MYSQL_PASSWORD"),
		Host:      os.Getenv("MYSQL_HOST"),
		Port:      os.Getenv("MYSQL_PORT"),
		DBName:    os.Getenv("MYSQL_DATABASE"),
		Charset:   os.Getenv("MYSQL_CHARSET"),
		ParseTime: os.Getenv("MYSQL_PARSE_TIME"),
		Loc:       os.Getenv("MYSQL_LOC"),
	}, nil
}

// NewDB はデータベース接続を初期化します
func NewDB(cfg *DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
		cfg.ParseTime,
		cfg.Loc,
	)

	// sqlboiler用に標準のsql.DBを開く
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// コネクションプールの設定
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	// 実際に接続確認
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

