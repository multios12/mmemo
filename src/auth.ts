import express from 'express';
import fs from 'fs';
import jwt from 'jsonwebtoken';
import path from 'path';
/** 秘密鍵ファイルパス */
var keyFilePath = path.join(process.cwd(), './data/secretKey');
/** 秘密鍵文字列 */
var secretKey: Buffer;

// 秘密鍵を読み込む。ファイルが存在していなければ作成する
if (!fs.existsSync(keyFilePath)) {
    fs.writeFileSync(keyFilePath, createRandomString(20));
}
secretKey = fs.readFileSync(keyFilePath);

/**
 * 指定されたトークンを検証する
 * @param token トークン
 */
export function verify(token: string, callback: jwt.VerifyCallback) {
    return jwt.verify(token, secretKey, callback)
}

/**
 * トークンを検証するExpressミドルウェア
 * @param req リクエスト
 * @param res レスポンス
 * @param next NextFunction
 */
export function verifyMiddleware(req: express.Request, res: express.Response, next: Function) {

    if (req.path == '/api/login') {
        next();
        return
    }

    if ((req.path == '/' || req.path == '/script.js')) {
        next();
        return
    }

    var token = req.headers.authorization;
    token = token ? token.replace('Bearer ', '') : token;

    verify(token, function (err: Error, decoded: object) {
        if (err && req.path.indexOf('/api') > -1) {
            res.sendStatus(401);
            return
        }

        next();
    });
}

/**
 * 指定された値からトークンを作成する
 * @param payload トークンを作成する値
 */
export function sign(payload: object) {
    return jwt.sign(payload, secretKey, {expiresIn: '1d'});
}

/**
 * ランダムな文字列を生成する
 * @param len 作成する文字列の長さ
 */
function createRandomString(len: number) {
    var str = "abcdefghijklmnopqrstuvwxyz";
    var strLen = str.length;
    var result = "";

    for (var i = 0; i < len; i++) {
        result += str[Math.floor(Math.random() * strLen)];
    }

    return result;
}
