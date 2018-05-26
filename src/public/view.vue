<!-- ビューコンテナ -->
<template>
  <main role="main" class="container">
    <b-card title="view" v-model="targetItem">
      <b-form-group label="name"     >{{targetItem.name}}</b-form-group>
      <b-form-group label="date"     >{{targetItem.date}}</b-form-group>
      <b-form-group label="shop"     >{{targetItem.shop}}</b-form-group>
      <b-form-group label="home page">{{targetItem.page}}</b-form-group>
      <b-form-group label="play"     >{{targetItem.play}}</b-form-group>
      <b-form-group label="talk"     >{{targetItem.talk}}</b-form-group>
      <b-card-footer>
        <b-button>edit</b-button>
        <b-button>close</b-button>
      </b-card-footer>
    </b-card>
  </main>
</template>

<script lang="ts">
import Vue from "vue";
export default Vue.extend({
  data() {
    return {};
  },
  props:['targetItem'],
  methods: {
    cancel: () => this.$emit(),
    save: () => this.$emit()
  }
});
</script>
