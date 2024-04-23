package main

import (
	"embed"
	"flag"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"

	"github.com/multios12/mmemo/pkg/diary"
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

	// モジュールの初期化
	println("test")
	//memo.Initial(router, dataPath)
	diary.Initial(router, dataPath)

	router.Run(port)
}
