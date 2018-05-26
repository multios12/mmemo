<!-- リストコンテナ -->
<template>
  <main role="main">
  <b-container>
    <div id="listSection" class="table-responsive">
      <b-table striped hover :items="items" :fields="fields">
          <template slot="name"><a href="javascript:" @click="selectRow(item)">{{name}}</a></template>
          <template slot="note" slot-scope="data">
            <b-button default @click="index"><i class="fas fa-minus-square"></i>delete</b-button>
          </template>
      </b-table>
    </div>
  </b-container>
  </main>
</template>
<script>
import Vue from "vue";
export default {
  data() {
    return {
      items: [{ date: "2018/01/01", name: "test", shop: "test" }],
      fields: {
        date: { label: "date", sortable: true },
        name: { label: "name", sortable: true },
        shop: { label: "shop", sortable: true },
        note: { label: "*", sortable: false }
      }
    };
  },
  methods: {
    selectRow: function() {
      this.$emit('childs-event', 'hello!');
    },
    deleteRow: function(index) {
      this.items.splice(index);
    },
    showViewModal: function(index) {
      this.parentNode.getAttribute("data-id");
    },
    showListContainer: function() {
      // $("#listContainer").fadeIn();

      var request = new XMLHttpRequest();
      request.onreadystatechange = function() {
        if (this.readyState != 4 || this.status != 200) {
          return;
        }
        if (!this.response) {
          return;
        }
        this.items;
      };

      request.open("GET", "./memos", true);
      request.send();
    }
  }
};
</script>
