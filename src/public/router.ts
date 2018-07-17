import Vue from 'vue';
import VueRouter from 'vue-router';
import editComponent from './components/page-edit.vue';
import loginComponent from './components/page-login.vue';
import mainComponent from './components/page-list.vue';
import viewComponent from './components/page-view.vue';
import auth from './auth';
import Axios from '../../node_modules/axios';

Vue.use(VueRouter);

var routes: any = [
    { component: mainComponent, path: '/', meta: { requiresAuth: true } },
    { component: editComponent, path: '/add', meta: { requiresAuth: true } },
    { component: loginComponent, path: '/login', meta: { requiresAuth: false } },
    { component: viewComponent, path: '/:id/', meta: { requiresAuth: true } },
    { component: editComponent, path: '/:id/edit', meta: { requiresAuth: true } },
];

var router = new VueRouter({ routes: routes });

// /**
//  * グローバルナビゲーションガード beforeイベント
//  */
// router.beforeEach(async (to, from, next) => {
//     // requireAuthメタフィールドによって、認証が必要か否かを判断する
//     if (!to.matched.some(record => record.meta.requiresAuth)) {
//         next();
//         return
//     }


//     // 認証が必要であれば、ログインページにリダイレクトする
//     if (!auth.loggedIn) {
//         var url = "/api/login";
//         console.log("開始");
//         var response = await Axios.get(url);

//         console.log(response.status);
//         if (response.status = 200) {
//             next();
//             return
//         }

//         next({ path: '/login', query: { redirect: to.fullPath } });
//     } else {
//         next();
//     }
// });

export default router;