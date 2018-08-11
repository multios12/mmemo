import BootstrapVue from "bootstrap-vue";
import Vue from "vue";

// import './bootstrap.min.css'
import "bootstrap-vue/dist/bootstrap-vue.css";
import "bootstrap/dist/css/bootstrap.css";

import navComponent from "./components/app-nav.vue";
import router from "./router";

Vue.use(BootstrapVue);

const app = new Vue({
    router,
    // tslint:disable-next-line:object-literal-sort-keys
    components: { "app-nav": navComponent },
    data: {
        isNav: true,
        settings: {},
    },
    created() {
        this.isNav = this.$route.path !== "/login";
    },
}).$mount("#app");
