import Vue from "vue";
import Router from "vue-router";
import App from "@/App.vue";
import Register from "@/Register.vue";
import ForgotPassword from "@/ForgotPassword.vue";

Vue.use(Router);

export default new Router({
  mode: "history",
  routes: [
    {
      path: "/login",
      name: "Login",
      component: App,
    },
    {
      path: "/register",
      name: "Register",
      component: Register,
    },
    {
      path: "/forgot-password",
      name: "ForgotPassword",
      component: ForgotPassword,
    },
    {
      path: "*",
      redirect: "/login", // 默认跳转到登录页面
    },
  ],
});