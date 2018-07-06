import Vue from 'vue';
import VueRouter from 'vue-router';
import editComponent from './components/page-edit.vue';
import mainComponent from './components/page-list.vue';
import viewComponent from './components/page-view.vue';

Vue.use(VueRouter);
var routes: any = [
    { component: mainComponent, path: ''},
    { component: editComponent, path: '/add'},
    { component: viewComponent, path: '/:id/' },
    { component: editComponent, path: '/:id/edit' },
];

export default new VueRouter({ routes: routes });