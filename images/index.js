var fs = require('fs');
var settings = require('../settings');
var functions = {};

/**
 * 指定されたIDのファイルリストを返す
 * @param {string} id ファイルリストを返すID
 */
functions.readList = async (id) => {
    if (!id) return;

    var idPath = settings.CreateIdPath(id);
    console.log(idPath);
    const exists = fs.exists(idPath);
    if (!exists) return [];
    return await fs.readdir(idPath);
};

/**
 * 
 * @param {string} id ファイルリストを返すID
 * @param {string} fileName ファイル名
 */
functions.readFile = async (id, fileName) => {
    var filePath = settings.CreateImagePath(id, fileName);
    return await fs.readFile(filePath);
};

/**
 * 
 * @param {string} id ファイルリストを返すID
 * @param {string} fileName ファイル名
 * @param {string} body ファイルの内容
 */
functions.writeFile = async (id, fileName, body) => {
    if (!body) return;

    var idPath = settings.CreateIdPath(id);
    const exists = await fs.exists(idPath);
    if (! exists) {
        await fs.mkdir(idPath);
    }

    var filePath = settings.CreateImagePath(id, fileName);
    await fs.writeFile(filePath, body);
    return filePath;
};

/**
 * 
 * @param {string} id ファイルリストを返すID
 * @param {array} json 
 */
functions.writeFiles = async (id, json) => {
    for( var i = 0; i < json.length; i++) {
        var element = json[i];
        await functions.writeFile(id, element.fileName , element.body);
    }
};

/**
 * 
 * @param {string} id ファイルリストを返すID
 * @param {string} fileName ファイル名
 */
functions.removeFile = async (id, fileName) => {
    if (!id) return;
    var idPath = settings.CreateIdPath(id, fileName);
    await fs.delete(idPath);
};

module.exports = functions;
