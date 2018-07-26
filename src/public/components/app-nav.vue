<template>
    <b-navbar type="dark" class="navbar navbar-expand-md navbar-dark bg-dark fixed-top" v-if="this.$route.path != '/login'">
      <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>
      <b-navbar-brand><router-link class="text-white" style="text-decoration: none;" to="/">hmemo</router-link></b-navbar-brand>
      <b-collapse is-nav id="nav_collapse" class="mr-auto">
        <router-link variant="success" to="/add"><i class="fas fa-plus-square"></i>add</router-link>
      </b-collapse>
      <b-navbar-nav>
        <router-link variant="secondary" to="/logout"><i class="fas fa-sign-out-alt"></i></router-link>
      </b-navbar-nav>
    </b-navbar>
</template>

<script lang="ts">
import Vue from 'vue'
export default Vue.extend({
    
})
</script>
