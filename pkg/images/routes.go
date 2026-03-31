package images

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var imagesPath string

func Initial(router *gin.Engine, dataPath string) error {
	imagesPath = path.Join(dataPath, "images")

	if _, err := os.Stat((imagesPath)); err != nil {
		if err := os.MkdirAll(imagesPath, 0755); err != nil {
			return err
		}
	}

	// ルーティング
	router.POST("/api/images", postImage)
	router.GET("/api/images/:file", getImage)

	return nil
}

// 一時保存
func postImage(c *gin.Context) {
	filename := path.Join(imagesPath, "tmp_")
	var err error
	filename, err = createPath(filename)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	inFile, _, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	defer inFile.Close()

	outFile, err := os.Create(filename)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, inFile)
	if err != nil {
		err = fmt.Errorf("ファイルが保存できません: %w", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	filename = filepath.Base(filename)
	filename = fmt.Sprintf("/api/images/%s", filename)
	c.String(http.StatusOK, filename)
}

// 一時画像の取得
func getImage(c *gin.Context) {
	filename := path.Join(imagesPath, c.Param("file"))
	if b, err := os.ReadFile(filename); err == nil {
		c.Data(http.StatusOK, "image/png", b)
		return
	}
	c.Status(http.StatusNotFound)
}
