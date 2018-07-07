import express from 'express';
import db from '../db'

var router = express.Router();

router.get('/', (req, res) => {
    const sql = "SELECT * FROM memos";
    db.all(sql, (err: Error, rows: any[]) => res.json(rows));
});

router.get('/:id', (req, res) => {
    const sql = "SELECT * FROM memos WHERE id = ?";
    db.get(sql, [req.params.id], (err: Error, row: any) => res.json(row));
});

router.put('/', (req, res) => {
    var b = req.body;
    const sql = "INSERT INTO memos (name, date, shop, page, play, talk) VALUES (?, ?, ?, ?, ?, ?)";
    db.run(sql, [b.name, b.date, b.shop, b.page, b.play, b.talk], (err: Error) => res.status(200));
});

router.post('/:id', (req, res) => {
    const b = req.body;
    var sql ="UPDATE memos SET name=?, date =?, shop=?, page=?, play=?, talk=? WHERE id = ?"
    db.run(sql,[b.name, b.date, b.shop, b.page, b.play, b.talkb,req.params.id], (err) => res.status(200));
});

router.delete('/:id', (req, res) => {
    var sql = "DELETE FROM memos WHERE id=?";
    db.run(sql, req.params.id, (err) => res.status(200));
});

module.exports = router;
