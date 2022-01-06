<template>
  <v-app-bar app dark>
    <v-toolbar-title class="headline">
      <span>hmemo</span>
    </v-toolbar-title>
    <v-btn to="add"><i class="fas fa-plus-square"></i>add</v-btn>
    <v-spacer></v-spacer>
    <v-btn to="login" icon><v-icon class="fas fa-power-off"></v-icon></v-btn>
  </v-app-bar>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
@Component
export default class TopNav extends Vue {}
</script>
