package memo

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var setting SettingModel

func Initial(router *gin.Engine, dataPath string, s SettingModel) error {
	setting = s
	// DB初期化
	if err := dbOpen(dataPath); err != nil {
		fmt.Sprintln(err)
		return err
	}
	memos := findMemos("")
	log.Println("info: memo[dataPath=" + filepath.Join(dataPath, "memo") + "]")
	log.Println("info: memo[count=" + strconv.Itoa(len(memos)) + "]")

	// ルーティング
	router.GET("/api/memos", getSetting)
	router.GET("/api/memos/:category", getMemos)
	router.GET("/api/memos/:category/:id", getMemosId)
	router.PUT("/api/memos/:category", putMemos)
	router.POST("/api/memos/:category/:id", postMemosId)
	router.DELETE("/api/memos/:category/:id", deleteMemosId)

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

func putMemos(c *gin.Context) {
	var b Memo
	err := c.ShouldBindJSON(&b)
	if err == nil {
		b.Category = c.Param("category")
		upsertMemo(b)
	}
	createResponse(c, nil, err)
}

func postMemosId(c *gin.Context) {
	var b Memo
	err := c.ShouldBindJSON(&b)
	if err == nil {
		validate := validator.New() //インスタンス生成
		err = validate.Struct(b)
		b.Category = c.Param("category")
		if err == nil {
			b.Id, _ = strconv.Atoi(c.Param("id"))
			upsertMemo(b)
		}
	}
	createResponse(c, nil, err)
}

func deleteMemosId(c *gin.Context) {
	deleteMemo(c.Param("id"))
	createResponse(c, nil, nil)
}

func createResponse(c *gin.Context, memos []Memo, err error) {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, memos)
}
