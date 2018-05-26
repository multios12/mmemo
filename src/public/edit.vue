<!-- 編集コンテナ -->
<template>
  <main role="main" class="container">
    <b-container  id="editContainer">
      <b-card title="edit">
          <b-form id="editForm" name="editForm" class="form-horizontal" data-toggle="validator" role="form">
            <input name="_id" type="hidden" value="" />
            <b-form-group label="name"     ><b-input required v-model="item.name"></b-input></b-form-group>
            <b-form-group label="date"     ><b-input required v-model="item.date"></b-input></b-form-group>
            <b-form-group label="shop"     ><b-input required v-model="item.shop"></b-input></b-form-group>
            <b-form-group label="home page"><b-input required v-model="item.page"></b-input></b-form-group>
            <b-form-group label="play"     ><b-textarea       v-model="item.play"></b-textarea></b-form-group>
            <b-form-group label="talk"     ><b-textarea       v-model="item.talk"></b-textarea></b-form-group>
            <!-- <div class="form-group">
              <div class="col-sm-11">
                <label class="control-label col-sm-1" style="padding:25px;">pictures:</label>
              </div>
              <div class="col-sm-11" style="overflow: hidden" id="imageGroup">
                <div id="imageBox" class="pictureSelectBox" draggable="false">
                  <span class="glyphicon glyphicon-upload"></span>
                  <input type="file" id="selectPicture" name="selectPicture" accept="image/jpeg, image/gif, image/png" class="inputFile" />
                </div>
              </div>
            </div> -->
            <div class="form-footer">
              <b-button variant="primary" id="editOkButton" @click="regist"><span class="glyphicon glyphicon-ok"></span>OK</b-button>
              <b-button default id="editCloseButton" type="button" @click="cancel">cancel</b-button>
            </div>
          </b-form>
      </b-card>
    </b-container>
  </main>    
</template>

<script>
import Vue from "vue";
export default {
  methods: {
    regist: function() {
      this.$emit("change", "list-component");
    },
    /**
     * メインコンテナをリストに変更する
     */
    cancel: function() {
      this.$emit("change", "list-component");
    },
    /**
     * ランダムな文字列を作成する
     * @param length 文字列の長さ
     */
    createRandomString: function(length) {
      var c = "abcdefghijklmnopqrstuvwxyz0123456789";
      var cl = c.length;
      var r = "";
      for (var i = 0; i < length; i++) {
        r += c[Math.floor(Math.random() * cl)];
      }
      return r;
    }
  }
};
</script>
