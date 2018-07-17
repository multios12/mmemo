import express from 'express';
import fs from 'fs';
import path from 'path';
var app = express();

app.use(require('morgan')('dev'));
app.use(require('compression')());

app.use(express.json());
app.use(express.urlencoded({ extended: false }));
var sessionkey = { secret: 'secret', resave: false, saveUninitialized: false, cookie: { maxage: 1000 * 60 * 3600 } };
//app.use(require("express-session")(sessionkey));

app.use(express.static(path.join(__dirname, 'public')));
app.use('/api/', require('./routes/api-login'));

// 以下、認証が必要なURL
app.use('/api/memos', require('./routes/api-memos'));

app.use(function (req, res, next) { next(require('http-errors')(404)) });
app.use(function (err: any, req: any, res: any, next: any) {
  fs.readFile("./dist/public/error.html", 'utf-8', function (readerr, data) {
    data = data.replace('<%= message %>', err.message).replace('<%= error.stack %>', err.stack).replace('<%= error.status %>', err.status);
    res.send(data);
  })
});

module.exports = app;
