<template>
  <v-container>
    <v-flex xs12 sm8 offset-sm2 pb-4>
      <v-card>
        <v-alert class="alert" type="error" effect="dark">
          username/password invalid.
        </v-alert>
        <form action=".">
          <v-container>
            <v-text-field
              v-model="username"
              label="username"
              autofocus
            ></v-text-field>
            <v-text-field
              type="password"
              v-model="password"
              label="Password"
            ></v-text-field>
            <v-btn color="success" @click="login" round>Login</v-btn>
          </v-container>
        </form>
      </v-card>
    </v-flex>
  </v-container>
</template>

<script lang="ts">
import { Component, Prop, Vue, Watch } from "vue-property-decorator";
import UserStore from "@/userStore";
import axios from "@/axiosForApi";
import router from "@/router";

@Component
export default class LoginView extends Vue {
  /** 入力されたユーザ名 */
  private username: string = "";
  /** 入力されたパスワード */
  private password: string = "";

  /** ビュー生成イベント */
  private created() {
    localStorage.removeItem("token");
  }

  /** ログインボタンクリックイベント */
  private async login() {
    const result = await UserStore.login(this.username, this.password);
    const alert = document.querySelector(".alert") as HTMLDivElement;
    alert.style.display = result ? "none" : "block";

    if (result) {
      const form = document.querySelector("form") as HTMLFormElement;
      form.submit();
    }
  }
}
</script>
