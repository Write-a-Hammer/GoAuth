package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"goauth/config"
	"goauth/handlers"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 读取环境变量
	apiBaseURL := os.Getenv("API_BASE_URL")
	log.Printf("API Base URL: %s", apiBaseURL)

	// 加载配置
	cfg := config.LoadConfig()

	// 初始化数据库
	InitDatabase(cfg)
	defer CloseDatabase()

	// 设置路由
	http.HandleFunc("/auth/login", handlers.LoginHandler(cfg))
	http.HandleFunc("/api/bing-image", func(w http.ResponseWriter, r *http.Request) {
		// 添加 CORS 头
		w.Header().Set("Access-Control-Allow-Origin", cfg.AllowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// 处理 OPTIONS 请求（预检请求）
		if r.Method == http.MethodOptions {
			return
		}

		log.Printf("处理请求: %s", r.URL.Path)

		// 转发 Bing API 请求
		resp, err := http.Get("https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1")
		if err != nil {
			http.Error(w, "Failed to fetch Bing image", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		io.Copy(w, resp.Body)
	})

	http.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		// 添加 CORS 头
		w.Header().Set("Access-Control-Allow-Origin", cfg.AllowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// 处理 OPTIONS 请求（预检请求）
		if r.Method == http.MethodOptions {
			return
		}

		// 返回配置
		configResponse := map[string]string{
			"TurnstileSiteKey": cfg.TurnstileSiteKey, // 从配置中读取
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(configResponse)
	})

	http.HandleFunc("/auth/forgot-password", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "仅支持 POST 请求", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			Email string `json:"email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Printf("请求解析失败: %v", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		log.Printf("收到找回账号请求，邮箱: %s", request.Email)
		// 其他逻辑...

		// 检查邮箱是否存在
		user, err := FindUserByEmail(request.Email)
		if err != nil {
			http.Error(w, "邮箱未注册", http.StatusNotFound)
			return
		}

		// 生成重置令牌
		resetToken := GenerateResetToken() // 自定义方法生成令牌
		user.ResetToken = resetToken
		user.ResetTokenExpires = time.Now().Add(1 * time.Hour) // 1小时后过期
		if err := SaveUser(user); err != nil {
			http.Error(w, "Failed to save user", http.StatusInternalServerError)
			return
		}

		// 发送重置链接到邮箱
		resetLink := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("API_BASE_URL"), resetToken)
		if err := SendEmail(request.Email, "找回账号", fmt.Sprintf("点击以下链接重置密码：%s", resetLink)); err != nil {
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "重置链接已发送到您的邮箱"})
	})

	// 启动服务器
	port := ":8080"
	log.Printf("服务器启动，监听端口 %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// 初始化数据库连接
var db *gorm.DB

// CloseDatabase 关闭数据库连接
func CloseDatabase() {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("关闭数据库时出错: %v", err)
		return
	}
	sqlDB.Close()
}

func InitDatabase(cfg config.Config) {
	var err error
	if cfg.DBType == "sqlite" {
		db, err = gorm.Open(sqlite.Open(cfg.SQLitePath), &gorm.Config{})
	} else if cfg.DBType == "mysql" {
		db, err = gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
	}
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 自动迁移模型
	db.AutoMigrate(&User{})
}

// User 模型定义
type User struct {
	ID                uint   `gorm:"primaryKey"`
	Email             string `gorm:"uniqueIndex"`
	ResetToken        string
	ResetTokenExpires time.Time
}

func FindUserByEmail(email string) (*User, error) {
	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GenerateResetToken() string {
	token := make([]byte, 16)
	rand.Read(token)
	return hex.EncodeToString(token)
}

func SaveUser(user *User) error {
	result := db.Save(user)
	return result.Error
}

func SendEmail(to, subject, body string) error {
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	if smtpServer == "" || smtpPort == "" || smtpUser == "" || smtpPassword == "" {
		return fmt.Errorf("SMTP 配置不完整")
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpServer)
	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))
	addr := fmt.Sprintf("%s:%s", smtpServer, smtpPort)

	return smtp.SendMail(addr, auth, smtpUser, []string{to}, msg)
}
