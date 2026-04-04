package images

import (
	"fmt"
	"path"
	"regexp"

	"github.com/multios12/mmemo/pkg/store"
)

var tmpImagePattern = regexp.MustCompile(`!\[[^\]]*\]\((/images/tmp_[^)\s]+)(?:\s+"[^"]*")?\)`)

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

func moveImage(srcURL string, destPrefix string) (string, error) {
	srcPath := srcURL[len("/images/"):]
	image, err := store.FindImage(srcPath)
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
	if err := store.DeleteImage(srcPath); err != nil {
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
		newPath, err := moveImage(subMatches[1], destPrefix)
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
