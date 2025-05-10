<template>
  <div class="forgot-password-container">
    <div class="forgot-password-box">
      <h1>找回密码</h1>
      <form @submit.prevent="handleForgotPassword">
        <div class="input-group">
          <input type="email" v-model="email" placeholder="请输入您的邮箱" required />
        </div>
        <button type="submit">发送重置链接</button>
      </form>
    </div>
  </div>
</template>

<script>
import axios from "axios";

const apiBase = import.meta.env.VITE_API_BASE;

export default {
  name: "ForgotPassword",
  data() {
    return {
      email: "",
    };
  },
  methods: {
    async handleForgotPassword() {
      try {
        const response = await axios.post(`${apiBase}/auth/forgot-password`, {
          email: this.email,
        });
        console.log("重置链接已发送:", response.data);
        alert("重置链接已发送到您的邮箱，请检查！");
      } catch (error) {
        console.error("发送失败:", error.response?.data || error.message);
        alert("发送失败，请重试！");
      }
    },
  },
};
</script>

<style scoped>
/* 找回密码页面样式 */
.forgot-password-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: linear-gradient(135deg, #0078d7, #005a9e); /* 渐变背景 */
  color: #fff;
  font-family: "Roboto", Arial, sans-serif;
}

.forgot-password-box {
  background: rgba(255, 255, 255, 0.9); /* 半透明白色背景 */
  padding: 40px;
  border-radius: 15px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.2);
  text-align: center;
  width: 400px;
}

.forgot-password-box h1 {
  margin-bottom: 20px;
  font-size: 28px;
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
  border-color: #005a9e;
  box-shadow: 0 0 5px rgba(0, 90, 158, 0.5);
}

button {
  width: 100%;
  padding: 12px 15px;
  font-size: 16px;
  font-weight: bold;
  color: white;
  background-color: #005a9e;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.2s ease;
}

button:hover {
  background-color: #004080;
  transform: scale(1.02);
}

button:active {
  transform: scale(0.98);
}
</style>
