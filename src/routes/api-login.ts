import { sign } from '../auth';
import express from 'express';
import db from '../db';
var router = express.Router();

router.post('/login', (req: express.Request, res: express.Response) => {

    if (req.body.username != process.env.username || req.body.password != process.env.password) {
        res.json({token:undefined});
        return
    }

    var token = sign({ username: req.body.username });
    res.json({ token: token });
});

module.exports = router;
