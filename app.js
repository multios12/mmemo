var express = require("express");
var bodyParser = require('body-parser');
var app = express();
var mongoClient = require("mongodb").MongoClient;
var mongoUri = "mongodb://livaz:27017/hmemo";

app.set('view engine', 'ejs');
app.use('/js', express.static('js'));
app.use('/css', express.static('css'));

/* 2. listen()メソッドを実行して3000番ポートで待ち受け。*/
var server = app.listen(3000, function () {
    console.log("Node.js is listening to PORT:" + server.address().port);

    if (process.env.MONGO_URI) {
        mongoUri = process.env.MONGO_URI;
    }
    console.log("MONGO_URI=" + mongoUri);
});

// "/"へのGETリクエストでindex.ejsを表示する。拡張子（.ejs）は省略されていることに注意。
app.get("/", function (req, res, next) {
    res.render("index", {});
});

app.get("/detail", function (req, res, next) {
    console.log(req.url);
    res.render("detail");
});

app.post("/detail", bodyParser.json(), function (req, res, next) {
    if (!req.body) return res.sendStatus(400);
    console.log(req.url);
    console.log(req.body);
    // MongoDB へ 接続
    mongoClient.connect(mongoUri, (error, db) => {
        // コレクションにドキュメントを挿入
        var collection = db.collection("memos");
        collection.insertOne(req.body
            , (error, result) => {
                db.close();
            });
    });

    res.sendStatus(200);
});

app.get("/list.json", function (req, res, next) {
    // MongoDB へ 接続
    mongoClient.connect(mongoUri, (error, db) => {
        // コレクションに含まれるドキュメントをすべて取得
        var collection = db.collection("memos");
        collection.find().toArray((error, documents) => {
            res.json(documents);
        });
    });
})

app.delete("/list.json", bodyParser.json(), function (req, res, next) {
    if (!req.body) return res.sendStatus(400);

    console.log(req.url);
    console.log(req.body);

    // MongoDB へ 接続
    mongoClient.connect(mongoUri, (error, db) => {
        req.body.forEach(function (element) {
            db.memos.remove({ _id: element });
        }, this);

    });

})