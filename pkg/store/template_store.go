package store

import (
	"errors"

	"gorm.io/gorm"
)

func FindTemplates(category string) ([]Template, error) {
	var templates []Template
	result := db.
		Where("category = ?", category).
		Order("display_order asc, id asc").
		Find(&templates)
	return templates, result.Error
}

func SaveTemplate(category string, template Template) (Template, error) {
	var current Template
	result := db.Where("category = ? and name = ?", category, template.Name).First(&current)
	if result.Error == nil {
		current.Value = template.Value
		current.Name = template.Name
		return current, db.Save(&current).Error
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return Template{}, result.Error
	}

	var maxOrder int
	db.Model(&Template{}).
		Where("category = ?", category).
		Select("coalesce(max(display_order), 0)").
		Scan(&maxOrder)

	template.Category = category
	template.DisplayOrder = maxOrder + 1
	if err := db.Create(&template).Error; err != nil {
		return Template{}, err
	}
	return template, nil
}

func SeedTemplates(category string, templates []Template) error {
	existing, err := FindTemplates(category)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return nil
	}

	for index, template := range templates {
		template.Category = category
		template.DisplayOrder = index + 1
		if err := db.Create(&template).Error; err != nil {
			return err
		}
	}
	return nil
}
