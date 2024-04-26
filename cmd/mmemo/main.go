package main

import (
	"embed"
	"flag"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/multios12/mmemo/pkg/diary"
	"github.com/multios12/mmemo/pkg/memo"
)

//go:embed static/*
var static embed.FS

var port string
var dataPath string

func init() {
	// 環境変数またはコマンドライン引数の読み込み
	flag.StringVar(&port, "p", ":3000", "Webサーバが使用するポートを指定します")
	flag.StringVar(&dataPath, "d", "./data", "")
	flag.Parse()

}

func main() {
	// ルーティング
	router := gin.Default()

	router.GET("/", getStatic)
	router.GET("/index.html", getStatic)
	router.GET("/favicon.ico", getStatic)

	// モジュールの初期化
	memo.Initial(router, dataPath)
	diary.Initial(router, dataPath)

	router.Run(port)
}

// スタティックリソース GET API
func getStatic(c *gin.Context) {
	p := "static" + c.Request.URL.Path
	c.FileFromFS(p, http.FS(static))
}
