<!-- 編集コンテナ -->
<template>
  <v-flex xs12 sm8 offset-sm2 pb-4>
    <v-alert v-if="errorMessage" type="error" transition="scale-translation">{{
      errorMessage
    }}</v-alert>
    <v-card>
      <v-container>
        <v-form
          id="editForm"
          name="editForm"
          data-toggle="validator"
          role="form"
        >
          <v-text-field name="_id" type="hidden" value="" />
          <v-text-field label="name" required v-model="item.name" />
          <v-text-field label="shop" required v-model="item.shop" />
          <v-text-field label="home page" required v-model="item.page" />
          <v-text-field label="date" required v-model="item.date" type="date" />
          <v-textarea label="play" v-model="item.play" />
          <v-textarea label="talk" v-model="item.talk" />
          <v-card-actions>
            <v-btn color="primary" @click="regist"
              ><i class="far fa-check-circle"></i>OK</v-btn
            >
            <v-btn to="/">cancel</v-btn>
          </v-card-actions>
        </v-form>
      </v-container>
    </v-card>
  </v-flex>
</template>
<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import axios from "../axiosForApi";
import router from "../router";
@Component
export default class Edit extends Vue {
  private errorMessage: string = "";
  private item: {
    id: string;
    name: string;
    date: string;
    shop: string;
    page: string;
    play: string;
    talk: string;
  } = { id: "", name: "", date: "", shop: "", page: "", play: "", talk: "" };

  public created() {
    this.show();
  }
  public show() {
    const self = this;
    const url: string = this.$route.path;
    if (url == "/add") {
      this.item = {
        id: "",
        name: "",
        date: "",
        shop: "",
        page: "",
        play: "",
        talk: "",
      };
    } else {
      const res = axios
        .get(`./memos/${this.$route.params.id}`)
        .then((res) => (self.item = res.data));
    }
  }
  public async regist() {
    const self = this;
    if (!this.item.id) {
      await axios
        .put("./memos", this.item)
        .catch((res) => (this.errorMessage = "登録が失敗しました。"));
    } else {
      await axios.post(`./memos/${self.item.id}`, this.item);
    }
    this.errorMessage = "";
    router.push("/");
  }
}
</script>
