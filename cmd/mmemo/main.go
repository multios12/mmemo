package main

import (
	"embed"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/multios12/mmemo/pkg/diary"
	"github.com/multios12/mmemo/pkg/images"
	"github.com/multios12/mmemo/pkg/memo"
)

//go:embed static/*
var static embed.FS

var port string
var dataPath string
var setting memo.SettingModel

const DEFAULT_SETTING_JSON string = "{\"Diary\": {\"Name\": \"日記\"},\"Categories\": [{\"Key\": \"sample\",\"Name\": \"サンプル\",\"UseDate\": true,\"UseTag\": true,\"Template\": \"## データ１\n----\n## データ2\n----\",\"Titles\": {\"Name\": \"名前\",\"Date\": \"日付\",\"Tags\": \"タグ\",\"Value\": \"情報\"}}]}"

func init() {
	// 環境変数またはコマンドライン引数の読み込み
	flag.StringVar(&port, "p", ":3000", "Webサーバが使用するポートを指定します")
	flag.StringVar(&dataPath, "d", "./data", "")
	flag.Parse()

}

func main() {
	// 設定ファイルの読み込み
	loadSettingJson()

	// ルーティング
	router := gin.Default()

	router.GET("/", getStatic)
	router.GET("/index.html", getStatic)
	router.GET("/favicon.ico", getStatic)

	// モジュールの初期化
	memo.Initial(router, dataPath, setting)
	diary.Initial(router, dataPath)
	images.Initial(router, dataPath)
	router.Run(port)
}

// スタティックリソース GET API
func getStatic(c *gin.Context) {
	p := "static" + c.Request.URL.Path
	c.FileFromFS(p, http.FS(static))
}

// 設定ファイルの読み込み
func loadSettingJson() {

	// データファイルパスの確認と、存在しない場合は作成
	if _, err := os.Stat((dataPath)); err != nil {
		os.Mkdir(dataPath, os.ModeDir)
	}

	// 設定ファイルの読み込み、存在しない場合はサンプルファイルをもとに作成
	filename := filepath.Join(dataPath, "settings.json")
	if _, err := os.Stat(filename); err != nil {
		b, err := static.ReadFile("static/.default.settings.json")

		if err != nil {
			panic(err)
		}
		os.WriteFile(filename, b, os.ModePerm)
	}

	if b, err := os.ReadFile(filename); err != nil {
		panic(err)
	} else if err := json.Unmarshal(b, &setting); err != nil {
		panic(err)
	}
}
