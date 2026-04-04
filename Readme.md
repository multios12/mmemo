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
4. 実行とデバッグで「Launch edge against localhost」を選択、実行

## デバッグ実行
vscode上での実行を前提。chromeを利用
フロントデバッグサーバ：ポート3000
バックエンドサーバ：ポート3001

開発環境では Go サーバを `data` ディレクトリをカレントディレクトリにして起動します。
そのため、開発用データは `data/settings.json` と `data/memo.db` に保存されます。
`data` ディレクトリおよび生成されるデータファイルは Git にコミットしません。

## 実行
`mmemo` は起動したカレントディレクトリに `settings.json` と `memo.db` を作成して利用します。

```
mmemo
mmemo -p :3001
```

ヘルプ:

```
mmemo -h
```

## settings.json

`settings.json` は mmemo の表示や入力ルールを決める設定ファイルです。

- 開発中は `data/settings.json` を編集します
- 初回起動時の雛形は `cmd/mmemo/static/.default.settings.json` です
- 本番実行時は、起動ディレクトリ直下の `settings.json` が使われます

### 役割

- `Holidays`
  カレンダーで祝日扱いにしたい日付一覧です
- `Categories`
  一覧・詳細・入力ルールを持つカテゴリ定義です

### 基本構造

```json
{
  "Holidays": [
    { "Date": "2026-01-01", "Name": "元日" }
  ],
  "Categories": [
    {
      "Key": "diary",
      "Name": "日記",
      "UseDate": true,
      "UseTag": true,
      "AllowMultipleEntriesPerDate": false,
      "Template": "## 今日の出来事\n----\n## 明日の予定\n----",
      "Fields": {
        "Name": "タイトル",
        "Date": "日付",
        "Tags": "タグ",
        "Value": "内容"
      }
    }
  ]
}
```

### Holidays

- `Date`
  `YYYY-MM-DD` 形式の日付
- `Name`
  祝日名。例: `元日`, `昭和の日`, `休日`

用途:

- リスト画面のカレンダーで、日曜日と同じ見た目で表示されます
- 名前自体は今後 tooltip や詳細表示にも使えるよう保持しています

参考:
- https://www8.cao.go.jp/chosei/shukujitsu/syukujitsu.csv
  <br/>祝日データ

### Categories

- `Key`
  内部キーです。URL や API で使われるので、作成後はなるべく変更しません
- `Name`
  画面に表示するカテゴリ名です
- `UseDate`
  `true` の場合、そのカテゴリは日付ベースで扱います
- `UseTag`
  `true` の場合、タグ入力を表示します
- `AllowMultipleEntriesPerDate`
  `false` の場合、同じカテゴリ・同じ日付のエントリは 1 件だけに制限されます
- `Template`
  新規作成時に本文へ入る初期テンプレートです
- `Fields`
  UI 上のラベル文言です
  `Outline` を指定すると、詳細画面の見出しラベルを変更できます
  `Date` を指定すると、詳細画面の日付ラベルを変更できます
  `Tags` を指定すると、詳細画面のタグラベルを変更できます
  `OutlineIcon` を指定すると、リスト画面の見出しモードで使うアイコンを変更できます

`OutlineIcon` の指定例:

- `folder-tree`
- `book-open`
- `file-lines`
- `note-sticky`
- `layer-group`
- `tags`

### よくある設定パターン

- 日記カテゴリ
  `UseDate: true`
  `AllowMultipleEntriesPerDate: false`
  1 日 1 件の運用向け
- 記録カテゴリ
  `UseDate: false`
  `AllowMultipleEntriesPerDate: true`
  日付に縛られないメモ向け

### 編集ルール

- `Date` は必ず `YYYY-MM-DD`
- JSON のキー名は Go / frontend で参照しているため変更しない
- カテゴリを追加するときは `Key` の重複を避ける
- 実運用データを変えたいときは `data/settings.json` を編集する
- 新規ユーザー向けデフォルトを変えたいときは `cmd/mmemo/static/.default.settings.json` も合わせて更新する

### エージェント向けメモ

- `Categories[].Key` は URL と API パスに使われるため、既存値の変更は破壊的です
- `AllowMultipleEntriesPerDate` は詳細画面の事前バリデーションと API 保存時の両方で使われます
- `Holidays` は `Calendar.svelte` に `HolidayDates` として渡されます
- `settings.json` に項目を追加・変更した場合は、この README の `settings.json` セクションも必ず更新してください
- 設定項目を追加したら次も更新すること:
  `pkg/entryapi/models.go`
  `front/models/settingType.ts`
  必要なら `cmd/mmemo/static/.default.settings.json`

## ビルド for Linux
> ./.devcontainer/build.sh

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
