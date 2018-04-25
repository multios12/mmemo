<template>
    <!-- ビューモーダル -->
  <div id="viewArea" class="modal fade" tabindex="-1">
    <div class="modal-dialog modal-content">
      <div class="modal-header">
        <button type="button" class="close" data-dismiss="modal">
          <span aria-hidden="true">&times;</span>
        </button>
        <h4 class="modal-title" id="viewTitle">view</h4>
      </div>
      <!-- 5.モーダルのボディ -->
      <div class="modal-body">
        <form id="viewForm" name="viewForm" class="form-horizontal">
          <input name="_id" type="hidden" />
          <div class="form-group">
            <label class="control-label col-sm-1">name:</label>
            <p class="form-control-static col-sm-11" id="viewname" size="14"></p>
          </div>
          <div class="form-group">
            <label class="control-label col-sm-1">date:</label>
            <p class="form-control-static col-sm-11" id="viewdate"></p>
          </div>
          <div class="form-group">
            <label class="control-label col-sm-1">shop:</label>
            <p class="form-control-static col-sm-3" style="width:40%" id="viewshoptype"></p>
            <p class="form-control-static col-sm-8" style="width:60%" id="viewshop"></p>
          </div>
          <div class="form-group">
            <label class="control-label col-sm-1">home page:</label>
            <p class="form-control-static col-sm-11" id="viewhomepage" style="font-size:1vmax" class="form-control"></p>
          </div>
          <div class="form-group">
            <label class="control-label col-sm-1">play:</label>
            <p class="form-control-static col-sm-11" id="viewplay" class="form-control"></p>
          </div>
          <div class="form-group">
            <label class="control-label col-sm-1">talk:</label>
            <p class="form-control-static col-sm-11" id="viewtalk"></p>
          </div>
        </form>
      </div>
      <div class="modal-footer">
        <button type="button" id="editButton" data-dismiss="modal" class="btn btn-lg btn-info">
          <span class="glyphicon glyphicon-edit"></span>
          edit</button>
        <input type="button" data-dismiss="modal" value="close" class="btn btn-lg btn-default" />
      </div>
    </div>
  </div>
</template>
