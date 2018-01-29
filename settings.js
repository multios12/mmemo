var fs = require('fs');
var path = require('path');
var settings = {};

/** データ保存パス */
var _datasPath;

/** データ保存パスを返す */
settings.DatasPath = () => {
    if(!_datasPath) {
        _datasPath = process.env.DATAS_PATH ? process.env.DATAS_PATH : path.join(__dirname, 'data');
        if(!fs.existsSync(_datasPath)) {
            fs.mkdirSync(_datasPath);
        }
    }
    return _datasPath;
};

/** 指定されたIDのデータ保存パスを返す */
settings.CreateIdPath = (id) => path.join(settings.DatasPath(), id);

/** 指定されたID・ファイル名のデータ保存パスを返す */
settings.CreateImagePath =  (id, fileName) => path.join(settings.DatasPath(), id, fileName);

/** memos.dbファイルのパスを返す */
settings.MemosDBPath = path.join(settings.DatasPath(), 'memos.db');

module.exports = settings;
