import Vue from "vue";
import BootstrapVue from "bootstrap-vue";
import axios from "axios";
import moment from "moment";

import 'bootstrap/dist/css/bootstrap.css'
//import './bootstrap.min.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import navComponent from './components/app-nav.vue';
import router from './router';

Vue.use(BootstrapVue);

const app = new Vue({
    router,
    components: { "app-nav": navComponent },
    data: { 
        settings: {},
        isNav:true 
    },
    created:function(){
        this.isNav = this.$route.path != '/login';
    },
}).$mount('#app')