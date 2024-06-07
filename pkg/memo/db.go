package memo

import (
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

func dbOpen(dataPath string) (err error) {
	filename := filepath.Join(dataPath, "memo.db")
	db, err = gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&Memo{})
	return err
}

func findMemos(category string) (memos []Memo) {
	if len(category) == 0 {
		db.Find(&memos)
	} else {
		db.Where("category = ?", category).Find(&memos)
	}
	return memos
}

func findMemoById(category string, id string) (memos []Memo) {
	db.Where("category = ? and id = ?", category, id).Find(&memos)
	return memos
}

func findMemosByMonth(category string, month string) (memos []Memo) {
	from := month + "-01"
	toDate, _ := time.Parse("2006-01-02", from)
	toDate = toDate.AddDate(0, 1, 0).AddDate(0, 0, -1)
	to := toDate.Format("2006-01-02")
	db.Where("category = ? and date >= ? and date <= ?", category, from, to).Find(&memos)
	return memos
}

func upsertMemo(m Memo) Memo {
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&m)
	return m
}

func deleteMemo(id string) {
	db.Delete(&Memo{}, id)
}
