<template>
  <div class="register-container">
    <div class="register-box">
      <h1>注册</h1>
      <form @submit.prevent="handleRegister">
        <div class="input-group">
          <input type="text" v-model="username" placeholder="用户名" required />
        </div>
        <div class="input-group">
          <input type="email" v-model="email" placeholder="邮箱" required />
        </div>
        <div class="input-group">
          <input type="password" v-model="password" placeholder="密码" required />
        </div>
        <div class="input-group">
          <input
            type="password"
            v-model="confirmPassword"
            placeholder="确认密码"
            required
          />
        </div>
        <div class="input-group">
          <input type="text" v-model="verificationCode" placeholder="验证码" required />
          <button type="button" @click="sendVerificationCode" class="send-code-button">
            发送验证码
          </button>
        </div>
        <div id="turnstile-container" class="turnstile"></div>
        <button type="submit">注册</button>
      </form>
    </div>
  </div>
</template>

<script>
import axios from "axios";

const apiBase = import.meta.env.VITE_API_BASE;

export default {
  name: "Register",
  data() {
    return {
      username: "",
      email: "",
      password: "",
      confirmPassword: "",
      verificationCode: "",
      turnstileToken: "", // 用于存储 Turnstile 验证的 token
    };
  },
  methods: {
    async sendVerificationCode() {
      if (!this.email) {
        alert("请输入邮箱后再发送验证码！");
        return;
      }
      try {
        const response = await axios.post(`${apiBase}/auth/send-verification-code`, {
          email: this.email,
        });
        console.log("验证码已发送:", response.data);
        alert("验证码已发送到您的邮箱，请检查！");
      } catch (error) {
        console.error("验证码发送失败:", error.response?.data || error.message);
        alert("验证码发送失败，请重试！");
      }
    },
    async handleRegister() {
      if (this.password !== this.confirmPassword) {
        alert("两次输入的密码不一致！");
        return;
      }

      if (!this.turnstileToken) {
        alert("请完成验证码验证！");
        return;
      }

      try {
        const response = await axios.post(`${apiBase}/auth/register`, {
          username: this.username,
          email: this.email,
          password: this.password,
          verificationCode: this.verificationCode,
          turnstileToken: this.turnstileToken, // 将 Turnstile token 发送到后端
        });
        console.log("注册成功:", response.data);
        alert("注册成功，请返回登录页面！");
        this.$router.push("/"); // 跳转到登录页面
      } catch (error) {
        console.error("注册失败:", error.response?.data || error.message);
        alert("注册失败，请重试！");
      }
    },
  },
  mounted() {
    // 加载 Cloudflare Turnstile 脚本
    const script = document.createElement("script");
    script.src = "https://challenges.cloudflare.com/turnstile/v0/api.js";
    script.async = true;
    script.onload = () => {
      this.$nextTick(() => {
        const container = document.querySelector("#turnstile-container");
        if (!container) return;
        window.turnstile.render(container, {
          sitekey: "your-turnstile-site-key", // 替换为你的 Turnstile Site Key
          callback: (token) => {
            this.turnstileToken = token; // 获取 Turnstile 验证成功后的 token
          },
        });
      });
    };
    document.head.appendChild(script);
  },
};
</script>

<style scoped>
/* 确保背景覆盖整个页面 */
html,
body {
  height: 100%; /* 确保 html 和 body 的高度为 100% */
  margin: 0; /* 移除默认外边距 */
  padding: 0; /* 移除默认内边距 */
  font-family: "Inter", Arial, sans-serif;
  background: linear-gradient(135deg, #e3f2fd, #bbdefb); /* 渐变背景 */
  overflow: hidden; /* 防止滚动条出现 */
}

/* 注册页面样式 */
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh; /* 确保容器高度覆盖整个视口 */
  margin: 0; /* 移除可能的外边距 */
}

.register-box {
  background: rgba(255, 255, 255, 0.9); /* 半透明白色背景 */
  padding: 40px;
  border-radius: 15px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.2);
  text-align: center;
  width: 400px;
}

.register-box h1 {
  margin-bottom: 30px;
  font-size: 24px;
  color: #333;
}

.input-group {
  margin-bottom: 20px;
  position: relative;
}

.input-group input {
  width: calc(100% - 120px); /* 调整输入框宽度，给按钮留出更多空间 */
  padding: 10px 12px; /* 缩小内边距 */
  font-size: 14px; /* 缩小字体大小 */
  border: 1px solid #ccc;
  border-radius: 8px;
  outline: none;
  transition: border-color 0.3s ease;
}

.input-group input:focus {
  border-color: #0078d7;
  box-shadow: 0 0 5px rgba(0, 120, 215, 0.5);
}

.send-code-button {
  position: absolute;
  right: 0; /* 按钮靠右对齐 */
  top: 0;
  height: 100%; /* 按钮高度与输入框一致 */
  padding: 0 10px; /* 缩小按钮的左右内边距 */
  font-size: 12px; /* 缩小按钮字体大小 */
  color: white;
  background-color: #0078d7;
  border: none;
  border-radius: 0 8px 8px 0; /* 右侧圆角 */
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.2s ease;
}

.send-code-button:hover {
  background-color: #005a9e;
  transform: scale(1.02);
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

.turnstile {
  margin: 20px 0; /* 给验证码组件添加间距 */
}
</style>
