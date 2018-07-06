import express from 'express';
import fs from 'fs';
import settings from '../config/settings';
var router = express.Router();

router.get('/:id', (req, res) => {
    if (!req.params.id) return res.sendStatus(404);

    var idPath = settings.CreateIdPath(req.params.id);
    console.log(idPath);
    if (!fs.existsSync(idPath)) return res.sendStatus(404);
    res.send(fs.readdirSync(idPath));
});

router.get('/:id/:fileName', (req, res) => {
    var filePath = settings.CreateImagePath(req.params.id, req.params.fileName);
    res.sendFile(filePath);
});

router.post('/:id/:fileName', (req, res) => {
    if (!req.body) return res.sendStatus(400);

    var idPath = settings.CreateIdPath(req.params.id);
    if (!fs.existsSync(idPath)) {
        fs.mkdirSync(idPath);
    }

    var filePath = settings.CreateImagePath(req.params.id, req.params.fileName);
    fs.writeFileSync(filePath, req.body);

    res.sendStatus(200);
});

router.delete('/:id', (req, res) => {
    if (!req.params.id) return res.sendStatus(400);

    var idPath = settings.CreateIdPath(req.params.id);
    fs.rmdirSync(idPath);
});
    /**
     * ランダムな文字列を作成する
     * @param length 文字列の長さ
     */
    var createRandomString= function(length:number ) {
        var c = "abcdefghijklmnopqrstuvwxyz0123456789";
        var cl = c.length;
        var r = "";
        for (var i = 0; i < length; i++) {
          r += c[Math.floor(Math.random() * cl)];
        }
        return r;
      }
module.exports = router;
