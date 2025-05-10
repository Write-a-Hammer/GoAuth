import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';

export default defineConfig({
  plugins: [vue()],
  base: '/', // 设置基础路径
  define: {
    'process.env': {}, // 注入空的 process.env 对象
  },
  build: {
    outDir: '../backend/static', // 打包输出到后端项目的静态文件目录
    emptyOutDir: true, // 打包前清空目录
  },
  server: {
    host: 'localhost',
    port: 5173,
    strictPort: true, // 确保使用固定端口
    hmr: {
      protocol: 'ws',
      host: 'localhost',
    },
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // 后端 API 地址
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
      },
    },
  },
  envDir: '../', // 确保这里指向正确的 .env 文件目录
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
});
