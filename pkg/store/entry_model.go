package store

type Entry struct {
	Id       int    `gorm:"primaryKey" ` // id
	Outline  string `validate:"min=1"` // outline
	Date     string `validate:"len=10"`  // 日付
	Category string // カテゴリ
	Tags     string // タグ
	Value    string // メモ
}

func (Entry) TableName() string {
	return "entries"
}
