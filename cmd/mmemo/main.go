package main

import (
	"embed"
	"encoding/json"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/multios12/mmemo/pkg/images"
	"github.com/multios12/mmemo/pkg/memo"
)

//go:embed static/*
var static embed.FS
var staticFiles fs.FS

var port string
var dataPath string
var setting memo.SettingModel

func init() {
	// 環境変数またはコマンドライン引数の読み込み
	flag.StringVar(&port, "p", ":3000", "Webサーバが使用するポートを指定します")
	flag.StringVar(&dataPath, "d", "./data", "")
	flag.Parse()

	var err error
	staticFiles, err = fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}

}

func main() {
	// 設定ファイルの読み込み
	loadSettingJson()

	// ルーティング
	router := http.NewServeMux()

	router.HandleFunc("GET /", getStatic)
	router.HandleFunc("GET /index.html", getStatic)
	router.HandleFunc("GET /favicon.ico", getStatic)
	router.HandleFunc("GET /manifest.json", getStatic)

	// モジュールの初期化
	if err := memo.Initial(router, dataPath, setting); err != nil {
		log.Fatal(err)
	}
	if err := images.Initial(router, dataPath); err != nil {
		log.Fatal(err)
	}
	if err := http.ListenAndServe(port, withRecovery(withLogging(router))); err != nil {
		log.Fatal(err)
	}
}

// スタティックリソース GET API
func getStatic(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		r = r.Clone(r.Context())
		r.URL.Path = "/index.html"
	}
	http.FileServer(http.FS(staticFiles)).ServeHTTP(w, r)
}

// 設定ファイルの読み込み
func loadSettingJson() {

	// データファイルパスの確認と、存在しない場合は作成
	if _, err := os.Stat((dataPath)); err != nil {
		if err := os.MkdirAll(dataPath, 0755); err != nil {
			panic(err)
		}
	}

	// 設定ファイルの読み込み、存在しない場合はサンプルファイルをもとに作成
	filename := filepath.Join(dataPath, "settings.json")
	if _, err := os.Stat(filename); err != nil {
		b, err := static.ReadFile("static/.default.settings.json")

		if err != nil {
			panic(err)
		}
		if err := os.WriteFile(filename, b, os.ModePerm); err != nil {
			panic(err)
		}
	}

	if b, err := os.ReadFile(filename); err != nil {
		panic(err)
	} else if err := json.Unmarshal(b, &setting); err != nil {
		panic(err)
	}
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("panic: %v", recovered)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
