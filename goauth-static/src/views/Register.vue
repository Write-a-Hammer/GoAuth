<template>
  <div class="register-container">
    <div class="register-box">
      <h1>注册</h1>
      <form @submit.prevent="handleRegister">
        <div class="input-group">
          <input type="text" v-model="username" placeholder="用户名" required />
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
      password: "",
      confirmPassword: "",
    };
  },
  methods: {
    async handleRegister() {
      if (this.password !== this.confirmPassword) {
        alert("两次输入的密码不一致！");
        return;
      }

      try {
        const response = await axios.post(`${apiBase}/auth/register`, {
          username: this.username,
          password: this.password,
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
};
</script>

<style scoped>
/* 注册页面样式 */
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f0f2f5;
}

.register-box {
  background: rgba(255, 255, 255, 0.9);
  padding: 40px;
  border-radius: 15px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.2);
  text-align: center;
  width: 350px;
}

.register-box h1 {
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
</style>
