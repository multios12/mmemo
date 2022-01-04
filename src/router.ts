import Vue from "vue";
import VueRouter from "vue-router";

import editComponent from "./views/page-edit.vue";
import mainComponent from "./views/page-list.vue";
import viewComponent from "./views/page-view.vue";

Vue.use(VueRouter);

const routes: any = [
  { component: mainComponent, path: "/" },
  { component: editComponent, path: "/add" },
  { component: viewComponent, path: "/:id/" },
  { component: editComponent, path: "/:id/edit" }
];

const router = new VueRouter({ routes });

export default router;
