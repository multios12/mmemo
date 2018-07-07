import Vue from "vue";
import BootstrapVue from "bootstrap-vue";
import moment from "moment";
import axios from "axios";

import 'bootstrap/dist/css/bootstrap.css'
//import './bootstrap.min.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import navComponent from './components/app-nav.vue';
import router from './router'

Vue.use(BootstrapVue);

const app = new Vue({
    router,
    components: { "app-nav": navComponent },
    data: {
        settings: {}
    },
    created: function () {
        var self = this;
        axios.get("./settings").then(value => {
            self.settings = value.data
            self.selectedMonth = moment(new Date()).format("YYYY-MM");
        });
    },
    methods: {
    }
}).$mount('#app')