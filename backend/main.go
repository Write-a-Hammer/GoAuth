package main

import (
	"embed"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		panic("无法加载 .env 文件")
	}

	// 读取环境变量
	apiBase := os.Getenv("API_BASE")
	if apiBase == "" {
		apiBase = "http://localhost:8080"
	}

	// 提供静态文件服务
	http.Handle("/", http.FileServer(http.FS(staticFiles)))

	// 启动服务
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	println("服务启动，访问地址：http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}
