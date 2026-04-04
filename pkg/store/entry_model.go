package store

import "time"

type Entry struct {
	Id        int       `gorm:"primaryKey"` // id
	Outline   string    `validate:"min=1"`  // outline
	Date      string    `validate:"len=10"` // 日付
	Category  string    // カテゴリ
	Tags      string    // タグ
	Value     string    // メモ
	CreatedAt time.Time `gorm:"<-:create"` // 作成日時
	UpdatedAt time.Time // 更新日時
}

func (Entry) TableName() string {
	return "entries"
}
