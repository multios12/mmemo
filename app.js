var express = require('express');
var app = express();
var path = require('path');
var bodyParser = require('body-parser');
var morgan = require('morgan');

app.use(morgan('dev',{ immediate: true }));
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

app.use('/index.html', express.static(path.join(__dirname, 'dist/views/index.html')));
app.use('/main.js', express.static(path.join(__dirname, 'dist/main.js')));
app.use('/memos', require('./routes/memos'));

module.exports = app;
