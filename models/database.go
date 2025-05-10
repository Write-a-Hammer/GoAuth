package models

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
    _ "github.com/mattn/go-sqlite3"
    "goauth/config"
)

var DB *sql.DB

func InitDatabase(cfg config.Config) {
    var err error

    if cfg.DBType == "mysql" {
        DB, err = sql.Open("mysql", cfg.MySQLDSN)
    } else if cfg.DBType == "sqlite" {
        DB, err = sql.Open("sqlite3", cfg.SQLitePath)
    } else {
        log.Fatal("不支持的数据库类型")
    }

    if err != nil {
        log.Fatalf("数据库连接失败: %v", err)
    }

    // 测试数据库连接
    err = DB.Ping()
    if err != nil {
        log.Fatalf("数据库连接测试失败: %v", err)
    }

    fmt.Println("数据库连接成功")
}

func CloseDatabase() {
    if DB != nil {
        DB.Close()
    }
}