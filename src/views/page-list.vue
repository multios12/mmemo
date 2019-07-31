<!-- リストコンテナ -->
<template>
  <v-data-table :headers="headers" :items="items"> </v-data-table>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import axios from "@/axiosForApi";
import router from "../router";
@Component
export default class EditView extends Vue {
  private items = [];
  private headers = [
    { text: "date", value: "date" },
    { text: "name", value: "name" },
    { text: "shop", value: "shop" }
  ];
  public created() {
    this.showListContainer();
  }
  public rowClicked(item: any, index: number, event: any) {
    router.push(`/${item.id}/`);
  }
  public async deleteRow(index: any) {
    await axios.delete(`./memos/${this.items[index]}`);
    this.showListContainer();
  }
  public showListContainer() {
    axios
      .get("./memos")
      .then(res => (this.items = res.data))
      .catch(r => router.push("/login"));
  }
}
</script>
