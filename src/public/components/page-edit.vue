<!-- 編集コンテナ -->
<template>
  <main role="main" class="container">
    <b-alert v-if="errorMessage">{{errorMessage}}</b-alert>
    <b-card title="edit">
      <b-form id="editForm" name="editForm" class="form-horizontal" data-toggle="validator" role="form">
        <input name="_id" type="hidden" value="" />
        <b-form-group label="name"     ><b-input required v-model="targetItem.name"/></b-form-group>
        <b-form-group label="shop"     ><b-input required v-model="targetItem.shop"/></b-form-group>
        <b-form-group label="home page"><b-input required v-model="targetItem.page"/></b-form-group>
        <b-form-group label="date"     ><b-input required v-model="targetItem.date" type="date"/></b-form-group>
        <b-form-group label="play"     ><b-textarea       v-model="targetItem.play"/></b-form-group>
        <b-form-group label="talk"     ><b-textarea       v-model="targetItem.talk"/></b-form-group>
        <div class="form-footer">
          <b-button variant="primary" @click="regist"><i class="far fa-check-circle"></i>OK</b-button>
          <router-link class="btn btn-secondary active" to="/">cancel</router-link>
        </div>
      </b-form>
  </b-card>
  </main>    
</template>
<script lang="ts">
import Vue from "vue";
import axios from "../axiosForApi";
import router from "../router";
export default Vue.extend({
  data() {
    return {
      errorMessage: undefined,
      targetItem: {},
    };
  },
  created: function() {
    this.show();
  },
  methods: {
    show: function() {
      var self = this;
      var url: string = this.$route.path;
      if (url == '/add') {
        this.targetItem = {};
      } else {
        var res = axios.get(`./api/memos/${this.$route.params.id}`)
          .then(res => self.targetItem = res.data)
          .catch(r => router.push('/login'));
      }
    },
    regist: async function() {
      var self = this;
      if (!this.targetItem.id) {
        await axios
          .put("./api/memos", this.targetItem)
          .catch(res => (this.errorMessage = "登録が失敗しました。"));
      } else {
        await axios
          .post(`./api/memos/${self.targetItem.id}`, this.targetItem)
          .catch(r => router.push('/login'));
      }
      this.errorMessage = undefined;
      router.push('/');
    },
  }
});
</script>

<style>
main {
  padding-top: 4.5rem;
}
</style>