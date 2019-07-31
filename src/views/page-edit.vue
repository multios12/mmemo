<!-- 編集コンテナ -->
<template>
  <v-container>
    <v-alert v-if="errorMessage" type="error" transition="scale-translation">{{
      errorMessage
    }}</v-alert>
    <v-flex xs12 sm8 offset-sm2 pb-4>
      <v-card>
        <v-container>
          <v-form
            id="editForm"
            name="editForm"
            data-toggle="validator"
            role="form"
          >
            <v-text-field name="_id" type="hidden" value="" />
            <v-text-field label="name" required v-model="targetItem.name" />
            <v-text-field label="shop" required v-model="targetItem.shop" />
            <v-text-field
              label="home page"
              required
              v-model="targetItem.page"
            />
            <v-text-field
              label="date"
              required
              v-model="targetItem.date"
              type="date"
            />
            <v-text-field label="play" v-model="targetItem.play" />
            <v-text-field label="talk" v-model="targetItem.talk" />
            <v-card-actions>
              <v-btn variant="primary" @click="regist"
                ><i class="far fa-check-circle"></i>OK</v-btn
              >
              <v-btn to="/">cancel</v-btn>
            </v-card-actions>
          </v-form>
        </v-container>
      </v-card>
    </v-flex>
  </v-container>
</template>
<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import axios from "../axiosForApi";
import router from "../router";
@Component
export default class Edit extends Vue {
  private errorMessage: string = "";
  private targetItem: {
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
      this.targetItem = {
        id: "",
        name: "",
        date: "",
        shop: "",
        page: "",
        play: "",
        talk: ""
      };
    } else {
      const res = axios
        .get(`./memos/${this.$route.params.id}`)
        .then(res => (self.targetItem = res.data))
        .catch(r => router.push("/login"));
    }
  }
  public async regist() {
    const self = this;
    if (!this.targetItem.id) {
      await axios
        .put("./memos", this.targetItem)
        .catch(res => (this.errorMessage = "登録が失敗しました。"));
    } else {
      await axios
        .post(`./memos/${self.targetItem.id}`, this.targetItem)
        .catch(r => router.push("/login"));
    }
    this.errorMessage = "";
    router.push("/");
  }
}
</script>
