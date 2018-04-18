import jQuery from 'jquery';
const $ = jQuery;
window.jQuery = jQuery;

import dt from 'datatables.net';
import 'datatables.net-dt/css/jquery.datatables.css';
import 'bootstrap/dist/js/bootstrap.js';
import 'bootstrap/dist/css/bootstrap.css';
import 'bootstrap-datepicker';
import moment from 'moment';
import 'moment/locale/ja';
import '../css/main.css';

var createObjectURL = (window.URL && window.URL.createObjectURL) ?
    (file) => window.URL.createObjectURL(file) :
    (window.webkitURL && window.webkitURL.createObjectURL) ?
        (file) => window.webkitURL.createObjectURL(file) : undefined;

window.onload = function () {
    editCloseButton_onclick();

    document.getElementById('registButton').onclick = () => showEditContainer();
    document.getElementById('editButton').onclick = () => showEditContainer(document.getElementsByName('viewForm')[0]._id.value);
    $('#date').datepicker({
        format: 'yyyy/mm/dd',
        language: 'ja',
        todayHighlight: true
    });

    document.getElementById('editOkButton').onclick = editOkButton_onclick;
    document.getElementById('editCloseButton').onclick = editCloseButton_onclick;

    // 画像関連イベントのデータ設定
    document.getElementById('imageBox').addEventListener('click', imageBox_click);
    document.getElementById('imageBox').addEventListener('drop', imageBox_drop);
    document.getElementById('imageBox').addEventListener('dragover', img_dragover);
    document.getElementById('imageBox').addEventListener('dragleave', img_dragleave);
    document.getElementById('selectPicture').addEventListener('change', selectPicture_change);
};

// #region イベント
/** 編集コンテナ OKボタンクリックイベント */
var editOkButton_onclick = () => {
    insertMemos();
    return false;
};

/** 編集コンテナ 閉じるボタンクリックイベント */
var editCloseButton_onclick = () => {
    document.getElementById('editContainer').style.display = 'none';
    $('#listContainer').fadeIn();

    var request = new XMLHttpRequest();
    request.onreadystatechange = function () {
        if (this.readyState != 4 || this.status != 200) {
            return;
        }
        if (!this.response) {
            return;
        }

        var table = document.getElementById('list');
        table.innerHTML = '<tr><th>date</th><th>name</th><th>shop</th><th></th></tr>';
        var datas = JSON.parse(this.response);

        datas.forEach(function (element) {
            var tr = table.insertRow(-1);
            tr.setAttribute('data-id', element._id);
            tr.style.cursor = 'pointer';
            var td1 = tr.insertCell(-1);
            td1.onclick = function(){ showViewModal(this.parentNode.getAttribute('data-id'));};
            td1.innerText = element.date ? element.date : '';
            var td2 = tr.insertCell(-1);
            td2.innerText = element.name ? element.name : '';
            td2.onclick = () => showViewModal(this.parentNode.getAttribute('data-id'));
            var td3 = tr.insertCell(-1);
            td3.innerText = element.shop ? element.shop : '';
            td3.onclick = () => showViewModal(this.parentNode.getAttribute('data-id'));
            var td4 = tr.insertCell(-1);
            td4.innerHTML = '<button type=\'button\' class=\'btn btn-sm btn-default\'>delete</button>';
            td4.children[0].onclick = () => deletememos(this.parentNode.parentNode.getAttribute('data-id'));
        }, this);
    };

    request.open('GET', './memos', true);
    request.send();
};

/** イメージ選択ボックス クリックイベント */
var imageBox_click = () => document.getElementById('selectPicture').click();

/**
 * イメージ選択ボックス ドロップ中イベント
 * @param {DragEvent} ev イベントデータ 
 */
var imageBox_drop = (ev) => {
    ev.preventDefault();
    var text = ev.dataTransfer.getData('text');
    document.getElementById(text).remove();

    var target = ev.target.id == 'imageBoxIcon' ? document.getElementById('imageBox') : ev.target;
    target.classList.remove('pictureover');
};

