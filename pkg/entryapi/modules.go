package entryapi

import (
	"fmt"
	"path"
	"strconv"

	entryimage "github.com/multios12/mmemo/pkg/image"
)

// 一時保存画像をdiaryデータパスに移動
func Move(value string, category string, id string) (string, error) {
	idNumber, _ := strconv.Atoi(id)
	id = fmt.Sprintf("%05d", idNumber)
	newDirPath := path.Join(category, id) + "/"
	imageTemplate := fmt.Sprintf("![イメージ](/api/%s/%s/images/%%s)", category, id)
	d, err := entryimage.MoveTempImages(value, newDirPath, imageTemplate)
	if err != nil {
		return d, err
	}
	return d, nil
}
