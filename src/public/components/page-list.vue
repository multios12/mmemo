<!-- リストコンテナ -->
<template>
  <main role="main">
  <b-container>
    <div class="table-responsive">
      <b-table striped hover :items="items" :fields="fields">
          <template slot="name" slot-scope="data"><a href="javascript:" @click="selectRow(data)">{{data.value}}</a></template>
          <template slot="note" slot-scope="data">
            <b-button default @click="deleteRow(data.item)"><i class="fas fa-minus-square"></i>delete</b-button>
          </template>
      </b-table>
    </div>
  </b-container>
  </main>
</template>

<script lang="ts">
import Vue from "vue";
import axios from "axios";
export default Vue.extend({
  data() {
    return {
      items: [],
      fields: {
        date: { label: "date", sortable: true },
        name: { label: "name", sortable: true },
        shop: { label: "shop", sortable: true },
        note: { label: "*", sortable: false }
      }
    };
  },
  props: ["targetItem"],
  created: function() {
    this.showListContainer();
  },
  methods: {
    selectRow: function(row: any) {
      this.$emit("select", row.item);
    },
    deleteRow: function(index: any) {
      this.items.splice(index);
    },
    showViewModal: function(index: any) {
      this.parentNode.getAttribute("data-id");
    },
    showListContainer: function() {
      axios.get("./api/memos").then(res => (this.items = res.data));
    }
  }
});
</script>