/**
 * ファイル選択 変更イベント
 * @param {*} e イベントデータ
 */
var selectPicture_change = (e) => {
    var file = e.target.files[0];
    if (!file) return;

    var fileReader = new FileReader();
    fileReader.onload = function (e) {
        var url = createObjectURL ? createObjectURL(file) : e.target.result;
        var img = document.createElement('img');
        img.id = createRandomString(12);
        img.src = url;
        img.draggable = true;
        img.classList.add('picture');
        img.addEventListener('dragstart', img_dragstart);
        img.addEventListener('dragend', img_dragend);
        img.addEventListener('dragover', img_dragover);
        img.addEventListener('dragleave', img_dragleave);
        img.addEventListener('drop', img_drop);
        document.getElementById('imageGroup').appendChild(img);
    };
    fileReader.readAsDataURL(document.getElementById('selectPicture').files[0]);
};
// #endregion

// #region イメージ イベント
/**
 * 画像 ドラッグ開始イベント
 * @param {DragEvent} ev イベントデータ 
 */
var img_dragstart = ev => {
    ev.dataTransfer.setData('text', ev.target.id);
    document.getElementById('imageBoxIcon').classList.remove('glyphicon-upload');
    document.getElementById('imageBoxIcon').classList.add('glyphicon-remove');
    document.getElementById('imageBoxIcon').style.color = '#fbb';

    document.getElementById('imageBoxIcon').dragDrop = true;
};

/**
 * 画像 ドラッグ完了イベント
 */
var img_dragend = () => {
    document.getElementById('imageBoxIcon').classList.remove('glyphicon-remove');
    document.getElementById('imageBoxIcon').classList.add('glyphicon-upload');
    document.getElementById('imageBoxIcon').style.color = '#bbb';

    document.getElementById('imageBoxIcon').draggable = false;
};

/**
 * 画像 ドラッグ重なりイベント
 * @param {DragEvent} ev イベントデータ 
 */
var img_dragover = ev => {
    ev.preventDefault();
    var target = ev.target.id == 'imageBoxIcon' ? document.getElementById('imageBox') : ev.target;
    ev.dataTransfer.dropEffect = 'move';
    target.classList.add('pictureover');
    // dropArea.classList.add('dragover');
};

/**
 * 画像 ドラッグ解放イベント
 * @param {DragEvent} ev イベントデータ
 */
var img_dragleave = ev => {
    var target = ev.target.id == 'imageBoxIcon' ? document.getElementById('imageBox') : ev.target;
    target.classList.remove('pictureover');
};

/**
 * 画像 ドロップイベント
 * @param {DragEvent} ev イベントデータ 
 */
var img_drop = ev => {
    ev.preventDefault();
    var text = ev.dataTransfer.getData('text');
    ev.target.classList.remove('pictureover');
    document.getElementById('imageGroup').insertBefore(document.getElementById(text), ev.target);
};
// #endregion

/** 
 * アラートメッセージを表示する
 * @param {String} message メッセージ 
 * @param {Boolean} isError trueの場合、エラーアラートを表示する
 */
function showMessage(message, isError) {
    var e = document.getElementById('message');
    e.classList.remove('alert-danger');
    e.classList.remove('alert-info');
    e.classList.add(isError ? 'alert-danger' : 'alert-info');
    e.innerText = message;
    e.style.display = 'block';
}

/**
 * アラートメッセージを閉じる
 */
function closeMessage() {
    var e = document.getElementById('message');
    e.classList.remove('alert-danger');
    e.classList.remove('alert-info');
    e.innerText = '';
    e.style.display = 'none';
}

/**
 * 指定したIDのデータをエディタコンテナに表示する
 * @param {String} id 
 */
function showEditContainer(id) {
    document.getElementById('listContainer').style.display = 'none';
    $('#editContainer').fadeIn();

    if (!id) {
        document.editForm.reset();
        document.editForm._id.value = '';
        closeMessage();
        return;
    }

    var request = new XMLHttpRequest();
    request.onloadend = function () {
        analyzeResponses(request, function () {
            jsonToForm(request.response, document.editForm);
        });
    };

    request.open('GET', './memos?id=' + id, true);
    request.setRequestHeader('Content-Type', 'application/json');
    request.send();
    closeMessage();
}

