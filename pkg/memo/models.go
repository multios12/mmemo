package memo

type SettingModel struct {
	Diary      DiaryModel
	Categories []CategoryModel // カテゴリ
}
type DiaryModel struct {
	Name string
}

// 情報種別
type CategoryModel struct {
	Key      string // キー
	Name     string `validate:"min=1"` // 種類名
	UseDate  bool   // 日付の使用・表示
	UseTag   bool   // タグ表示
	Template string // テンプレート
}

// メモ
type Memo struct {
	Id       int    `gorm:"primaryKey" ` // id
	Name     string `validate:"min=1"`   // name
	Date     string `validate:"len=10"`  // 日付
	Category string // カテゴリ
	Tags     string // タグ
	Value    string // メモ
}
