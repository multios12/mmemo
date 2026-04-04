package entryapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	entryimage "github.com/multios12/mmemo/pkg/image"
	"github.com/multios12/mmemo/pkg/store"
	"gorm.io/gorm"
)

var setting SettingModel

func Initial(router *http.ServeMux, s SettingModel) error {
	setting = s
	if err := store.Open("."); err != nil {
		log.Printf("entry init failed: %v", err)
		return err
	}
	entries := store.FindEntries("")
	log.Printf("info: entry[dataPath=.]")
	log.Printf("info: entry[count=%d]", len(entries))

	router.HandleFunc("GET /settings", getSetting)
	router.HandleFunc("GET /api/{category}", getEntries)
	router.HandleFunc("POST /api/{category}/images/tmp", postTempImage)
	router.HandleFunc("GET /api/{category}/images/tmp/{file}", getTempImage)
	router.HandleFunc("GET /api/{category}/{id}", getEntry)
	router.HandleFunc("PUT /api/{category}", putEntry)
	router.HandleFunc("POST /api/{category}/{id}", postEntry)
	router.HandleFunc("DELETE /api/{category}/{id}", deleteEntry)
	router.HandleFunc("POST /api/{category}/{id}/images", postEntryImage)
	router.HandleFunc("GET /api/{category}/{id}/images/{file}", getImage)

	return nil
}

func getSetting(w http.ResponseWriter, _ *http.Request) {
	latestSetting, err := loadSettingFromFile()
	if err == nil {
		setting = latestSetting
	}
	writeJSON(w, http.StatusOK, setting)
}

func loadSettingFromFile() (SettingModel, error) {
	b, err := os.ReadFile("settings.json")
	if err != nil {
		return SettingModel{}, err
	}

	var latest SettingModel
	if err := json.Unmarshal(b, &latest); err != nil {
		return SettingModel{}, err
	}
	return latest, nil
}

