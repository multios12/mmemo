<template>
  <b-card title="view" v-model="targetItem">
    <b-form-group label="name">{{ targetItem.name }}</b-form-group>
    <b-form-group label="date">{{ targetItem.date }}</b-form-group>
    <b-form-group label="shop">{{ targetItem.shop }}</b-form-group>
    <b-form-group label="home page">{{ targetItem.page }}</b-form-group>
    <b-form-group label="play">{{ targetItem.play }}</b-form-group>
    <b-form-group label="talk">{{ targetItem.talk }}</b-form-group>
    <b-card-footer>
      <router-link class="btn btn-primary active" to="./edit">edit</router-link>
      <router-link class="btn btn-secondary active" to="/">close</router-link>
    </b-card-footer>
  </b-card>
</template>

<script lang="ts">
import axios from "../axiosForApi";
import { AxiosResponse } from "axios";
import { Component, Prop, Vue } from "vue-property-decorator";
import router from "../router";
@Component
export default class viewView extends Vue {
  private targetItem: {} = {};
  private created() {
    this.show();
  }
  private show() {
    const self = this;
    axios
      .get(`../api/memos/${this.$route.params.id}`)
      .then(res => (self.targetItem = res.data))
      .catch(r => router.push("/login"));
  }
}
</script>
