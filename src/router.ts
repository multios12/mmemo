import Vue from "vue";
import VueRouter from "vue-router";
import userStore from "@/userStore";

import editComponent from "./views/page-edit.vue";
import mainComponent from "./views/page-list.vue";
import loginComponent from "./views/Login.vue";
import viewComponent from "./views/page-view.vue";

Vue.use(VueRouter);

const routes: any = [
  { component: mainComponent, path: "/", meta: { requiresAuth: true } },
  { component: editComponent, path: "/add", meta: { requiresAuth: true } },
  { component: loginComponent, path: "/login", meta: { requiresAuth: false } },
  { component: viewComponent, path: "/:id/", meta: { requiresAuth: true } },
  { component: editComponent, path: "/:id/edit", meta: { requiresAuth: true } }
];

const router = new VueRouter({ routes });

/**
 * グローバルナビゲーションガード beforeイベント
 */
router.beforeEach(async (to, from, next) => {
  if (to.path === "/logout") {
    // ログアウト処理を実施し、ログイン画面を表示する
    userStore.logout();
    next({ path: "/login" });
  } else if (!userStore.token() && to.path !== "/login") {
    // localStorageにtokenが保存されていなければ、ログイン画面を表示する
    next({ path: "/login" });
  } else {
    next();
  }
});

export default router;
