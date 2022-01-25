package main

import (
	"embed"
	"flag"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "modernc.org/sqlite"

	modules "github.com/multios12/hmemo/modules"
)

//go:embed static/*
var static embed.FS

func main() {
	// 環境変数またはコマンドライン引数の読み込み
	modules.Filename = *flag.String("filename", "db", "sqlite db file")
	port := flag.String("port", ":3000", "server port")
	flag.Parse()

	// DB初期化
	modules.DbInit()

	// ルーティング
	router := gin.Default()

	// ルーティング for 静的リソース
	var f = func(c *gin.Context) {
		c.FileFromFS("static"+c.Request.URL.Path, http.FS(static))
	}
	router.GET("/", f)
	router.GET("/favicon.ico", f)
	router.GET("/static/:dir/:file", f)

	router.GET("/api/memos", func(c *gin.Context) {
		memos, err := modules.RowsToMemo("")
		createResponse(c, memos, err)
	})

	// ルーティング for API
	router.GET("/api/memos/:id", func(c *gin.Context) {
		memos, err := modules.RowsToMemo(c.Param("id"))
		createResponse(c, memos, err)
	})

	router.PUT("/api/memos/", func(c *gin.Context) {
		var b modules.Memo
		err := c.ShouldBindJSON(&b)
		if err == nil {
			modules.UpsertMemo(b)
		}
		createResponse(c, nil, err)
	})

	router.POST("/api/memos/:id", func(c *gin.Context) {
		var b modules.Memo
		err := c.ShouldBindJSON(&b)
		if err == nil {
			validate := validator.New() //インスタンス生成
			err = validate.Struct(b)
			if err == nil {
				b.Id = c.Param("id")
				modules.UpsertMemo(b)
			}
		}
		createResponse(c, nil, err)
	})

	router.DELETE("/api/memos/:id", func(c *gin.Context) {
		err := modules.DeleteMemo(c.Param("id"))
		createResponse(c, nil, err)
	})

	router.Run(*port)
}

func createResponse(c *gin.Context, memos []modules.Memo, err error) {
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, memos)
}
