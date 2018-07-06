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
        settings: {},
        selectedMonth: null,
        selectedDate: null,
        balance: 0
    },
    created: function () {
        var self = this;
        axios.get("./settings").then(value => {
            self.settings = value.data
            self.selectedMonth = moment(new Date()).format("YYYY-MM");
        });
    },
    methods: {
        monthChanged: function (value: string) {
            this.selectedMonth = value;
        },
        dateChanged: function (value: string) {
            this.selectedDate = value;
        },
        balanceClicked: function (value: number) {
            this.balance = value;
        }
    }
}).$mount('#app')