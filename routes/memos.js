var express = require('express');
var fs = require('fs');
var router = express.Router();

var NeDB = require('nedb');

var settings = require('../settings');
var memos = new NeDB({ filename: settings.MemosDBPath, autoload: true });

router.get('/', (req, res) => {
    var q = req.query.id ? { '_id': req.query.id } : {};
    memos.find(q, (err, docs) => res.json(docs));
});


router.post('/', (req, res) => {
    if (!req.body) return res.sendStatus(400);

    var filter = req.body._id ? { '_id': req.body._id } : undefined;
    delete req.body._id;

    // IDの作成
    if (req.body.qid == undefined) {
        req.body.qid = createRandomString(18);
    }
    var idPath = settings.CreateIdPath(req.body.qid);
    if (!fs.existsSync(idPath)) {
        fs.mkdirSync(idPath);
    }

    var body = writeImages(req.body);
    if (filter) {
        memos.update(filter, body, () => res.sendStatus(200));
    } else {
        memos.insert(body);
        res.sendStatus(200);
    }
});

router.delete('/', (req, res) => {
    if (!req.query) return res.sendStatus(400);
    memos.remove({ '_id': req.query._id }, () => res.sendStatus(200));
    var idPath = settings.CreateIdPath(req.body.qid);
    if (!fs.existsSync(idPath)) {
        fs.rmdirSync(settings.CreateIdPath(idPath));
    }
});

var writeImages = (body) => {
    if (body.images) {
        var fileNames = [];

        body.imageFileNames =[];
        for (var i = 0; i < body.images.length; i++) {
            var s = body.images[i].split(',');
            
            var filename;
            if(s[0] ==  "data:image/png;base64") {
                filename = (i + 1) + ".png";
            } else {
                filename = (i + 1);
            }

            var decode = new Buffer(s[1], 'base64');
            fs.writeFileSync(settings.CreateImagePath(body.qid, filename), decode);
            body.imageFileNames.push(filename);

        }
        delete body.images;
    }
    return body;
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


module.exports = router;
