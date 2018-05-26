<template>
  <div>
    <b-navbar type="dark" class="navbar navbar-expand-md navbar-dark bg-dark fixed-top">
      <b-navbar-brand @click="show('list-component')">{{title}}</b-navbar-brand>
      <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>

      <b-collapse is-nav id="nav_collapse">
        <b-button variant="success" @click="editingItem={}; show('edit-component')"><i class="fas fa-plus-square"></i>add</b-button>
      </b-collapse>
    </b-navbar>

    <component :is="currentPage" :targetItem="editingItem" @change="show" @select="select"></component>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import listComponent from "./list.vue"
import editComponent from "./edit.vue"
import viewComponent from "./view.vue"

Vue.component("list-component", listComponent)
Vue.component("edit-component", editComponent)
Vue.component("view-component", viewComponent)

export default Vue.extend({
  data() {
    return {
      title: "hmemo",
      currentPage: "list-component",
      editingItem: {},
    };
  },
  methods: {
    show: function(id:String) {
      this.currentPage = id;
      console.log(id);
    },
    select: function(i:any) {
      this.editingItem =i;
      this.show('edit-component');
    }
  }
})
</script>

<style>
#app {
  margin-top: 60px;
}
</style>
