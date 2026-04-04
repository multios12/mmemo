package store

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var db *gorm.DB

func Open(dataPath string) error {
	filename := filepath.Join(dataPath, "memo.db")

	var err error
	db, err = gorm.Open(sqlite.Open(filename), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := migrateSchema(); err != nil {
		return err
	}

	return db.AutoMigrate(&Entry{}, &Image{})
}

func migrateSchema() error {
	migrator := db.Migrator()

	if migrator.HasTable("memos") && !migrator.HasTable("entries") {
		if err := migrator.RenameTable("memos", "entries"); err != nil {
			return fmt.Errorf("rename table memos -> entries: %w", err)
		}
	}

	if migrator.HasTable("entries") && migrator.HasColumn("entries", "name") && !migrator.HasColumn("entries", "title") {
		if err := migrator.RenameColumn("entries", "name", "title"); err != nil {
			return fmt.Errorf("rename column entries.name -> title: %w", err)
		}
	}

	return nil
}

func FindEntries(category string) (entries []Entry) {
	if category == "" {
		db.Find(&entries)
	} else {
		db.Where("category = ?", category).Find(&entries)
	}
	return entries
}

func FindEntryByID(category string, id string) (Entry, error) {
	var entry Entry
	result := db.Where("category = ? and id = ?", category, id).First(&entry)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return Entry{}, gorm.ErrRecordNotFound
	}
	return entry, result.Error
}

func FindEntriesByMonth(category string, month string) (entries []Entry, err error) {
	from := month + "-01"
	toDate, err := time.Parse("2006-01-02", from)
	if err != nil {
		return nil, err
	}
	to := toDate.AddDate(0, 1, -1).Format("2006-01-02")

	result := db.Where("category = ? and date >= ? and date <= ?", category, from, to).Find(&entries)
	return entries, result.Error
}

func FindEntryByDate(category string, day string) (Entry, error) {
	var entry Entry
	result := db.Where("category = ? and date = ?", category, day).First(&entry)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return Entry{}, gorm.ErrRecordNotFound
	}
	return entry, result.Error
}

func FindEntriesByDateRange(category string, month string) ([]Entry, error) {
	from := month + "-01"
	toDate, err := time.Parse("2006-01-02", from)
	if err != nil {
		return nil, err
	}
	to := toDate.AddDate(0, 1, -1).Format("2006-01-02")

	var entries []Entry
	result := db.Where("category = ? and date >= ? and date <= ?", category, from, to).Order("date desc").Find(&entries)
	return entries, result.Error
}

func FindMonths(category string) ([]string, error) {
	type result struct {
		Month string
	}

	var rows []result
	err := db.Model(&Entry{}).
		Select("substr(date, 1, 7) as month").
		Where("category = ?", category).
		Group("substr(date, 1, 7)").
		Order("month desc").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	months := make([]string, 0, len(rows))
	for _, row := range rows {
		if row.Month != "" {
			months = append(months, row.Month)
		}
	}
	return months, nil
}

func UpsertEntry(entry Entry) Entry {
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&entry)
	return entry
}

func UpsertEntryByDate(category string, entry Entry) error {
	var current Entry
	result := db.Where("category = ? and date = ?", category, entry.Date).First(&current)
	if result.Error == nil {
		entry.Id = current.Id
		return db.Save(&entry).Error
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}
	return db.Create(&entry).Error
}

func DeleteEntry(category string, id string) {
	db.Where("category = ? and id = ?", category, id).Delete(&Entry{})
}

func DeleteEntryByDate(category string, day string) error {
	return db.Where("category = ? and date = ?", category, day).Delete(&Entry{}).Error
}

func SaveImage(image Image) error {
	return db.Save(&image).Error
}

func FindImage(path string) (Image, error) {
	var image Image
	result := db.Where("path = ?", path).First(&image)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return Image{}, gorm.ErrRecordNotFound
	}
	return image, result.Error
}

func DeleteImage(path string) error {
	return db.Where("path = ?", path).Delete(&Image{}).Error
}

func DeleteImagesByPrefix(prefix string) error {
	return db.Where("path LIKE ?", prefix+"%").Delete(&Image{}).Error
}

func NextImagePath(prefix string) (string, error) {
	var images []Image
	if err := db.Select("path").Where("path LIKE ?", prefix+"%").Find(&images).Error; err != nil {
		return "", err
	}

	exists := make(map[string]struct{}, len(images))
	for _, image := range images {
		exists[image.Path] = struct{}{}
	}

	for i := 1; i < 999; i++ {
		candidate := fmt.Sprintf("%s%03d.png", prefix, i)
		if _, ok := exists[candidate]; !ok {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("image path limit reached: %s", strings.TrimSuffix(prefix, "/"))
}
