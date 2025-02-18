# mmemo

----------------------------------------------------------------
## 開発

### 必要なソフトウェア
* Docker Desktop
* Visual Studio Code
* VIsual Studio拡張機能：Remote - Containers
* VIsual Studio拡張機能：Remote Development

### 開発環境立ち上げ手順
1. vscodeでフォルダを開く
2. CTRL+SHIFT+Pを押下して、コマンドパレット表示し、「Reopen in Container」を実行し、devContainerを開く
3. 実行とデバッグで「go API Server」を選択、実行
4. 実行とデバッグで「debug react」を選択、実行

## デバッグ実行
vscode上での実行を前提。chromeを利用
フロントデバッグサーバ：ポート3000
バックエンドサーバ：ポート3001

## ビルド for Linux
> ./.devcontainer/build.sh

-------------------------------------------------------------

## nitialize setting
> npm i -g yarn create-react-app

## create new react project
> create-react-app --typescvript

## create new go project
> mkdir srv
> cd srv
> go mod init main
> go get github.com/gin-gonic/gin
> go get github.com/go-playground/validator/v10

wget https://golang.org/dl/go1.22.2.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.22.2.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

### git comment
add   :新規機能追加
update:機能修正（バグ修正以外）
fix   :バグ修正
remove:削除
update: dependencies
  外部モジュール更新

## 正規表現メモ

* markdown画像の取得
```
(?:!\[([^[]+)\])(?:\((?:([^()\s]+)(?:\s"((?:[^"]*\\")*[^"]*)"\s*)?)\))$
```

* markdown画像の取得(URLが、"/api/images/tmp_"から始まるもののみ)
```
(?:!\[([^[]+)\])(?:\((?:(\/api\/images\/tmp_[^()\s]+)(?:\s"((?:[^"]*\\")*[^"]*)"\s*)?)\))$
```

```
![イメージ](/api/diary/2024/06/06/images/00000004.png)
![イメージ](/api/diary/2024/06/06/images/00000004.png "テスト")

![イメージ](/api/images/tmp_00000004.png)
![イメージ](/api/images/tmp_00000004.png "テスト")
```

git tag -a v1.2.3 -m ''; git push origin --tags
