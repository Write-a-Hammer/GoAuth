<template>
  <div class="login-container">
    <div class="login-box">
      <h1>登录</h1>
      <form @submit.prevent="handleLogin">
        <div class="input-group">
          <input type="text" v-model="username" placeholder="用户名" required />
        </div>
        <div class="input-group">
          <input type="password" v-model="password" placeholder="密码" required />
        </div>
        <div class="input-group">
          <div id="turnstile-container"></div>
        </div>
        <button type="submit">登录</button>
      </form>
      <p>
        没有账号？
        <router-link to="/register">立即注册</router-link>
      </p>
      <p>
        忘记密码？
        <router-link to="/forgot-password">找回账号</router-link>
      </p>
      <p v-if="loginError" class="error-message">{{ loginError }}</p>
    </div>
  </div>
</template>

<script>
import axios from "axios";

const apiBase = import.meta.env.VITE_API_BASE;

export default {
  name: "Login",
  data() {
    return {
      username: "",
      password: "",
      turnstileToken: "",
      loginError: "",
    };
  },
  methods: {
    async handleLogin() {
      this.loginError = "";
      if (!this.turnstileToken) {
        this.loginError = "请完成验证码验证！";
        return;
      }
      try {
        await axios.post(`${apiBase}/auth/login`, {
          username: this.username,
          password: this.password,
          turnstileToken: this.turnstileToken,
        });
        alert("登录成功！");
        this.$router.push("/dashboard"); // 跳转到仪表盘页面
      } catch (error) {
        this.loginError = "登录失败，请检查用户名或密码。";
      }
    },
  },
  mounted() {
    axios.get(`${apiBase}/api/config`).then((response) => {
      const siteKey = response.data.TurnstileSiteKey;
      if (!siteKey) return;
      const script = document.createElement("script");
      script.src = "https://challenges.cloudflare.com/turnstile/v0/api.js";
      script.async = true;
      script.onload = () => {
        this.$nextTick(() => {
          const container = document.querySelector("#turnstile-container");
          if (!container) return;
          window.turnstile.render(container, {
            sitekey: siteKey,
            callback: (token) => {
              this.turnstileToken = token;
            },
          });
        });
      };
      document.head.appendChild(script);
    });
  },
};
</script>

<style scoped>
/* 确保背景覆盖整个页面 */
html,
body {
  height: 100%; /* 设置 html 和 body 的高度为 100% */
  margin: 0; /* 移除默认外边距 */
  font-family: "Inter", Arial, sans-serif;
  background: linear-gradient(135deg, #e3f2fd, #bbdefb); /* 渐变背景 */
}

/* 登录页面样式 */
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh; /* 确保容器高度覆盖整个视口 */
  margin: 0; /* 移除可能的外边距 */
}

.login-box {
  background: rgba(255, 255, 255, 0.7); /* 半透明白色背景 */
  padding: 40px;
  border-radius: 15px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.2);
  text-align: center;
  width: 350px;
  backdrop-filter: blur(10px); /* 添加模糊效果 */
}

.login-box h1 {
  margin-bottom: 30px;
  font-size: 24px;
  color: #333;
}

.input-group {
  margin-bottom: 20px;
}

.input-group input {
  width: 100%;
  padding: 12px 15px;
  font-size: 16px;
  border: 1px solid #ccc;
  border-radius: 8px;
  outline: none;
  transition: border-color 0.3s ease;
}

.input-group input:focus {
  border-color: #0078d7;
  box-shadow: 0 0 5px rgba(0, 120, 215, 0.5);
}

button {
  width: 100%;
  padding: 12px 15px;
  font-size: 16px;
  font-weight: bold;
  color: white;
  background-color: #0078d7;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.2s ease;
}

button:hover {
  background-color: #005a9e;
  transform: scale(1.02);
}

button:active {
  transform: scale(0.98);
}

.error-message {
  margin-top: 15px;
  color: #e74c3c;
  font-size: 14px;
}
</style>
