var request = require('supertest');
var images = require('../images');
var fs = require('fs');

describe('images readListメソッドテスト', () => {

    it('IDが指定されていない場合、404を返す', (done) => {
        var i = images.readList('1');
        images.readList().then((a) => {
            console.log(a);
        });
        console.log(i);
        done();
    });

});