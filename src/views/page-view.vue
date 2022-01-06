<template>
  <v-flex xs12 sm8 offset-sm2 pb-4>
    <v-card v-model="item">
      <v-card-title>{{ item.name }}</v-card-title>
      <v-card-subtitle>{{ item.date }}</v-card-subtitle>
      <v-card-text>
        <v-text-field name="_id" type="hidden" value="" />
        <v-text-field label="shop" filled disabled v-model="item.shop" />
        <v-text-field label="home page" filled disabled v-model="item.page" />
        <v-textarea label="play" filled disabled v-model="item.play" />
        <v-textarea label="talk" filled disabled v-model="item.talk" />
      </v-card-text>
      <v-card-actions>
        <v-btn color="primary" to="./edit">edit</v-btn>
        <v-btn to="/">close</v-btn>
      </v-card-actions>
    </v-card>
  </v-flex>
</template>

<script lang="ts">
import axios from "../axiosForApi";
import { Component, Prop, Vue } from "vue-property-decorator";
import router from "../router";
@Component
export default class viewView extends Vue {
  private item: {} = {};
  private created() {
    this.show();
  }
  private show() {
    const self = this;
    axios
      .get(`../api/memos/${this.$route.params.id}`)
      .then((res) => (self.item = res.data));
  }
}
</script>
