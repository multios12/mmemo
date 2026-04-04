package main

import (
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/multios12/mmemo/pkg/entryapi"
)

//go:embed static/*
var static embed.FS
var staticFiles fs.FS

var port string
var setting entryapi.SettingModel

func init() {
	flag.StringVar(&port, "p", ":3000", "Webサーバが使用するポートを指定します")
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage: %s [-p :3000]\n", os.Args[0])
		fmt.Fprintln(out, "")
		fmt.Fprintln(out, "起動ディレクトリに settings.json と memo.db を作成して利用します。")
		fmt.Fprintln(out, "")
		flag.PrintDefaults()
	}
	flag.Parse()

	var err error
	staticFiles, err = fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}

}

func main() {
	dataPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dbExists := fileExists(filepath.Join(dataPath, "memo.db"))

	// 設定ファイルの読み込み
	loadSettingJson()

	// ルーティング
	router := http.NewServeMux()

	router.HandleFunc("GET /", getStatic)
	router.HandleFunc("GET /index.html", getStatic)
	router.HandleFunc("GET /favicon.ico", getStatic)
	router.HandleFunc("GET /manifest.json", getStatic)

	// モジュールの初期化
	if err := entryapi.Initial(router, setting); err != nil {
		log.Fatal(err)
	}
	if !dbExists {
		if err := entryapi.SeedDefaultTemplates(defaultTemplateSetting()); err != nil {
			log.Printf("default template seed failed: %v", err)
		}
	}
	if err := http.ListenAndServe(port, withRecovery(withLogging(router))); err != nil {
		log.Fatal(err)
	}
}

// スタティックリソース GET API
func getStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		serveEmbeddedFile(w, "index.html")
		return
	} else if shouldServeSPA(path) {
		serveEmbeddedFile(w, "index.html")
		return
	}

	r = r.Clone(r.Context())
	r.URL.Path = path
	http.FileServer(http.FS(staticFiles)).ServeHTTP(w, r)
}

func serveEmbeddedFile(w http.ResponseWriter, name string) {
	b, err := fs.ReadFile(staticFiles, name)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if contentType := mime.TypeByExtension(filepath.Ext(name)); contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func shouldServeSPA(path string) bool {
	if strings.Contains(filepath.Base(path), ".") {
		return false
	}

	_, err := fs.Stat(staticFiles, strings.TrimPrefix(path, "/"))
	return errors.Is(err, fs.ErrNotExist)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func defaultTemplateSetting() entryapi.SettingModel {
	return entryapi.SettingModel{
		Categories: []entryapi.CategoryModel{
			{
				Key: "diary",
				Templates: []entryapi.TemplateModel{
					{Name: "通常日記", Value: "## 今日の出来事\n----\n## 明日の予定\n----"},
					{Name: "ふりかえり", Value: "## 良かったこと\n----\n## 改善したいこと\n----\n## 次にやること\n----"},
				},
			},
			{
				Key: "sample",
				Templates: []entryapi.TemplateModel{
					{Name: "基本", Value: "## データ1\n----\n## データ2\n----"},
					{Name: "打ち合わせ", Value: "## 議題\n----\n## 決定事項\n----\n## 宿題\n----"},
				},
			},
		},
	}
}

// 設定ファイルの読み込み
func loadSettingJson() {
	dataPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// データファイルパスの確認と、存在しない場合は作成
	if _, err := os.Stat(dataPath); err != nil {
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
