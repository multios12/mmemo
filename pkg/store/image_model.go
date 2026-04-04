package store

type Image struct {
	Path        string `gorm:"primaryKey"`
	ContentType string
	Data        []byte
}

func (Image) TableName() string {
	return "images"
}
