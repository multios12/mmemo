package store

import "time"

type Template struct {
	ID           int    `gorm:"primaryKey"`
	Category     string `gorm:"uniqueIndex:idx_template_category_name"`
	Name         string `gorm:"uniqueIndex:idx_template_category_name"`
	Value        string
	DisplayOrder int
	CreatedAt    time.Time `gorm:"<-:create"`
	UpdatedAt    time.Time
}

func (Template) TableName() string {
	return "templates"
}
