import express from "express";
import fs from "fs";
import Settings from "../config/settings";
const router = express.Router();

router.get("/:id", (req, res) => {
    if (!req.params.id) { return res.sendStatus(404); }

    const idPath = Settings.CreateIdPath(req.params.id);
    console.log(idPath);
    if (!fs.existsSync(idPath)) { return res.sendStatus(404); }
    res.send(fs.readdirSync(idPath));
});

router.get("/:id/:fileName", (req, res) => {
    const filePath = Settings.CreateImagePath(req.params.id, req.params.fileName);
    res.sendFile(filePath);
});

router.post("/:id/:fileName", (req, res) => {
    if (!req.body) { return res.sendStatus(400); }

    const idPath = Settings.CreateIdPath(req.params.id);
    if (!fs.existsSync(idPath)) {
        fs.mkdirSync(idPath);
    }

    const filePath = Settings.CreateImagePath(req.params.id, req.params.fileName);
    fs.writeFileSync(filePath, req.body);

    res.sendStatus(200);
});

router.delete("/:id", (req, res) => {
    if (!req.params.id) { return res.sendStatus(400); }

    const idPath = Settings.CreateIdPath(req.params.id);
    fs.rmdirSync(idPath);
});

module.exports = router;
