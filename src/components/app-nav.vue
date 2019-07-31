<template>
  <v-toolbar app>
    <v-toolbar-title class="headline">
      <span>hmemo</span>
    </v-toolbar-title>
    <v-btn to="add"><i class="fas fa-plus-square"></i>add</v-btn>
    <v-sepacer></v-sepacer>
    <v-btn to="login" icon><v-icon class="fas fa-power-off"></v-icon></v-btn>
  </v-toolbar>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
@Component
export default class TopNav extends Vue {}
</script>
