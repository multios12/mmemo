package diary

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/multios12/mmemo/pkg/images"
)

// gin.Engineインスタンスにルーティングを設定して返す
func Initial(router *gin.Engine, dataPath string) {
	diaryPath, _ = filepath.Abs(dataPath)
	diaryPath = filepath.Join(diaryPath, "diary")
	if _, err := os.Stat(diaryPath); err != nil {
		os.Mkdir(diaryPath, 0777)
	}
	router.GET("/api/diary/:year/:month", getMonth)
	router.GET("/api/diary/:year/:month/:day", getDetail)
	router.POST("/api/diary/:year/:month/:day", postDetail)
	router.DELETE("/api/diary/:year/:month/:day", deleteDetail)
	router.GET("/api/diary/:year/:month/:day/images/:file", getImage)
}

func getMonth(c *gin.Context) {
	month := c.Param("year") + c.Param("month")
	var m = readListFile(month)
	m.WritedMonths = getWritedMonths()
	c.JSON(200, m)
}

func getDetail(c *gin.Context) {
	day := c.Param("year") + "-" + c.Param("month") + "-" + c.Param("day")
	if l := readDetail(day); len(l.Day) == 0 {
		c.Status(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, l)
	}
}

func postDetail(c *gin.Context) {
	day := c.Param("year") + "-" + c.Param("month") + "-" + c.Param("day")

	var detail detailModel
	if err := c.ShouldBindJSON(&detail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if detail.Outline = strings.TrimSpace(detail.Outline); detail.Outline == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "outline is not found."})
	}
	detail.Day = day

	// 一時保存画像をdiaryデータパスに移動
	newDirPath := c.Param("year") + c.Param("month") + c.Param("day")
	newDirPath = path.Join(diaryPath, newDirPath) + "/"
	imageTemplate := "![イメージ](/api/diary/" + c.Param("year") + "/" + c.Param("month") + "/" + c.Param("day") + "/images/%s)"
	d, err := images.MoveTempImages(detail.Detail, newDirPath, imageTemplate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	detail.Detail = d

	if err = detail.writeDetailFile(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func deleteDetail(c *gin.Context) {
	day := c.Param("year") + "-" + c.Param("month") + "-" + c.Param("day")
	if err := removeDetail(day); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.Status(http.StatusOK)
	}
}

func getImage(c *gin.Context) {
	filename := c.Param("year") + c.Param("month") + c.Param("day")
	filename = path.Join(diaryPath, filename, c.Param("file"))
	if b, err := os.ReadFile(filename); err == nil {
		c.Data(http.StatusOK, "image/png", b)
	}
	c.Status(http.StatusNotFound)
}
