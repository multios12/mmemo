package store

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

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
