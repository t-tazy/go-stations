package config

import "os"

const (
	defaultPort   = ":8080"
	defaultDBPath = ".sqlite3/todo.db"
)

type Config struct {
	Port     string
	DbPath   string
	UserID   string
	Password string
}

// 環境変数から取得した情報を詰め込む
func New() *Config {
	cfg := Config{}
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	cfg.Port = port

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}
	cfg.DbPath = dbPath

	cfg.UserID = os.Getenv("BASIC_AUTH_USER_ID")
	cfg.Password = os.Getenv("BASIC_AUTH_PASSWORD")

	return &cfg
}
