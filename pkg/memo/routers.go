package memo

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/multios12/mmemo/pkg/images"
	"github.com/multios12/mmemo/pkg/store"
	"gorm.io/gorm"
)

var setting SettingModel
var dataPath string

func Initial(router *http.ServeMux, d string, s SettingModel) error {
	setting = s
	dataPath, _ = filepath.Abs(d)
	if err := store.Open(dataPath); err != nil {
		log.Printf("memo init failed: %v", err)
		return err
	}
	entries := store.FindEntries("")
	log.Printf("info: memo[dataPath=%s]", dataPath)
	log.Printf("info: memo[count=%d]", len(entries))

	router.HandleFunc("GET /settings", getSetting)
	router.HandleFunc("GET /api/{category}", getEntries)
	router.HandleFunc("GET /api/{category}/{id}", getEntry)
	router.HandleFunc("PUT /api/{category}", putEntry)
	router.HandleFunc("POST /api/{category}/{id}", postEntry)
	router.HandleFunc("DELETE /api/{category}/{id}", deleteEntry)
	router.HandleFunc("GET /api/{category}/{id}/images/{file}", getImage)

	return nil
}

func getSetting(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, setting)
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
	_ = images.DeleteImagesByPrefix(path.Join(category, paddedID) + "/")
	if category == "diary" && legacyDiaryDate != "" {
		_ = images.DeleteImagesByPrefix(path.Join(category, strings.ReplaceAll(legacyDiaryDate, "-", "")) + "/")
	}
	w.WriteHeader(http.StatusOK)
}

func getImage(w http.ResponseWriter, r *http.Request) {
	imagePath := imagePathForRoute(r.PathValue("category"), r.PathValue("id"), r.PathValue("file"))
	image, err := store.FindImage(imagePath)
	if errors.Is(err, gorm.ErrRecordNotFound) && r.PathValue("category") == "diary" {
		if legacyPath, legacyErr := legacyDiaryImagePath(r.PathValue("id"), r.PathValue("file")); legacyErr == nil {
			image, err = store.FindImage(legacyPath)
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

	contentType := image.ContentType
	if contentType == "" {
		contentType = "image/png"
	}
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(image.Data)
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
