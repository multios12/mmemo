package store

import "time"

// Entry は日記やメモの本文を表す保存モデルです。
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

// TableName は Entry の保存先テーブル名を返します。
func (Entry) TableName() string {
	return "entries"
}

// Image はアップロード画像を表す保存モデルです。
type Image struct {
	Path        string `gorm:"primaryKey"`
	ContentType string
	Data        []byte
	CreatedAt   time.Time `gorm:"<-:create"` // 作成日時
	UpdatedAt   time.Time // 更新日時
}

// TableName は Image の保存先テーブル名を返します。
func (Image) TableName() string {
	return "images"
}

// Template はカテゴリごとのテンプレートを表す保存モデルです。
type Template struct {
	ID           int    `gorm:"primaryKey"`
	Category     string `gorm:"uniqueIndex:idx_template_category_name"`
	Name         string `gorm:"uniqueIndex:idx_template_category_name"`
	Value        string
	Tags         string
	DisplayOrder int
	CreatedAt    time.Time `gorm:"<-:create"`
	UpdatedAt    time.Time
}

// TableName は Template の保存先テーブル名を返します。
func (Template) TableName() string {
	return "templates"
}
