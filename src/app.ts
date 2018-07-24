import express from 'express';
import fs from 'fs';
import path from 'path';
import {verifyMiddleware} from './auth';
var app = express();

app.use(require('morgan')('dev'));
app.use(require('compression')());
app.use(express.json());
app.use(express.urlencoded({ extended: false }));

app.use(express.static(path.join(__dirname, 'public','script.js')));
app.use(verifyMiddleware);
app.use(express.static(path.join(__dirname, 'public')));
app.use('/api/', require('./routes/api-login'));
app.use('/api/memos', require('./routes/api-memos'));

app.use(function (req, res, next) { next(require('http-errors')(404)) });
app.use(function (err: any, req: any, res: any, next: any) {
  fs.readFile("./dist/public/error.html", 'utf-8', function (readerr, data) {
    data = data.replace('<%= message %>', err.message).replace('<%= error.stack %>', err.stack).replace('<%= error.status %>', err.status);
    res.send(data);
  })
});

module.exports = app;
