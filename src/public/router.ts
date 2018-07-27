import Vue from 'vue';
import VueRouter from 'vue-router';
import userStore from './userStore';

import editComponent from './components/page-edit.vue';
import loginComponent from './components/page-login.vue';
import mainComponent from './components/page-list.vue';
import viewComponent from './components/page-view.vue';

Vue.use(VueRouter);

var routes: any = [
    { component: mainComponent, path: '/', meta: { requiresAuth: true } },
    { component: editComponent, path: '/add', meta: { requiresAuth: true } },
    { component: loginComponent, path: '/login', meta: { requiresAuth: false } },
    { component: viewComponent, path: '/:id/', meta: { requiresAuth: true } },
    { component: editComponent, path: '/:id/edit', meta: { requiresAuth: true } },
];

var router = new VueRouter({ routes: routes });

/**
 * グローバルナビゲーションガード beforeイベント
 */
router.beforeEach(async (to, from, next) => {
    // localStorageにtokenが保存されていなければ、ログイン画面を表示する
    if (to.path == '/logout') {
        userStore.logout();
        next({ path: '/login' });
        return
    } else if (!userStore.token() && (to.path != '/login')) {
        next({ path: '/login' });
        return
    }

    next();
});

export default router;
