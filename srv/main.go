package main

import (
	"embed"
	"flag"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"

	modules "github.com/multios12/hmemo/modules"
)

//go:embed static/*
var static embed.FS

func main() {
	modules.Filename = *flag.String("filename", "db", "sqlite db file")
	port := flag.String("port", ":3000", "server port")
	flag.Parse()

	modules.DbInit()

	router := gin.Default()

	var f = func(c *gin.Context) {
		c.FileFromFS("static"+c.Request.URL.Path, http.FS(static))
	}

	router.GET("/", f)
	router.GET("/favicon.ico", f)
	router.GET("/:dir/:file", f)

	router.GET("/api/memos", func(c *gin.Context) {
		memos, err := modules.RowsToMemo("")
		createResponse(c, memos, err)
	})

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
			modules.UpsertMemo(b)
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
	if memos == nil {
		c.Status(http.StatusOK)
	} else {
		c.JSON(http.StatusOK, memos)
	}
}
