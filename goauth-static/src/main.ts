import { createApp } from 'vue';
import App from './App.vue';
import router from './router/index.ts'; // 显式指定扩展名

createApp(App).use(router).mount('#app');