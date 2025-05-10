package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBType             string
	MySQLDSN           string
	SQLitePath         string
	AllowedOrigin      string
	TurnstileSecretKey string
	TurnstileSiteKey   string // 添加 TurnstileSiteKey 字段
}

func LoadConfig() Config {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("加载 .env 文件失败: %v", err)
	}

	// 从环境变量中读取配置
	cfg := Config{
		DBType:             os.Getenv("DB_TYPE"),
		MySQLDSN:           os.Getenv("MYSQL_DSN"),
		SQLitePath:         os.Getenv("SQLITE_PATH"),
		AllowedOrigin:      os.Getenv("ALLOWED_ORIGIN"),
		TurnstileSecretKey: os.Getenv("TURNSTILE_SECRET_KEY"),
		TurnstileSiteKey:   os.Getenv("TURNSTILE_SITE_KEY"), // 加载 TurnstileSiteKey
	}

	// 检查必要的配置是否存在
	if cfg.DBType == "" {
		log.Fatal("环境变量 DB_TYPE 未设置")
	}
	if cfg.DBType == "mysql" && cfg.MySQLDSN == "" {
		log.Fatal("环境变量 MYSQL_DSN 未设置")
	}
	if cfg.DBType == "sqlite" && cfg.SQLitePath == "" {
		log.Fatal("环境变量 SQLITE_PATH 未设置")
	}
	if cfg.AllowedOrigin == "" {
		log.Fatal("环境变量 ALLOWED_ORIGIN 未设置")
	}
	if cfg.TurnstileSecretKey == "" {
		log.Fatal("环境变量 TURNSTILE_SECRET_KEY 未设置")
	}
	if cfg.TurnstileSiteKey == "" {
		log.Fatal("环境变量 TURNSTILE_SITE_KEY 未设置")
	}

	return cfg
}
