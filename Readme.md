# HMemo

## dockerでの実行
~~~
docker build --tag hmemo:0.7 .
~~~

~~~
 docker run -it --name hmemo -p 3001:3000 -e MONGO_URI=mongodb://mongo:27017/hmemo --link mongo:mongo hmemo:0.7
~~~

## 環境変数
MONGO_URI：データ保存用mongoDB接続URIを変更します。
