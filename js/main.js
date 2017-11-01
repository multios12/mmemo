
$('#date').datepicker({
    format: "yyyy/mm/dd",
    language: "ja",
    todayHighlight: true
});

function showMessage(message, isError) {
    var e = document.getElementById("message");
    e.classList.remove("alert-danger");
    e.classList.remove("alert-info");
    e.classList.add(isError ? "alert-danger" : "alert-info");
    e.innerText = message;
    e.style.display = "block";
}

function closeMessage() {
    var e = document.getElementById("message");
    e.classList.remove("alert-danger");
    e.classList.remove("alert-info");
    e.innerText = "";
    e.style.display = "none";
}

function showmemos() {
    var request = new XMLHttpRequest();
    request.onreadystatechange = function () {
        if (this.readyState != 4 || this.status != 200) {
            return;
        }
        if (!this.response) {
            return;
        }

        var table = document.getElementById("list");
        table.innerHTML = "<tr><th>date</th><th>name</th><th>shop</th><th></th></tr>";
        var datas = JSON.parse(this.response);

        datas.forEach(function (element) {
            var tr = table.insertRow(-1);
            tr.setAttribute("data-id", element._id);
            tr.style.cursor = "pointer";
            var td1 = tr.insertCell(-1);
            td1.onclick = function () { showViewModal(this.parentNode.getAttribute("data-id")); }
            td1.innerText = element.date ? element.date : "";
            var td2 = tr.insertCell(-1);
            td2.innerText = element.name ? element.name : "";
            td2.onclick = function () { showViewModal(this.parentNode.getAttribute("data-id")); }
            var td3 = tr.insertCell(-1);
            td3.innerText = element.shop ? element.shop : "";
            td3.onclick = function () { showViewModal(this.parentNode.getAttribute("data-id")); }
            var td4 = tr.insertCell(-1);
            td4.innerHTML = "<button type='button' class='btn btn-sm btn-default' onclick=\"deletememos('" + element._id + "');\">delete</button>";
        }, this);
    }

    request.open('GET', './memos', true);
    request.send();
};

$('#editArea').on('hide.bs.modal', function (e) {
    //return false;
});

$('#editForm').on('invalid.bs.validator', function(e){
    document.getElementById("editOkButton").removeAttribute("data-dismiss"); 
});

$('#editForm').on('valid.bs.validator', function(e){
    document.getElementById("editOkButton").setAttribute("data-dismiss", "modal"); 
 });
 
$('#editForm').validator().on('submit', function (e) {
    if (e.isDefaultPrevented()) {
        // handle the invalid form...
        $(this).off;
        return false;
    } else {
        // everything looks good!
        insertMemos();
    }
});

function insertMemos() {
    var request = new XMLHttpRequest();
    request.onloadend = function () {
        analyzeResponse(request, function () {
            showMessage("登録が完了しました。", false);
        });
    }

    request.open('POST', './memos', true);
    request.setRequestHeader("Content-Type", "application/json");
    var data = parseJson($('#editForm').serializeArray());
    request.send(JSON.stringify(data));
};

function showEditModal(id) {
    if (!id) {
        document.editForm.reset();
        //document.editForm.date.value = moment(new Date()).format("YYYY/MM/DD");
        $('#editArea').modal();
        closeMessage();
        return;
    }

    var request = new XMLHttpRequest();
    request.onloadend = function () {
        analyzeResponse(request, function () {
            jsonToForm(request.response, document.editForm);
            $('#editArea').modal();
        });
    }

    request.open('GET', './memos?id=' + id, true);
    request.setRequestHeader("Content-Type", "application/json");
    request.send();
    closeMessage();

}

function showViewModal(id) {
    closeMessage();
    var request = new XMLHttpRequest();
    request.onloadend = function () {
        analyzeResponse(request, function () {
            var data = JSON.parse(request.response)[0];
            viewForm._id.value = data._id;
            document.getElementById("viewname").innerText = data.name ? data.name : "";
            document.getElementById("viewdate").innerText = data.date ? data.date : "";
            document.getElementById("viewshoptype").innerText = data.shoptype ? data.shoptype : "";
            document.getElementById("viewshop").innerText = data.shop ? data.shop : "";
            document.getElementById("viewhomepage").innerText = data.homepage ? data.homepage : "";
            document.getElementById("viewplay").innerText = data.play ? data.play : "";
            document.getElementById("viewtalk").innerText = data.talk ? data.talk : "";

            $('#viewArea').modal();
        });
    }

    request.open('GET', './memos?id=' + id, true);
    request.setRequestHeader("Content-Type", "application/json");
    request.send();
}

function deletememos(id) {
    var request = new XMLHttpRequest();
    request.onloadend = function () {
        analyzeResponse(request, function () {
            showMessage("削除しました。", false);
        });
    }
    request.open('DELETE', './memos?_id=' + id, true);
    request.setRequestHeader("Content-Type", "application/json");
    request.send();
};

/** レスポンスを処理する */
var analyzeResponse = function (request, successAction) {
    if (request.status == 200) {
        successAction();
        showmemos();
    } else {
        showMessage(request.status + ":" + request.response, false);
    }
}

var parseJson = function (data) {
    var returnJson = {};
    for (idx = 0; idx < data.length; idx++) {
        if (data[idx].value == "") {
            continue;
        }
        returnJson[data[idx].name] = data[idx].value
    }
    return returnJson;
}

var jsonToForm = function (json, form) {
    form.reset();
    var data = JSON.parse(json)[0];
    Object.keys(data).forEach(function (key) {

        form.elements[key].value = data[key] ? data[key] : "";
    });
}
