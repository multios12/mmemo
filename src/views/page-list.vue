<!-- リストコンテナ -->
<template>
  <v-card class="mx-lg-18">
    <v-list class="mx-lg-4">
      <v-list-item v-for="(item, index) in items" :key="item.id">
        <v-list-item-content @click="rowClicked(index)">
          <v-list-item-subtitle v-text="item.date"></v-list-item-subtitle>
        </v-list-item-content>
        <v-list-item-content @click="rowClicked(index)">
          <v-list-item-title v-text="item.name"></v-list-item-title>
        </v-list-item-content>
        <v-list-item-content>
          <v-icon @click="deleteRow(index)">mdi-delete</v-icon>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-card>
</template>

<script lang="ts">
import { Component, Prop, Vue } from "vue-property-decorator";
import axios from "@/axiosForApi";
import router from "../router";
@Component
export default class EditView extends Vue {
  private items: { id: string; name: string; date: string }[] = [];
  private headers = [
    { text: "date", value: "date" },
    { text: "name", value: "name" },
    { text: "shop", value: "shop" },
  ];
  public created() {
    this.showListContainer();
  }
  public rowClicked(index: number) {
    router.push(`/${this.items[index].id}/`);
  }
  public async deleteRow(index: number) {
    await axios.delete(`./memos/${this.items[index].id}`);
    this.showListContainer();
  }
  public showListContainer() {
    axios.get("./memos").then((res) => (this.items = res.data));
  }
}
</script>
