var request = require('supertest');
var app = require('../app');
var fs = require('fs');

describe('images getリクエスト テスト', () => {

    it('IDが指定されていない場合、404を返す', (done) => {
        request(app).get('/images//').expect(404)
            .end((err) => {
                if (err) throw err;
                done();
            });
    });

    it('存在しないIDが指定されている場合、404を返す', (done) => {
        request(app).get('/images/xxxxxx/').expect(404)
            .end((err) => {
                if (err) throw err;
                done();
            });
    });

    it('存在するIDが指定されている場合、ファイル名リストを返す', done => {
        if(!fs.existsSync('./data/')) fs.mkdirSync('./data/');
        if(!fs.existsSync('./data/1/')) fs.mkdirSync('./data/1/');
        fs.copyFileSync('./test/resources/jpeg01.jpg', './data/1/jpeg01.jpg');
        request(app).get('/images/1').expect(200)
            .end((err) => {
                if (err) throw err;
                done();
            });
    });
});
