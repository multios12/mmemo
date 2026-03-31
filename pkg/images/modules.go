package images

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
)

var tmpImagePattern = regexp.MustCompile(`!\[[^\]]*\]\((/api/images/tmp_[^)\s]+)(?:\s+"[^"]*")?\)`)

func Rename(srcUrl string, newDirPath string) (newPath string, err error) {
	srcUrl = srcUrl[len("/api/images/"):]
	oldPath := path.Join(imagesPath, srcUrl)
	if _, err = os.Stat(oldPath); err != nil {
		return "", err
	}

	if err := os.MkdirAll(newDirPath, 0755); err != nil {
		return "", err
	}

	newPath, err = createPath(newDirPath)
	if err != nil {
		return "", err
	}
	if err = os.Rename(oldPath, newPath); err != nil {
		return "", err
	}
	return newPath, nil
}

func createPath(prefix string) (string, error) {
	for i := 1; i < 999; i++ {
		n := prefix + fmt.Sprintf("%03d.png", i)
		if _, err := os.Stat(n); err != nil {
			return n, nil
		}
	}
	return "", errors.New("image path limit reached")
}

// 一時保存画像をdiaryデータパスに移動
func MoveTempImages(detail string, newDirPath string, imageTemplate string) (string, error) {
	matches := tmpImagePattern.FindAllStringSubmatch(detail, -1)
	replacements := make(map[string]string, len(matches))
	for _, subMatches := range matches {
		if newPath, err := Rename(subMatches[1], newDirPath); err != nil {
			return detail, err
		} else {
			replacements[subMatches[0]] = fmt.Sprintf(imageTemplate, newPath[len(newPath)-7:])
		}
	}

	return tmpImagePattern.ReplaceAllStringFunc(detail, func(match string) string {
		if replacement, ok := replacements[match]; ok {
			return replacement
		}
		return match
	}), nil
}
