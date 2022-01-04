<template>
  <v-app>
    <app-nav />
    <v-content>
      <router-view />
    </v-content>
  </v-app>
</template>

<script lang="ts">
import Vue from "vue";
import appNav from "./components/app-nav.vue";

export default Vue.extend({
  name: "App",
  components: {
    appNav
  },
  data: () => ({
    //
  })
});
</script>