func getEntries(w http.ResponseWriter, r *http.Request) {
	category := r.PathValue("category")
	if r.URL.Query().Get("months") == "1" {
		months, err := store.FindMonths(category)
		if err != nil {
			writeErrorJSON(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, months)
		return
	}

	var entries []store.Entry
	var err error
	if month := r.URL.Query().Get("month"); month != "" {
		entries, err = store.FindEntriesByMonth(category, month)
	} else {
		entries = store.FindEntries(category)
	}
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	response := make([]entryPayload, 0, len(entries))
	for _, entry := range entries {
		response = append(response, entryToPayload(entry))
	}
	writeJSON(w, http.StatusOK, response)
}

func getEntry(w http.ResponseWriter, r *http.Request) {
	entry, err := findEntry(r.PathValue("category"), r.PathValue("id"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, entryToPayload(entry))
}

func putEntry(w http.ResponseWriter, r *http.Request) {
	saveEntry(w, r, "")
}

func postEntry(w http.ResponseWriter, r *http.Request) {
	saveEntry(w, r, r.PathValue("id"))
}

func saveEntry(w http.ResponseWriter, r *http.Request, pathID string) {
	category := r.PathValue("category")

	var payload entryPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	entry := payloadToEntry(category, payload)
	if err := validatePathID(&entry, pathID); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}
	if err := validateEntryDateDuplication(category, entry); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	isNewEntry := entry.Id == 0
	if isNewEntry {
		entry = store.UpsertEntry(entry)
	}

	value, err := Move(entry.Value, category, strconv.Itoa(entry.Id))
	if err != nil {
		if isNewEntry {
			store.DeleteEntry(category, strconv.Itoa(entry.Id))
		}
		writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	entry.Value = value
	entry = store.UpsertEntry(entry)
	w.WriteHeader(http.StatusOK)
}

func validateEntryDateDuplication(category string, entry store.Entry) error {
	for _, categorySetting := range setting.Categories {
		if categorySetting.Key != category {
			continue
		}
		if categorySetting.AllowMultipleEntriesPerDate || entry.Date == "" {
			return nil
		}

		current, err := store.FindEntryByDate(category, entry.Date)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
		if entry.Id != 0 && current.Id == entry.Id {
			return nil
		}
		return errors.New("同じ日付のエントリは登録できません")
	}
	return nil
}

func validatePathID(entry *store.Entry, pathID string) error {
	if pathID == "" {
		return nil
	}

	id, err := strconv.Atoi(pathID)
	if err != nil {
		return err
	}
	if entry.Id != 0 && entry.Id != id {
		return errors.New("path id and body id do not match")
	}
	entry.Id = id
	return nil
}

func findEntry(category string, id string) (store.Entry, error) {
	return store.FindEntryByID(category, id)
}

func deleteEntry(w http.ResponseWriter, r *http.Request) {
	category := r.PathValue("category")
	id := r.PathValue("id")
	var legacyDiaryDate string

	if category == "diary" {
		entry, err := store.FindEntryByID(category, id)
		if err == nil {
			legacyDiaryDate = entry.Date
		}
	}

	store.DeleteEntry(category, id)
	n, _ := strconv.Atoi(id)
	paddedID := fmt.Sprintf("%05d", n)
	_ = entryimage.DeleteImagesByPrefix(path.Join(category, paddedID) + "/")
	if category == "diary" && legacyDiaryDate != "" {
		_ = entryimage.DeleteImagesByPrefix(path.Join(category, strings.ReplaceAll(legacyDiaryDate, "-", "")) + "/")
	}
	w.WriteHeader(http.StatusOK)
}

func getImage(w http.ResponseWriter, r *http.Request) {
	imagePath := imagePathForRoute(r.PathValue("category"), r.PathValue("id"), r.PathValue("file"))
	storedImage, err := store.FindImage(imagePath)
	if errors.Is(err, gorm.ErrRecordNotFound) && r.PathValue("category") == "diary" {
		if legacyPath, legacyErr := legacyDiaryImagePath(r.PathValue("id"), r.PathValue("file")); legacyErr == nil {
			storedImage, err = store.FindImage(legacyPath)
		}
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	contentType := storedImage.ContentType
	if contentType == "" {
		contentType = "image/png"
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(storedImage.Data)
}

func postTempImage(w http.ResponseWriter, r *http.Request) {
	imageURL, err := saveUploadedImage(r, func(contentType string, data []byte) (string, error) {
		return entryimage.SaveTempImage(r.PathValue("category"), contentType, data)
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(imageURL))
}

func getTempImage(w http.ResponseWriter, r *http.Request) {
	storedImage, err := entryimage.FindTempImage(r.PathValue("category"), r.PathValue("file"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	contentType := storedImage.ContentType
	if contentType == "" {
		contentType = "image/png"
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(storedImage.Data)
}

func postEntryImage(w http.ResponseWriter, r *http.Request) {
	imageURL, err := saveUploadedImage(r, func(contentType string, data []byte) (string, error) {
		return entryimage.SaveEntryImage(r.PathValue("category"), r.PathValue("id"), contentType, data)
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(imageURL))
}

func imagePathForRoute(category string, id string, file string) string {
	n, err := strconv.Atoi(id)
	if err == nil {
		return path.Join(category, fmt.Sprintf("%05d", n), file)
	}
	return path.Join(category, id, file)
}

func legacyDiaryImagePath(id string, file string) (string, error) {
	entry, err := store.FindEntryByID("diary", id)
	if err != nil {
		return "", err
	}
	return path.Join("diary", strings.ReplaceAll(entry.Date, "-", ""), file), nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErrorJSON(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func saveUploadedImage(
	r *http.Request,
	save func(contentType string, data []byte) (string, error),
) (string, error) {
	inFile, header, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer inFile.Close()

	data, err := io.ReadAll(inFile)
	if err != nil {
		return "", fmt.Errorf("ファイルが保存できません: %w", err)
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/png"
	}
	contentType, data, err = entryimage.NormalizeImageForStorage(contentType, data)
	if err != nil {
		return "", err
	}

	return save(contentType, data)
}
