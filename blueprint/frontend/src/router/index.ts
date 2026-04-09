import { createRouter, createWebHistory } from "vue-router";
import LoginView from "@/views/LoginView.vue";
import HomeView from "@/views/HomeView.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/auth",
      component: LoginView,
      beforeEnter: () => {
        if (localStorage.getItem("token")) return "/";
      },
    },
    {
      path: "/",
      component: HomeView,
      beforeEnter: () => {
        if (!localStorage.getItem("token")) return "/auth";
      },
    },
  ],
});

export default router;
