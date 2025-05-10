<template>
  <div :style="appStyle" class="app-container">
    <div v-if="isLoading" class="loading-overlay">
      <div class="spinner"></div>
      <p>加载每日一图中...</p>
    </div>
    <router-view v-else />
  </div>
</template>

<script>
import axios from "axios";

const apiBase = import.meta.env.VITE_API_BASE;

export default {
  name: "App",
  data() {
    return {
      backgroundImage: "", // 存储每日一图的 URL
      isLoading: true, // 是否正在加载图片
    };
  },
  computed: {
    appStyle() {
      return {
        backgroundImage: `url(${this.backgroundImage || "/fallback-image.jpg"})`,
        backgroundSize: "cover",
        backgroundPosition: "center",
        backgroundRepeat: "no-repeat",
        minHeight: "100vh",
        filter: this.isLoading ? "blur(10px)" : "none", // 加载时模糊背景
        transition: "filter 0.5s ease",
      };
    },
  },
  methods: {
    async fetchDailyImage() {
      try {
        const response = await axios.get(`${apiBase}/api/bing-image`);
        const imageUrl = `https://www.bing.com${response.data.images[0].url}`;
        const img = new Image();
        img.src = imageUrl;
        img.onload = () => {
          this.backgroundImage = imageUrl;
          this.isLoading = false; // 图片加载完成
        };
      } catch (error) {
        console.error("每日一图加载失败:", error);
        this.backgroundImage = "/fallback-image.jpg"; // 加载失败时使用默认图片
        this.isLoading = false;
      }
    },
  },
  mounted() {
    this.fetchDailyImage();
  },
};
</script>

<style>
/* 全局样式 */
body {
  margin: 0;
  font-family: "Roboto", Arial, sans-serif;
  background-color: #f0f2f5;
  color: #333;
  line-height: 1.6;
}

.app-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f0f2f5;
  color: #333;
  position: relative;
}

.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background-color: rgba(255, 255, 255, 0.8);
  z-index: 10;
}

.spinner {
  width: 50px;
  height: 50px;
  border: 5px solid #ccc;
  border-top-color: #0078d7;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

a {
  color: #0078d7;
  text-decoration: none;
  transition: color 0.3s ease;
}

a:hover {
  color: #005a9e;
}

button {
  font-family: inherit;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.3s ease, transform 0.2s ease;
}

input {
  font-family: inherit;
  font-size: 16px;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  outline: none;
  transition: border-color 0.3s ease, box-shadow 0.3s ease;
}

input:focus {
  border-color: #0078d7;
  box-shadow: 0 0 5px rgba(0, 120, 215, 0.5);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  background: rgba(255, 255, 255, 0.8); /* 半透明背景 */
  border-radius: 15px;
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.2);
}

h1,
h2,
h3,
h4,
h5,
h6 {
  margin: 0;
  font-weight: 600;
}

p {
  margin: 0 0 1em;
}

ul {
  list-style: none;
  padding: 0;
}

li {
  margin: 0 0 0.5em;
}
</style>
