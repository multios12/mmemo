var express = require("express");
var bodyParser = require('body-parser');
var app = express();
var mongo = require("mongodb");
var mongoClient = mongo.MongoClient;
var mongoUri = "mongodb://livaz:27017/hmemo2";
var compression = require('compression');
const path = require('path');
var logger = require('morgan');

app.use(logger('dev'));
app.use(compression({ threshold: 0, level: 9, memLevel: 9 }));
app.use('/main.js', express.static(path.join(path.dirname(process.argv[1]), "js/main.js")));

/* listen()メソッドを実行して3000番ポートで待ち受け。*/
var server = app.listen(3000, function () {
    console.log("Node.js is listening to PORT:" + server.address().port);
    mongoUri = process.env.MONGO_URI ? process.env.MONGO_URI : mongoUri;
    console.log("MONGO_URI=" + mongoUri);
});

app.get("/", function (req, res, next) {
    var filePath = path.join(path.dirname(process.argv[1]), "views/index.html");
    res.sendFile(filePath);
});

//#region Json API [memo]------------------------------------------------------
var collectionName = "memos";
var pathName = "/" + collectionName;
app.get(pathName, function (req, res, next) {
    // MongoDB へ 接続
    mongoClient.connect(mongoUri, function (error, db) {
        var q = req.query.id ? { "_id": new mongo.ObjectId(req.query.id) } : {};
        db.collection(collectionName).find(q).toArray(function (error, documents) {
            res.json(documents);
        });
    });
})

app.post(pathName, bodyParser.json(), function (req, res, next) {
    if (!req.body) return res.sendStatus(400);
    console.log(req.url);
    console.log(req.body);
    // MongoDB へ 接続
    mongoClient.connect(mongoUri, function (error, db) {
        var collection = db.collection(collectionName);
        var filter = req.body._id ? { "_id": new mongo.ObjectId(req.body._id) } : undefined;
        delete req.body._id;

        if (filter) {
            collection.updateOne(filter, req.body, function (error, result) {
                db.close();
                res.sendStatus(200);
            });
        } else {
            collection.insertOne(req.body, function (error, result) {
                db.close();
                res.sendStatus(200);
            });
        }
    });

});

app.delete(pathName, bodyParser.json(), function (req, res, next) {
    if (!req.query) return res.sendStatus(400);
    console.log("DELETE:" + req.url);

    // MongoDB へ 接続
    mongoClient.connect(mongoUri, function (error, db) {
        db.collection(collectionName).deleteOne({ "_id": new mongo.ObjectId(req.query._id) }, function (error, result) {
            db.close();
            res.sendStatus(200);
        });
    });
})
//#endregion

function stringsToObjectIds(strings) {
    for (var i = 0; i < strings.length; i++) {
        strings[i] = new mongo.ObjectId(strings[i]);
    }
    return strings;
}