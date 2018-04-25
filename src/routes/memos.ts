import express from 'express';
import NeDB from 'nedb';
import settings from '../config/settings';

var router = express.Router();
var memos = new NeDB({ filename: settings.MemosDBPath(), autoload: true });

router.get('/', (req, res) => {
    var q = req.query.id ? { '_id': req.query.id } : {};
    memos.find(q, (err: any, docs: any) => res.json(docs));
});

router.post('/', (req, res) => {
    if (!req.body) return res.sendStatus(400);

    var filter = req.body._id ? { '_id': req.body._id } : undefined;
    delete req.body._id;

    if (filter) {
        memos.update(filter, req.body, undefined, () => res.sendStatus(200));
    } else {
        memos.insert(req.body);
        res.sendStatus(200);
    }
});

router.delete('/', (req, res) => {
    if (!req.query) return res.sendStatus(400);
    memos.remove({ '_id': req.query._id }, () => res.sendStatus(200));
});

module.exports = router;
