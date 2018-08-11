import express from "express";
import { sign } from "../auth";

const router = express.Router();

router.post("/login", (req: express.Request, res: express.Response) => {
    const token = sign(req.body.username, req.body.password);
    res.json({ token });
});

module.exports = router;
