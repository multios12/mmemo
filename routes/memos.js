var express = require('express');
var router = express.Router();

var NeDB = require('nedb');
var memos = new NeDB({ filename: 'memos.db', autoload: true });

router.get('/',  (req, res, next) => {
    var q = req.query.id ? { "_id": req.query.id } : {};
    memos.find(q,
        function (err, docs) {
            console.log("[FIND]");
            console.log(docs);
            res.json(docs);
        }
    );
})

router.post('/', (req, res, next) => {
    if (!req.body) return res.sendStatus(400);
    console.log(req.url);
    console.log(req.body);

    var filter = req.body._id ? { "_id": req.body._id } : undefined;
    delete req.body._id;

    if (filter) {
        memos.update(filter, req.body, (error, result) => {
            res.sendStatus(200);
        });
    } else {
        memos.insert(req.body);
        res.sendStatus(200);
    }
})

router.delete('/', (req, res, next) => {
    if (!req.query) return res.sendStatus(400);
    console.log("DELETE:" + req.url);
    memos.remove({ "_id": req.query._id }, (error, result) => {
        res.sendStatus(200);
    });
})

module.exports = router;
