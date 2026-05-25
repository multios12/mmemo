package web

import (
	"fmt"
	"path"
	"strconv"

	entryimage "github.com/multios12/mmemo/pkg/image"
)

// 一時保存画像をdiaryデータパスに移動
func Move(value string, category string, id string) (string, error) {
	displayID := id
	idNumber, _ := strconv.Atoi(id)
	storageID := fmt.Sprintf("%05d", idNumber)
	newDirPath := path.Join(category, storageID) + "/"
	imageTemplate := fmt.Sprintf("![イメージ](./%s/%s/%%s)", category, displayID)
	d, err := entryimage.MoveTempImages(value, newDirPath, imageTemplate)
	if err != nil {
		return d, err
	}
	return d, nil
}
