package memo

import (
	"fmt"
	"path"
	"strconv"

	"github.com/multios12/mmemo/pkg/images"
)

// 一時保存画像をdiaryデータパスに移動
func Move(value string, category string, id string) (string, error) {
	idNumber, _ := strconv.Atoi(id)
	id = fmt.Sprintf("%05d", idNumber)
	newDirPath := path.Join(dataPath, category, id) + "/"
	imageTemplate := fmt.Sprintf("![イメージ](/api/memos/%s/%s", category, id) + "/%s)"
	d, err := images.MoveTempImages(value, newDirPath, imageTemplate)
	if err != nil {
		return d, err
	}
	return d, nil
}
