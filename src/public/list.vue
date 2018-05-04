  <!-- リストコンテナ -->
<template>
  <main role="main" class="container">
  <div id="listContainer" class="container">
    <div id="listSection" class="table-responsive">
      <table id="list" class="table table-striped table-hover">
        <thead>
          <tr>
            <th>date</th>
            <th>name</th>
            <th>shop</th>
            <th>*</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, index) in items" v-bind:key="item.no">
            <td>{{item.date}}</td>
            <td>{{item.name}}</td>
            <td>{{item.shop}}</td>
            <td><button type='button' class='btn btn-sm btn-default' @click="deleteRow(index)"><i class="fas fa-minus-square"></i>delete</button></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
  </main>
</template>
<script>
import Vue from "vue";
export default {
  data() {
    return {
      items: [{ date: "2018/01/01", name: "test", shop: "test" }]
    };
  },
  methods: {
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
