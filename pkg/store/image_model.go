package store

import "time"

type Image struct {
	Path        string `gorm:"primaryKey"`
	ContentType string
	Data        []byte
	CreatedAt   time.Time `gorm:"<-:create"` // 作成日時
	UpdatedAt   time.Time // 更新日時
}

func (Image) TableName() string {
	return "images"
}
