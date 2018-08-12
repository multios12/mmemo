import express from "express";
import fs from "fs";
import jwt from "jsonwebtoken";
import path from "path";

/** 秘密鍵ファイルパス */
const keyFilePath = path.join(process.cwd(), "./data/secretKey");

// 秘密鍵を読み込む。ファイルが存在していなければ作成する
if (!fs.existsSync(keyFilePath)) {
    console.log(`認証キー新規作成:${keyFilePath}`);
    fs.writeFileSync(keyFilePath, createRandomString(20));
}

/** 秘密鍵文字列 */
const secretKey: Buffer = fs.readFileSync(keyFilePath);

/**
 * 指定されたトークンを検証する
 * @param token トークン
 */
export function verify(token: string, callback: jwt.VerifyCallback) {
    return jwt.verify(token, secretKey, callback);
}

/**
 * トークンを検証するExpressミドルウェア
 * @param req リクエスト
 * @param res レスポンス
 * @param next NextFunction
 */
export function verifyMiddleware(req: express.Request, res: express.Response, next: FunctionConstructor) {

    if (req.path === "/api/login" || req.path === "/" || req.path === "/script.js") {
        next();
        return;
    }

    let token = req.headers.authorization;
    token = token ? token.replace("Bearer ", "") : token;

    verify(token, (err: Error, decoded: object) => {
        if (err !== undefined && err.name === "JsonWebTokenError") {
            console.log(`認証未完了　　　:${err.message}:token ${token}`);
        } else if (err !== undefined) {
            console.log(`認証エラー　　　:${err.message}:token ${token}`);
        }

        if (err && req.path.indexOf("/api") > -1) {
            res.sendStatus(401);
            return;
        }

        next();
    });
}

/**
 * 指定された値からトークンを作成する
 * @param payload トークンを作成する値
 */
export function sign(username: string, password: string) {
    if (username !== process.env.username || password !== process.env.password) {
        console.log(`認証失敗　　　　:${username}`);
        return undefined;
    }

    return jwt.sign({ username }, secretKey, { expiresIn: "1d" });
}

/**
 * ランダムな文字列を生成する
 * @param len 作成する文字列の長さ
 */
function createRandomString(len: number) {
    const str = "abcdefghijklmnopqrstuvwxyz";
    const strLen = str.length;
    let result = "";

    for (let i = 0; i < len; i++) {
        result += str[Math.floor(Math.random() * strLen)];
    }

    return result;
}
