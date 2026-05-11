package store

import (
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Open(dataPath string) error {
	filename := filepath.Join(dataPath, "memo.db")

	var err error
	db, err = gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		return err
	}

	return db.AutoMigrate(&Entry{}, &Image{}, &Template{})
}
