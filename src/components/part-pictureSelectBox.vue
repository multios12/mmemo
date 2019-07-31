<template>
  <b-form-group label="pictures" label-cols="sm-1">
    <div class="col-sm-11" style="overflow: hidden" id="imageGroup">
      <div id="imageBox" class="pictureSelectBox" draggable="false">
        <i class="glyphicon glyphicon-upload"></i>
        <b-form-file
          id="selectPicture"
          name="selectPicture"
          class="inputFile"
          accept="image/jpeg, image/gif, image/png"
        />
      </div>
    </div>
  </b-form-group>
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
