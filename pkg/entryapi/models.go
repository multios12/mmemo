package entryapi

import "github.com/multios12/mmemo/pkg/store"

type SettingModel struct {
	Holidays   []HolidayModel  // 祝日一覧
	Categories []CategoryModel // カテゴリ
}

type HolidayModel struct {
	Date string `validate:"min=1"` // 日付
	Name string `validate:"min=1"` // 名称
}

// 情報種別
type CategoryModel struct {
	Key                         string // キー
	Name                        string `validate:"min=1"` // 種類名
	UseDate                     bool   // 日付の使用・表示
	UseTag                      bool   // タグ表示
	AllowMultipleEntriesPerDate bool   // 同一日付の複数登録可否
	Fields                      map[string]string
	Template                    string // テンプレート
}

// 共有エントリ
type Entry = store.Entry
