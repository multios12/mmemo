package image

import (
	"fmt"
	"path"
	"regexp"
	"strconv"

	"github.com/multios12/mmemo/pkg/store"
)

var tmpImagePattern = regexp.MustCompile(`!\[[^\]]*\]\((/api/([^/\s]+)/images/tmp/([^)\s]+))(?:\s+"[^"]*")?\)`)

func createPath(prefix string) (string, error) {
	return store.NextImagePath(prefix)
}

func saveImage(imagePath string, contentType string, data []byte) error {
	return store.SaveImage(store.Image{
		Path:        imagePath,
		ContentType: contentType,
		Data:        data,
	})
}

func findImage(imagePath string) (store.Image, error) {
	return store.FindImage(imagePath)
}

func tempImageStorePath(category string, file string) string {
	return fmt.Sprintf("tmp_%s_%s", category, file)
}

func entryImagePrefix(category string, id string) string {
	idNumber, _ := strconv.Atoi(id)
	return path.Join(category, fmt.Sprintf("%05d", idNumber)) + "/"
}

func SaveTempImage(category string, contentType string, data []byte) (string, error) {
	imagePath, err := createPath(fmt.Sprintf("tmp_%s_", category))
	if err != nil {
		return "", err
	}
	if err := saveImage(imagePath, contentType, data); err != nil {
		return "", err
	}

	return fmt.Sprintf("/api/%s/images/tmp/%s", category, path.Base(imagePath)), nil
}

func SaveEntryImage(category string, id string, contentType string, data []byte) (string, error) {
	imagePath, err := createPath(entryImagePrefix(category, id))
	if err != nil {
		return "", err
	}
	if err := saveImage(imagePath, contentType, data); err != nil {
		return "", err
	}

	return fmt.Sprintf("/api/%s/%s/images/%s", category, id, path.Base(imagePath)), nil
}

func FindTempImage(category string, file string) (store.Image, error) {
	return findImage(tempImageStorePath(category, file))
}

func moveTempImage(category string, file string, destPrefix string) (string, error) {
	image, err := store.FindImage(tempImageStorePath(category, file))
	if err != nil {
		return "", err
	}

	destPath, err := createPath(destPrefix)
	if err != nil {
		return "", err
	}

	if err := store.SaveImage(store.Image{
		Path:        destPath,
		ContentType: image.ContentType,
		Data:        image.Data,
	}); err != nil {
		return "", err
	}
	if err := store.DeleteImage(tempImageStorePath(category, file)); err != nil {
		return "", err
	}

	return destPath, nil
}

func DeleteImagesByPrefix(prefix string) error {
	return store.DeleteImagesByPrefix(prefix)
}

// 一時保存画像を保存先に移動
func MoveTempImages(detail string, destPrefix string, imageTemplate string) (string, error) {
	matches := tmpImagePattern.FindAllStringSubmatch(detail, -1)
	replacements := make(map[string]string, len(matches))
	for _, subMatches := range matches {
		newPath, err := moveTempImage(subMatches[2], subMatches[3], destPrefix)
		if err != nil {
			return detail, err
		}
		replacements[subMatches[0]] = fmt.Sprintf(imageTemplate, path.Base(newPath))
	}

	return tmpImagePattern.ReplaceAllStringFunc(detail, func(match string) string {
		if replacement, ok := replacements[match]; ok {
			return replacement
		}
		return match
	}), nil
}
