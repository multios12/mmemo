var request = require('supertest');
var app = require('../app');
var fs = require('fs');
var assert = require('assert');

describe('memos getリクエスト テスト', () => {

    it('empty', (done) => {
        fs.copyFileSync('./test/resources/memos_empty.db', './memos.db');
        request(app).get('/memos/')
            .expect(200)
            .end((err, res) => {
                if (err) throw err;
                done();
            });
    });

    it('1record', () => {
        fs.copyFileSync('./test/resources/memos_1record.db', './memos.db');
        request(app).get('/memos/').expect(200)
            .end((err, res) => {
                assert.equal(1, res.body.length);
                if (err) throw err;
            });
    });
});

describe('Memos postリクエスト テスト', () => {
    it('insert', () => {
        var data = { name: '名前', date: '2017-01-01', shoptype: 'ショップA', shop: 'ショップB', homepage: 'http://example', play: '', talk: '' };

        fs.copyFileSync('./test/resources/memos_empty.db', './memos.db');
        request(app).post('/memos/')
            .send(data)
            .expect(200)
            .end((err, res) => {
                assert.equal({}, res.body);
            });

    });
});
