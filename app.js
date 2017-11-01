var express = require("express");
var bodyParser = require('body-parser');
var app = express();
var mongo = require("mongodb");
var mongoClient = mongo.MongoClient;
var mongoUri = "mongodb://livaz:27017/hmemo";

app.set('view engine', 'ejs');
app.use('/assets', express.static('assets'));
app.use('/css', express.static('css'));
app.use('/fonts', express.static('fonts'));
app.use('/js', express.static('js'));

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

//#region Json API [memo]------------------------------------------------------
app.get("/memos", function (req, res, next) {
    var q = req.query.id ? { "_id": new mongo.ObjectId(req.query.id) } : {};

    // MongoDB へ 接続
    mongoClient.connect(mongoUri, (error, db) => {
        // コレクションに含まれるドキュメントをすべて取得
        db.collection("memos").find(q).toArray((error, documents) => {
            res.json(documents);
        });
    });
})

app.post("/memos", bodyParser.json(), function (req, res, next) {
    if (!req.body) return res.sendStatus(400);
    console.log(req.url);
    console.log(req.body);
    // MongoDB へ 接続
    mongoClient.connect(mongoUri, (error, db) => {
        var collection = db.collection("memos");
        var filter = req.body._id ? { "_id": new mongo.ObjectId(req.body._id) } : undefined;
        delete req.body._id;

        if (filter) {
            collection.updateOne(filter, req.body, (error, result) => {
                db.close();
                res.sendStatus(200);
            });
        } else {
            collection.insertOne(req.body, (error, result) => {
                db.close();
                res.sendStatus(200);
            });
        }
    });

});

app.delete("/memos", bodyParser.json(), function (req, res, next) {
    if (!req.query) return res.sendStatus(400);
    console.log("DELETE:" + req.url);

    // MongoDB へ 接続
    mongoClient.connect(mongoUri, (error, db) => {
        db.collection("memos").deleteOne({ "_id": new mongo.ObjectId(req.query._id) }, (error, result) => {
            db.close();
            res.sendStatus(200);
        });
    });
})
//#endregion

function resultAction(error, result) {
    db.close();
    res.sendStatus(200);
}

function stringsToObjectIds(strings) {
    for (let i = 0; i < strings.length; i++) {
        strings[i] = new mongo.ObjectId(strings[i]);
    }
    return strings;
}