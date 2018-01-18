var assert = require('assert');
var path = require('path');

describe('settingsモジュールテスト', () => {
    var settings = require('../settings');

    it('CreateIdPathメソッド', () => {
        var idPath = path.dirname(__dirname);
        idPath = path.join(idPath,'data' ,'aaa');
        assert.equal(idPath, settings.CreateIdPath('aaa'));
    });

    it('CreateImagePathメソッド', () => {
        var idPath = path.dirname(__dirname);
        idPath = path.join(idPath,'data', 'aaa', 'file.name');
        assert.equal(idPath, settings.CreateImagePath('aaa', 'file.name'));
    });

    it('DatasPathプロパティ', () => {
        assert.notEqual(null, settings.DatasPath);
    });
});
