<template>
  <div>
    <nav class="navbar navbar-expand-md navbar-dark bg-dark fixed-top">
      <a class="navbar-brand" href="#list" @click="show('list-component')">{{ title }}</a>
      <button class="navbar-toggler" type="button"
          data-toggle="collapse" data-target="#navbarsExampleDefault" aria-controls="navbarsExampleDefault" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>

      <div class="collapse navbar-collapse" id="navbarsExampleDefault">
        <ul class="navbar-nav mr-auto">
          <li class="nav-item active">
            <a class="nav-link" href="#edit" @click="show('edit-component')"><i class="fas fa-plus-square"></i>add</a>
          </li>
        </ul>
      </div>
    </nav>

    <main role="main" class="container">
      <component :is="currentPage" @change="show"></component>
    </main>
  </div>
</template>

<script>
import Vue from 'vue'
import listComponent from './list.vue'
import editComponent from './edit.vue'
import viewComponent from './view.vue'

Vue.component('list-component', listComponent)
Vue.component('edit-component', editComponent)
Vue.component('view-component', viewComponent)
export default {
  data() {
    return {
      title: 'hmemo',
      currentPage : 'list-component'
    }
  },
  methods:{
    show: function (id) {
      this.currentPage = id
      console.log(id)
    }
  }
}
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}

h1,
h2 {
  font-weight: normal;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}
</style>
