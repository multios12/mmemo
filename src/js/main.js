import jQuery from 'jquery';
const $ = jQuery;
window.jQuery = jQuery;

import dt from 'datatables.net';
import 'datatables.net-dt/css/jquery.datatables.css';
import 'bootstrap/dist/js/bootstrap.js';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-datepicker';
import moment from 'moment';
import "moment/locale/ja";
import '../css/main.css';

window.onload = function () {
    showListContainer();

    document.getElementById("registButton").onclick = function () { showEditContainer(); };
    document.getElementById("editButton").onclick = function () { showEditContainer(viewForm._id.value); };
    $('#date').datepicker({
        format: "yyyy/mm/dd",
        language: "ja",
        todayHighlight: true
    });

    document.getElementById("editOkButton").onclick = function () {
        insertMemos();
        return false;
    }

    document.getElementById('editCloseButton').onclick = function () { showListContainer(); }
}

/** 
 * アラートメッセージを表示する
 * @param {*} message メッセージ 
 * @param {*} isError trueの場合、エラーアラートを表示する
 */
function showMessage(message, isError) {
    var e = document.getElementById("message");
    e.classList.remove("alert-danger");
    e.classList.remove("alert-info");
    e.classList.add(isError ? "alert-danger" : "alert-info");
    e.innerText = message;
    e.style.display = "block";
}

/**
 * アラートメッセージを閉じる
 */
function closeMessage() {
    var e = document.getElementById("message");
    e.classList.remove("alert-danger");
    e.classList.remove("alert-info");
    e.innerText = "";
    e.style.display = "none";
}

/**
 * リストコンテナを表示する
 */
function showListContainer() {
    document.getElementById("editContainer").style.display = "none";
    $("#listContainer").fadeIn();

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
            td4.innerHTML = "<button type='button' class='btn btn-sm btn-default'>delete</button>";
            td4.children[0].onclick = function () { deletememos(this.parentNode.parentNode.getAttribute("data-id")); }
        }, this);
    }

    request.open('GET', './memos', true);
    request.send();
};

/**
 * 指定したIDのデータをエディタコンテナに表示する
 * @param {String} id 
 */
function showEditContainer(id) {
    document.getElementById("listContainer").style.display = "none";
    $("#editContainer").fadeIn();

    if (!id) {
        document.editForm.reset();
        document.editForm._id.value = "";
        closeMessage();
        return;
    }

    var request = new XMLHttpRequest();
    request.onloadend = function () {
        analyzeResponses(request, function () {
            jsonToForm(request.response, document.editForm);
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
        analyzeResponses(request, function () {
            var data = JSON.parse(request.response)[0];
            viewForm._id.value = data._id;
            document.getElementById("viewname").innerText = data.name ? data.name : "";
            document.getElementById("viewdate").innerText = data.date ? data.date : "";
            document.getElementById("viewshoptype").innerText = data.shoptype ? data.shoptype : "";
            document.getElementById("viewshop").innerText = data.shop ? data.shop : "";
            document.getElementById("viewhomepage").innerHTML = data.homepage ? "<a href='" + data.homepage + "' target='_blank'>" + data.homepage + "</a>" : "";
            document.getElementById("viewplay").innerText = data.play ? data.play : "";
            document.getElementById("viewtalk").innerText = data.talk ? data.talk : "";

            $('#viewArea').modal();
        });
    }

    request.open('GET', './memos?id=' + id, true);
    request.setRequestHeader("Content-Type", "application/json");
    request.send();
}

function insertMemos() {
    var request = new XMLHttpRequest();
    request.onloadend = function () {
        analyzeResponses(request, function () {
            showMessage("登録が完了しました。", false);
            showListContainer();
        });
    }

    request.open('POST', './memos', true);
    request.setRequestHeader("Content-Type", "application/json");
    var data = parseJson($('#editForm').serializeArray());
    request.send(JSON.stringify(data));
};

function deletememos(id) {
    var request = new XMLHttpRequest();
    request.onloadend = function () {
        analyzeResponses(request, function () {
            showMessage("削除しました。", false);
        });
        showListContainer();
    }
    request.open('DELETE', './memos?_id=' + id, true);
    request.setRequestHeader("Content-Type", "application/json");
    request.send();
};

/** レスポンスを処理する */
function analyzeResponses(request, successAction) {
    if (request.status == 200) {
        successAction();
    } else {
        showMessage(request.status + ":" + request.response, false);
    }
}

var parseJson = function (data) {
    var returnJson = {};
    for (var idx = 0; idx < data.length; idx++) {
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
