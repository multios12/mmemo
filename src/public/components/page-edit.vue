<!-- 編集コンテナ -->
<template>
  <main role="main" class="container">
    <b-alert v-if="errorMessage">{{errorMessage}}</b-alert>
    <b-card title="edit">
      <b-form id="editForm" name="editForm" class="form-horizontal" data-toggle="validator" role="form">
        <input name="_id" type="hidden" value="" />
        <b-form-group label="name"     ><b-input required v-model="targetItem.name"/></b-form-group>
        <b-form-group label="date"     ><b-input required v-model="targetItem.date"/></b-form-group>
        <b-form-group label="shop"     ><b-input required v-model="targetItem.shop"/></b-form-group>
        <b-form-group label="home page"><b-input required v-model="targetItem.page"/></b-form-group>
        <b-form-group label="play"     ><b-textarea       v-model="targetItem.play"/></b-form-group>
        <b-form-group label="talk"     ><b-textarea       v-model="targetItem.talk"/></b-form-group>
        <div class="form-footer">
          <b-button variant="primary" id="editOkButton"    @click="regist"><i class="far fa-check-circle"></i>OK</b-button>
          <b-button variant="default" id="editCloseButton" @click="cancel">cancel</b-button>
        </div>
      </b-form>
  </b-card>
  </main>    
</template>
<script lang="ts">
import Vue from "vue";
import axios from "axios";
export default Vue.extend({
  data() {
    return {
      errorMessage: undefined
    };
  },
  props: ["targetItem"],
  methods: {
    regist: async function() {
      await axios
        .post("./api/memos")
        .catch(res => (this.errorMessage = "登録が失敗しました。"));
      this.errorMessage = undefined;
      this.$emit("change", "list-component");
    },
    cancel: function() {
      this.errorMessage = undefined;
      this.$emit("change", "list-component");
    }
  }
});
</script>
