package diary

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
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
	router.POST("/api/diary/images", postImage)
	router.GET("/api/diary/images/:file", getImage)
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
	} else {
		detail.Day = day
		if err = detail.writeDetailFile(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.Status(http.StatusOK)
	}
}

func deleteDetail(c *gin.Context) {
	day := c.Param("year") + "-" + c.Param("month") + "-" + c.Param("day")
	if err := removeDetail(day); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.Status(http.StatusOK)
	}
}

func postImage(c *gin.Context) {
	filename := path.Join(diaryPath, "_")

	for i := 0; i < 99999999; i++ {
		n := filename + fmt.Sprintf("%08d.png", i)
		if _, err := os.Stat(n); err != nil {
			filename = n
			break
		}
	}

	inFile, _, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

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
	filename = path.Join("./api/diary/images", filename)
	c.String(http.StatusOK, filename)
}

func getImage(c *gin.Context) {
	filename := path.Join(diaryPath, c.Param("file"))
	if b, err := os.ReadFile(filename); err == nil {
		c.Data(http.StatusOK, "image/png", b)
	}
	c.Status(http.StatusNotFound)
}
