<template>
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
</template>

<script>
export default {
  data() {
    return { picturePaths: [] };
  },
  methods: {
    imageBox_click: function() {
      document.getElementById("selectPicture").click();
    },
    selectPicture_change: function() {
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
    }
  }
};
</script>
