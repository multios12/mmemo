# mmemo

プライベートな日記＆メモを記録するWebアプリです。
ブログのように第三者に公開することは想定していません。

- 日記及び、メモを作成
- MarkDownによる内容入力が可能
- 画像のアップロード
- テンプレートの登録・利用が可能

## 実行方法

GitHub Releases から配布される資材を展開し、`mmemo` バイナリを実行してください。
必要に応じて、環境変数でビルド設定を変更できます。

Release 資材の展開後、次のように起動できます。

```sh
./mmemo
```

Docker で動かす場合は、`docs/Dockerfile.sample` を参考にdockerfileを作成してください。
サンプルを使用して試す場合、下記のコマンドで実行できます。

```sh
docker build -f docs/Dockerfile.sample -t mmemo-run .
docker run --rm -v "$(pwd)/data:/app/data" mmemo-run
```

| 環境変数   | 既定値       | 説明 |
| ---------- | -----------  | ---- |
| `HTML`     | `index.html` | フロントのビルド対象 HTML を指定します。通常は変更不要です。 |
| `BASE_URL` | `./`         | サブパス配下へ配置する前提でフロントをビルドするときに使います。 |

`mmemo` 本体の通常起動では、追加の環境変数は必要ありません。

ローカルで起動した場合、下記のページからアクセスできます。
※ 既定ではローカルポート `:3000` で起動します。

```
http://localhost:3000
```

----------------------------------------------------------------

## Markdownエディタのショートカット

詳細画面の Markdown エディタでは、`Cmd` または `Ctrl` を使ったショートカットに対応しています。

- `Cmd/Ctrl + B`
  太字
- `Cmd/Ctrl + I`
  斜体
- `Cmd/Ctrl + Shift + X`
  取り消し線
- `Cmd/Ctrl + K`
  リンク
- `Cmd/Ctrl + Alt + 1`
  見出し1
- `Cmd/Ctrl + Alt + 2`
  見出し2
- `Cmd/Ctrl + Alt + 3`
  見出し3
- `Cmd/Ctrl + Alt + 0`
  本文
- `Cmd/Ctrl + Shift + 7`
  番号リスト
- `Cmd/Ctrl + Shift + 8`
  リスト
- `Cmd/Ctrl + Shift + 9`
  引用
- `Cmd/Ctrl + Shift + |`
  コード

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
      "Templates": [
        {
          "Name": "通常日記",
          "Value": "## 今日の出来事\n----\n## 明日の予定\n----",
          "Tags": ["日記"]
        }
      ],
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
- `Templates`
  新規作成時に使うテンプレート一覧です。1 件なら自動反映、複数なら選択UIが表示されます
  各テンプレートでは `Tags` も指定でき、新規作成時に本文と一緒に初期タグとして反映されます
  本文中に `----ここまで前回内容で置換` という行を入れると、その行はテンプレート適用時の境界として使われます
  タグカード・見出しカードの新規ボタンから開いた場合は、そのカード内で直前に更新されたエントリ本文で、この行より上の部分を置き換えます
  直前エントリが見つからない場合は、その行より上のテンプレート本文をそのまま使います
- `Fields`
  UI 上のラベル文言です
  `Outline` を指定すると、詳細画面の見出しラベルを変更できます
  `Date` を指定すると、詳細画面の日付ラベルを変更できます
  `Tags` を指定すると、詳細画面のタグラベルを変更できます
  `OutlineIcon` を指定すると、リスト画面の見出しモードで使うアイコンを変更できます

`OutlineIcon` の指定例:

- `note`
  汎用メモ向け。`fa-note-sticky`
- `tree`
  階層的な分類向け。`fa-folder-tree`
- `book`
  日記や記録向け。`fa-book`
- `group`
  グループやまとまり向け。`fa-layer-group`
- `tag`
  タグ中心のカテゴリ向け。`fa-tags`
- `calendar`
  日付ベースのカテゴリ向け。`fa-calendar-days`
- `document`
  文書メモ向け。`fa-file-lines`
- `list`
  箇条書き中心のカテゴリ向け。`fa-list-ul`
- `person`
  人物メモ向け。`fa-user`
- `user`
  `person` と同じ。`fa-user`
- `address-card`
  連絡先やプロフィール向け。`fa-address-card`

互換キーとして、次の値も受け付けます。

- `note-sticky`
- `folder-tree`
- `layer-group`
- `calendar-days`
- `file-lines`

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
  `pkg/web/models.go`
  `front/models/settingType.ts`
  必要なら `cmd/mmemo/static/.default.settings.json`

