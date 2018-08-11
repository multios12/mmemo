import express from "express";
import fs from "fs";
import path from "path";
import { verifyMiddleware } from "./auth";
const app = express();

app.use(express.json());
app.use(express.urlencoded({ extended: false }));

app.use(express.static(path.join(__dirname, "public", "script.js")));
app.use(verifyMiddleware);
app.use(express.static(path.join(__dirname, "public")));
// tslint:disable-next-line:no-var-requires
app.use("/api/", require("./routes/api-login"));
// tslint:disable-next-line:no-var-requires
app.use("/api/memos", require("./routes/api-memos"));

app.use((req, res, next) => next(require("http-errors")(404)));
app.use((err: any, req: any, res: any, next: any) => {
  fs.readFile("./dist/public/error.html", "utf-8", (readerr, data) => {
    data = data.replace("<%= message %>", err.message)
      .replace("<%= error.stack %>", err.stack)
      .replace("<%= error.status %>", err.status);
    res.send(data);
  });
});

module.exports = app;
