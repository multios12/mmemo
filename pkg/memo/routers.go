package memo

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

var setting SettingModel
var dataPath string

func Initial(router *gin.Engine, d string, s SettingModel) error {
	setting = s
	dataPath, _ = filepath.Abs(d)
	// DB初期化
	if err := dbOpen(dataPath); err != nil {
		fmt.Sprintln(err)
		return err
	}
	memos := findMemos("")
	log.Printf("info: memo[dataPath=%s]", dataPath)
	log.Printf("info: memo[count=%d]", len(memos))

	for _, c := range setting.Categories {
		categoryPath := path.Join(dataPath, c.Key)
		if _, err := os.Stat(categoryPath); err != nil {
			os.Mkdir(categoryPath, 0777)
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
	memos := findMemoById(c.Param("category"), c.Param("id"))
	createResponse(c, memos, nil)
}

func postMemos(c *gin.Context) {
	var b Memo
	err := c.ShouldBindJSON(&b)
	if err == nil {
		if b.Id == 0 {
			b.Category = c.Param("category")
			b = upsertMemo(b)
		}
		// 一時保存画像をdiaryデータパスに移動
		b.Value, err = Move(b.Value, c.Param("category"), strconv.Itoa(b.Id))
		b = upsertMemo(b)
	}
	createResponse(c, nil, err)
}

func deleteMemosId(c *gin.Context) {
	deleteMemo(c.Param("id"))

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
	}
	c.Status(http.StatusNotFound)
}
