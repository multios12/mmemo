package images

import (
	"fmt"
	"io"
	"net/http"

	"gorm.io/gorm"
)

func Initial(router *http.ServeMux, _ string) error {
	// ルーティング
	router.HandleFunc("POST /images", postImage)
	router.HandleFunc("GET /images/{file}", getImage)

	return nil
}

// 一時保存
func postImage(w http.ResponseWriter, r *http.Request) {
	filename, err := createPath("tmp_")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	inFile, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer inFile.Close()

	data, err := io.ReadAll(inFile)
	if err != nil {
		err = fmt.Errorf("ファイルが保存できません: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/png"
	}
	contentType, data, err = normalizeImageForStorage(contentType, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := saveImage(filename, contentType, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename = fmt.Sprintf("/images/%s", filename)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(filename))
}

// 一時画像の取得
func getImage(w http.ResponseWriter, r *http.Request) {
	image, err := findImage(r.PathValue("file"))
	if err == nil {
		contentType := image.ContentType
		if contentType == "" {
			contentType = "image/png"
		}
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(image.Data)
		return
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}
