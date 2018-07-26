<!-- リストコンテナ -->
<template>
  <main role="main">
  <b-container>
    <div class="table-responsive">
      <b-table striped hover :items="items" :fields="fields" @row-clicked="rowClicked">
          <template slot="name" slot-scope="data">{{data.value}}</template>
          <template slot="note" slot-scope="data">
            <b-button default @click="deleteRow(data.item)" size="sm"><i class="fas fa-minus-square"></i>delete</b-button>
          </template>
      </b-table>
    </div>
  </b-container>
  </main>
</template>

<script lang="ts">
import Vue from "vue";
import axios from "../axiosForApi";
import router from '../router'

export default Vue.extend({
  data() {
    return {
      items: [],
      fields: {
        date: { label: "date", sortable: true },
        name: { label: "name", sortable: true },
        shop: { label: "shop", sortable: true },
        note: { label: "    ", sortable: false }
      }
    };
  },
  created: function() {
    this.showListContainer();
  },
  methods: {
    rowClicked: function(item: any, index: number, event: any) {
      router.push(`/${item.id}/`);
    },
    deleteRow: async function(index: any) {
      await axios.delete(`./api/memos/${this.items[index]}`);
      this.showListContainer();
    },
    showViewModal: function(index: any) {
      this.parentNode.getAttribute("data-id");
    },
    showListContainer: function() {
      axios.get("./api/memos").then(res => this.items = res.data).catch(r => router.push('/login'));
    }
  }
});
</script>
