<template>
  <div :style="appStyle" class="app-container">
    <div class="background-overlay"></div>
    <div v-if="isLoading" class="loading-overlay">
      <div class="spinner"></div>
      <p>加载每日一图中...</p>
    </div>
    <div v-else-if="backgroundImage === '/fallback-image.jpg'" class="error-message">
      <div class="error-card">
        <h2>加载失败</h2>
        <p>无法加载每日一图，请检查网络连接。</p>
        <button @click="retryFetch">重试</button>
      </div>
    </div>
    <router-view v-else />
  </div>
</template>

<script>
import { onMounted } from "vue";
import { useDailyImage } from "./composables/useDailyImage";
import "./assets/styles/AppStyles.css"; // 引入新的样式文件

const apiBase = import.meta.env.VITE_API_BASE;

export default {
  name: "App",
  setup() {
    const { backgroundImage, isLoading, appStyle, fetchDailyImage } = useDailyImage(
      apiBase
    );

    const retryFetch = () => {
      fetchDailyImage();
    };

    onMounted(() => {
      fetchDailyImage();
    });

    return {
      backgroundImage,
      isLoading,
      appStyle,
      retryFetch,
    };
  },
};
</script>
