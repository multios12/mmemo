package memo

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var setting SettingModel
var dataPath string

func Initial(router *gin.Engine, d string, s SettingModel) error {
	setting = s
	dataPath, _ = filepath.Abs(d)
	// DB初期化
	if err := dbOpen(dataPath); err != nil {
		log.Printf("memo init failed: %v", err)
		return err
	}
	memos := findMemos("")
	log.Printf("info: memo[dataPath=%s]", dataPath)
	log.Printf("info: memo[count=%d]", len(memos))

	for _, c := range setting.Categories {
		categoryPath := path.Join(dataPath, c.Key)
		if _, err := os.Stat(categoryPath); err != nil {
			if err := os.MkdirAll(categoryPath, 0755); err != nil {
				return err
			}
		}
	}

	// ルーティング
	router.GET("/api/memos", getSetting)
	router.GET("/api/memos/:category", getMemos)
	router.GET("/api/memos/:category/:id", getMemosId)
	router.PUT("/api/memos/:category", postMemos)
	router.POST("/api/memos/:category/:id", postMemos)
	router.DELETE("/api/memos/:category/:id", deleteMemosId)
	router.GET("/api/memos/:category/:id/:file", getImage)

	return nil
}

func getSetting(c *gin.Context) {
	c.JSON(http.StatusOK, setting)
}

func getMemos(c *gin.Context) {
	var memos []Memo
	if len(c.Query("month")) == 0 {
		memos = findMemos(c.Param("category"))
	} else {
		memos = findMemosByMonth(c.Param("category"), c.Query("month"))
	}
	createResponse(c, memos, nil)
}

func getMemosId(c *gin.Context) {
	memo, err := findMemoById(c.Param("category"), c.Param("id"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		createResponse(c, nil, err)
		return
	}
	c.JSON(http.StatusOK, memo)
}

func postMemos(c *gin.Context) {
	var b Memo
	err := c.ShouldBindJSON(&b)
	if err == nil {
		category := c.Param("category")
		if pathID := c.Param("id"); pathID != "" {
			id, convErr := strconv.Atoi(pathID)
			if convErr != nil {
				err = convErr
			} else if b.Id != 0 && b.Id != id {
				err = errors.New("path id and body id do not match")
			} else {
				b.Id = id
			}
		}
		if err == nil {
			b.Category = category
			isNewMemo := b.Id == 0
			if isNewMemo {
				b = upsertMemo(b)
			}
			// 一時保存画像をmemoデータパスに移動
			b.Value, err = Move(b.Value, category, strconv.Itoa(b.Id))
			if err != nil && isNewMemo {
				deleteMemo(category, strconv.Itoa(b.Id))
			}
			if err == nil {
				b = upsertMemo(b)
			}
		}
	}
	createResponse(c, nil, err)
}

func deleteMemosId(c *gin.Context) {
	deleteMemo(c.Param("category"), c.Param("id"))

	// 画像ディレクトリの削除
	n, _ := strconv.Atoi(c.Param("id"))
	id := fmt.Sprintf("%05d", n)
	dirname := path.Join(dataPath, c.Param("category"), id)
	if _, err := os.Stat(dirname); err == nil {
		os.RemoveAll(dirname)
	}

	createResponse(c, nil, nil)
}

func createResponse(c *gin.Context, memos []Memo, err error) {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, memos)
}

func getImage(c *gin.Context) {
	filename := path.Join(dataPath, c.Param("category"), c.Param("id"), c.Param("file"))
	if b, err := os.ReadFile(filename); err == nil {
		c.Data(http.StatusOK, "image/png", b)
		return
	}
	c.Status(http.StatusNotFound)
}
