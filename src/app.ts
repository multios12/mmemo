var createError = require('http-errors');
import compression from "compression";
import express from 'express';
import fs from 'fs';
var path = require('path');
var cookieParser = require('cookie-parser');
var logger = require('morgan');

var app = express();

app.use(compression());
app.use(logger('dev'));

// browser-sync Setup 
if (app.get('env') == 'development') {
  var browserSync = require('browser-sync');
  var connectBrowserSync = require('connect-browser-sync');

  var browserSyncConfigurations = { "files": path.join(__dirname, "./public/*") };
  app.use(connectBrowserSync(browserSync(browserSyncConfigurations)));
}

app.use(express.json());
app.use(express.urlencoded({ extended: false }));
app.use(express.static(path.join(__dirname, 'public')));

app.use('/', require('./routes/index'));

app.use(function (req: express.Request, res: express.Response, next: Function) { next(createError(404)) });

app.use(function (err: any, req: any, res: any, next: any) {
  res.locals.message = err.message;
  res.locals.error = req.app.get('env') === 'development' ? err : {};

  fs.readFile("./dist/public/error.html", 'utf-8', function (readerr, data) {
    data = data.replace('<%= message %>', err.message).replace('<%= error.stack %>', err.stack).replace('<%= error.status %>', err.status);
    res.send(data);
  })

});

module.exports = app;