function showViewModal(id) {

    closeMessage();
    var request = new XMLHttpRequest();
    request.onloadend = function () {
        analyzeResponses(request, function () {
            var data = JSON.parse(request.response)[0];
            document.getElementsByName('viewForm')[0]._id.value = data._id;
            document.getElementById('viewname').innerText = data.name ? data.name : '';
            document.getElementById('viewdate').innerText = data.date ? data.date : '';
            document.getElementById('viewshoptype').innerText = data.shoptype ? data.shoptype : '';
            document.getElementById('viewshop').innerText = data.shop ? data.shop : '';
            document.getElementById('viewhomepage').innerHTML =
                data.homepage ? `<a href='${data.homepage}' target='_blank'>${data.homepage}</a>` : '';
            document.getElementById('viewplay').innerText = data.play ? data.play : '';
            document.getElementById('viewtalk').innerText = data.talk ? data.talk : '';

            $('#viewArea').modal();
        });
    };

    request.open('GET', `./memos?id=${id}`, true);
    request.setRequestHeader('Content-Type', 'application/json');
    request.send();
}

var insertMemos = () => {
    var request = new XMLHttpRequest();
    request.onloadend = () => {
        analyzeResponses(request, () => {
            showMessage('登録が完了しました。', false);
            editCloseButton_onclick();
        });
    };

    request.open('POST', './memos', true);
    request.setRequestHeader('Content-Type', 'application/json');
    var data = parseJson($('#editForm').serializeArray());
    var elements = document.getElementsByClassName('picture');
    var images = [];
    for (var i = 0; i < elements.length; i++) {
        var image = imageToBase64(elements[i], 'image/png');
        images.push(image);
    }
    data['images'] = images;
    request.send(JSON.stringify(data));
};

var deletememos = (id) => {
    var request = new XMLHttpRequest();
    request.onloadend = function () {
        analyzeResponses(request, () => showMessage('削除しました。', false));
        editCloseButton_onclick();
    };
    request.open('DELETE', `./memos?_id=${id}`, true);
    request.setRequestHeader('Content-Type', 'application/json');
    request.send();
};

/** レスポンスを処理する */
var analyzeResponses = (request, successAction) => {
    if (request.status == 200) {
        successAction();
    } else {
        showMessage(request.status + ':' + request.response, false);
    }
};

/**
 * 指定した配列から、連想配列を作成します。
 * @param {Array} data 配列データ
 * @returns 連想配列 
 */
var parseJson = function (data) {
    var returnJson = {};
    for (var idx = 0; idx < data.length; idx++) {
        if (data[idx].value == '') {
            continue;
        }
        returnJson[data[idx].name] = data[idx].value;
    }
    return returnJson;
};

var jsonToForm = (json, form) => {
    form.reset();
    var data = JSON.parse(json)[0];
    Object.keys(data).forEach(function (key) {
        form.elements[key].value = data[key] ? data[key] : '';
    });
};

/**
 * 指定した長さのランダムな文字列を作成する。
 * @param {number} length 長さ
 */
var createRandomString = length => {
    var c = 'abcdefghijklmnopqrstuvwxyz0123456789';
    var cl = c.length;
    var r = '';
    for (var i = 0; i < length; i++) {
        r += c[Math.floor(Math.random() * cl)];
    }
    return r;
};

/**
 * 
 * @param {HTMLImageElement} img img要素
 * @param {string} mime_type MimeType
 */
var imageToBase64 = (img, mime_type) => {
    // New Canvas
    var canvas = document.createElement('canvas');
    canvas.width  = img.width;
    canvas.height = img.height;
    // Draw Image
    var ctx = canvas.getContext('2d');
    ctx.drawImage(img, 0, 0);
    // To Base64
    return canvas.toDataURL(mime_type);
};
