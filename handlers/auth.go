package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"goauth/config"
	"goauth/models"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	TurnstileResponse string `json:"cf-turnstile-response"`
}

type TurnstileVerifyResponse struct {
	Success bool `json:"success"`
}

func LoginHandler(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "请求解析失败", http.StatusBadRequest)
			return
		}

		// 验证 Turnstile
		if !verifyTurnstile(cfg.TurnstileSecretKey, req.TurnstileResponse) {
			http.Error(w, "Turnstile 验证失败", http.StatusForbidden)
			return
		}

		// 查询用户
		var hashedPassword string
		var userID int
		err = models.DB.QueryRow("SELECT id, password FROM users WHERE email = ?", req.Email).Scan(&userID, &hashedPassword)
		if err == sql.ErrNoRows {
			http.Error(w, "邮箱或密码错误", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}

		// 验证密码
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
		if err != nil {
			http.Error(w, "邮箱或密码错误", http.StatusUnauthorized)
			return
		}

		// 记录登录日志
		_, err = models.DB.Exec("INSERT INTO login_logs (user_id, login_time) VALUES (?, ?)", userID, time.Now())
		if err != nil {
			http.Error(w, "服务器错误", http.StatusInternalServerError)
			return
		}

		// 返回成功响应
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "登录成功")
	}
}

func verifyTurnstile(secretKey, response string) bool {
	url := "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	data := fmt.Sprintf("secret=%s&response=%s", secretKey, response)

	resp, err := http.Post(url, "application/x-www-form-urlencoded", io.NopCloser(strings.NewReader(data)))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var verifyResp TurnstileVerifyResponse
	err = json.NewDecoder(resp.Body).Decode(&verifyResp)
	if err != nil {
		return false
	}

	return verifyResp.Success
}
