<template>
  <v-app>
    <app-nav />
    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import appNav from "./components/app-nav.vue";

@Component({components:{appNav}})
export default class App extends Vue {


}
</script>
