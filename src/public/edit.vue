<!-- 編集コンテナ -->
<template>
  <div id="editContainer" style="display: none;" class="container">
    <div class="panel panel-primary">
      <div class="panel-heading">edit</div>
      <div class="panel-body">
        <form id="editForm" name="editForm" class="form-horizontal" data-toggle="validator" role="form">
          <input name="_id" type="hidden" value="" />
          <div class="form-group has-feedback">
            <label for="inputName" class="control-label col-sm-1">name:</label>
            <div class="col-sm-11">
              <input type="text" class="form-control input-lg" name="name" id="inputName" size="14" required/>
              <div class="help-block with-errors"></div>
            </div>
          </div>
          <div class="form-group">
            <label for="date" class="control-label col-sm-1">date:</label>
            <div class="col-sm-11">
              <div class="input-group date">
                <input id="date" name="date" type="text" class="form-control" required>
                <span class="input-group-addon">
                  <i class="glyphicon glyphicon-th"></i>
                </span>
              </div>
              <div class="help-block with-errors"></div>
            </div>
          </div>
          <div class="form-group">
            <label class="control-label col-sm-1">shop:</label>
            <div class="col-sm-11 form-inline">
              <input type="text" size="8" name="shoptype" class="form-control col-sm-3" style="width:40%" autocomplete="on" list="test"
              />
              <datalist id="test">
                <option value="ヘルス">
                  <option value="デリヘル">
                    <option value="ソープ">
              </datalist>
              <input type="text" name="shop" class="form-control col-sm-9" style="width:60%" />
            </div>
          </div>
          <div class="form-group">
            <label class="control-label col-sm-1">home page:</label>
            <div class="col-sm-11">
              <input type="url" name="homepage" style="font-size:1vmax" class="form-control" />
            </div>
          </div>
          <div class="form-group">
            <label class="control-label col-sm-1">play:</label>
            <div class="col-sm-11">
              <textarea name="play" class="form-control"></textarea>
            </div>
          </div>
          <div class="form-group">
            <label class="control-label col-sm-1">talk:</label>
            <div class="col-sm-11">
              <textarea name="talk" class="form-control"></textarea>
            </div>
          </div>
          <div class="form-group">
            <div class="col-sm-11">
              <label class="control-label col-sm-1" style="padding:25px;">pictures:</label>
            </div>
            <div class="col-sm-11" style="overflow: hidden" id="imageGroup">
              <div id="imageBox" class="pictureSelectBox" draggable="false">
                <span class="glyphicon glyphicon-upload"></span>
                <input type="file" id="selectPicture" name="selectPicture" accept="image/jpeg, image/gif, image/png" class="inputFile" />
              </div>
            </div>
          </div>
          <div class="form-footer">
            <button id="editOkButton" type="submit" class="btn btn-lg btn-primary">
              <span class="glyphicon glyphicon-ok" style="font-size:32"></span> OK</button>
            <button id="editCloseButton" type="button" data-dismiss="modal" class="btn btn-lg btn-default">cancel</button>
          </div>
        </form>
      </div>
    </div>
  </div>
    
</template>

<script>
export default {};

var createObjectURL =
  window.URL && window.URL.createObjectURL
    ? file => window.URL.createObjectURL(file)
    : window.webkitURL && window.webkitURL.createObjectURL
      ? file => window.webkitURL.createObjectURL(file)
      : undefined;

var picturePaths = [];
document
  .getElementById("imageBox")
  .addEventListener("click", () =>
    document.getElementById("selectPicture").click()
  );
document
  .getElementById("selectPicture")
  .addEventListener("change", function(e) {
    var file = e.target.files[0];
    if (!file) return;

    var reader = new FileReader();
    reader.onload = function(e) {
      var url = createObjectURL ? createObjectURL(file) : e.target.result;
      var img = document.createElement("img");
      img.id = createRandomString(12);
      img.src = url;
      img.draggable = true;
      img.classList.add("picture");
      img.addEventListener("dragstart", ev => {
        ev.dataTransfer.setData("text", ev.target.id);
      });
      // ドラッグ中の要素がドロップ要素に重なった時
      img.addEventListener("dragover", function(ev) {
        ev.preventDefault();
        ev.dataTransfer.dropEffect = "move";
        dropArea.classList.add("dragover");
      });
      img.addEventListener("drop", function(ev) {
        ev.preventDefault();
        document
          .getElementById("imageGroup")
          .insertBefore(
            document.getElementById(ev.dataTransfer.getData("text")),
            this
          );
      });
      document.getElementById("imageGroup").appendChild(img);
    };
    reader.readAsDataURL(this.files[0]);
  });

function createRandomString(length) {
  var c = "abcdefghijklmnopqrstuvwxyz0123456789";
  var cl = c.length;
  var r = "";
  for (var i = 0; i < length; i++) {
    r += c[Math.floor(Math.random() * cl)];
  }
  return r;
}
</script>
