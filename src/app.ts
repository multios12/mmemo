import express from 'express';
import fs from 'fs';
import path from 'path';
var app = express();
var createError = require('http-errors');

app.use(require('morgan')('dev'));
app.use(require('compression')());

// browser-sync Setup 
if (app.get('env') == 'development') app.use(require('connect-browser-sync')(require('browser-sync')({ "files": path.join(__dirname, "./public/*") })));

app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(express.static(path.join(__dirname, 'public')));
app.use('/memos', require('./routes/memos'));

app.use(function (req: express.Request, res: express.Response, next: Function) { next(createError(404)) });
app.use(function (err: any, req: any, res: any, next: any) {
  fs.readFile("./dist/public/error.html", 'utf-8', function (readerr, data) {
    data = data.replace('<%= message %>', err.message).replace('<%= error.stack %>', err.stack).replace('<%= error.status %>', err.status);
    res.send(data);
  })

});

module.exports = app;
