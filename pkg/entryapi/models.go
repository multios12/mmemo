package entryapi

import "github.com/multios12/mmemo/pkg/store"

type SettingModel struct {
	Categories []CategoryModel // カテゴリ
}

// 情報種別
type CategoryModel struct {
	Key      string // キー
	Name     string `validate:"min=1"` // 種類名
	UseDate  bool   // 日付の使用・表示
	UseTag   bool   // タグ表示
	Template string // テンプレート
}

// 共有エントリ
type Entry = store.Entry
