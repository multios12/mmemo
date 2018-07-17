<template>
    <b-navbar type="dark" class="navbar navbar-expand-md navbar-dark bg-dark fixed-top">
      <b-navbar-brand><router-link class="text-white" style="text-decoration: none;" to="/">hmemo</router-link></b-navbar-brand>
      <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>
      <b-collapse is-nav id="nav_collapse">
        <router-link variant="success" to="/add"><i class="fas fa-plus-square"></i>add</router-link>
      </b-collapse>
    </b-navbar>
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
    
})
</script>
