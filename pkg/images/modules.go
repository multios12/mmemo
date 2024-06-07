package images

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

func Rename(srcUrl string, newDirPath string) (newPath string, err error) {
	srcUrl = strings.Replace(srcUrl, "api/images/", "", 1)
	oldPath := path.Join(imagesPath, srcUrl)
	if _, err = os.Stat(oldPath); err != nil {
		return "", err
	}

	if _, err := os.Stat(newDirPath); err != nil {
		os.Mkdir(newDirPath, 0777)
	}

	newPath = createPath(newDirPath)
	if err = os.Rename(oldPath, newPath); err != nil {
		return "", err
	}
	return newPath, nil
}

func createPath(prefix string) string {
	for i := 1; i < 999; i++ {
		n := prefix + fmt.Sprintf("%03d.png", i)
		if _, err := os.Stat(n); err != nil {
			return n
		}
	}
	return ""
}

// 一時保存画像をdiaryデータパスに移動
func MoveTempImages(detail string, newDirPath string, imageTemplate string) (string, error) {
	re, _ := regexp.Compile(`(?:!\[([^[]+)\])(?:\((?:(\/api\/images\/tmp_[^()\s]+)(?:\s"((?:[^"]*\\")*[^"]*)"\s*)?)\))$`)
	matches := re.FindAllStringSubmatch(detail, 100)
	for _, subMatches := range matches {
		if newPath, err := Rename(subMatches[2], newDirPath); err != nil {
			return detail, err
		} else {
			newPath = newPath[len(newPath)-7:]
			n := fmt.Sprintf(imageTemplate, newPath)
			detail = strings.Replace(detail, subMatches[0], n, 1)
		}
	}

	return detail, nil
}
