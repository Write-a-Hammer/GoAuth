<template>
  <div v-if="isLoaded" class="login-container" :style="loginBgStyle">
    <div class="login-box">
      <h1>欢迎登录</h1>
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
      <p v-if="loginError" style="color: #e74c3c">{{ loginError }}</p>
    </div>
  </div>
  <div v-else class="loading-container">
    <div class="spinner"></div>
    <p>加载中...</p>
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
      backgroundImage: "",
      isLoaded: false,
      turnstileToken: "",
      loginError: "",
    };
  },
  computed: {
    loginBgStyle() {
      return {
        backgroundImage: `url(${this.backgroundImage || "/fallback-image.jpg"})`,
      };
    },
  },
  methods: {
    async fetchBingImage() {
      try {
        const response = await axios.get(`${apiBase}/api/bing-image`);
        const imageUrl = `https://www.bing.com${response.data.images[0].url}`;
        this.backgroundImage = imageUrl;
        const img = new Image();
        img.src = imageUrl;
        img.onload = () => {
          this.isLoaded = true;
        };
      } catch (error) {
        this.isLoaded = true;
      }
    },
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
      } catch (error) {
        this.loginError = "登录失败，请检查用户名或密码。";
      }
    },
  },
  mounted() {
    this.fetchBingImage();
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

<style scoped lang="css">
/* 样式代码 */
</style>
