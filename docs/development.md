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
4. 実行とデバッグで「Launch edge against localhost」を選択、実行

## デバッグ実行
vscode上での実行を前提。chromeを利用
フロントデバッグサーバ：ポート3000
バックエンドサーバ：ポート3001

開発環境では Go サーバを `data` ディレクトリをカレントディレクトリにして起動します。
そのため、開発用データは `data/settings.json` と `data/memo.db` に保存されます。
`data` ディレクトリおよび生成されるデータファイルは Git にコミットしません。

## テスト

Go テストはキャッシュ先を `/tmp/go-build` に固定した `./test.sh` を使います。

```sh
./test.sh
```

## ビルド for Linux
> ./.devcontainer/build.sh

`BASE_URL` を指定すると、サブパス配下へ配置する前提でフロントをビルドできます。

例:

```sh
./build.sh
BASE_URL=/mmemo/ ./build.sh
```

- `./build.sh`
  ルート配下で使う通常ビルドです
- `BASE_URL=/mmemo/ ./build.sh`
  `https://example.com/mmemo/` のようなサブパス配下へ配置するときに使います

-------------------------------------------------------------
## create new go project
> mkdir srv
> cd srv
> go mod init main

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

git tag -a v1.2.6 -m ''; git push origin --tags
