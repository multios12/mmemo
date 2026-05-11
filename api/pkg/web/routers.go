package web

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

// Initial は API ルートを登録し、設定とストアを初期化します。
func Initial(router *http.ServeMux, s SettingModel) error {
	if err := entryStore.Open("."); err != nil {
		log.Printf("entry init failed: %v", err)
		return err
	}
	setting = hydrateSetting(s)
	entries := entryStore.FindEntries("")
	log.Printf("info: entry[dataPath=.]")
	log.Printf("info: entry[count=%d]", len(entries))

	router.HandleFunc("GET    /api/settings", getSetting)
	router.HandleFunc("GET    /api/{category}", getEntries)
	router.HandleFunc("POST   /api/{category}/templates", saveTemplate)
	router.HandleFunc("DELETE /api/{category}/templates/{name}", deleteTemplate)
	router.HandleFunc("POST   /api/{category}/images/tmp", postTempImage)
	router.HandleFunc("GET    /api/{category}/images/tmp/{file}", getTempImage)
	router.HandleFunc("GET    /api/{category}/{id}", getEntry)
	router.HandleFunc("PUT    /api/{category}", putEntry)
	router.HandleFunc("POST   /api/{category}/{id}", postEntry)
	router.HandleFunc("DELETE /api/{category}/{id}", deleteEntry)
	router.HandleFunc("POST   /api/{category}/{id}/images", postEntryImage)
	router.HandleFunc("GET    /api/{category}/{id}/images/{file}", getImage)
	router.HandleFunc("DELETE /api/{category}/{id}/images/{file}", deleteImage)

	return nil
}

func getSetting(w http.ResponseWriter, _ *http.Request) {
	latestSetting, err := loadSettingFromFile()
	if err == nil {
		setting = hydrateSetting(latestSetting)
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
	normalizeSetting(&latest)
	return latest, nil
}

func normalizeSetting(target *SettingModel) {
	for categoryIndex := range target.Categories {
		category := &target.Categories[categoryIndex]
		for templateIndex := range category.Templates {
			template := &category.Templates[templateIndex]
			template.Name = strings.TrimSpace(template.Name)
			if template.Name == "" {
				template.Name = fmt.Sprintf("テンプレート%d", templateIndex+1)
			}
		}
	}
}

func hydrateSetting(base SettingModel) SettingModel {
	normalizeSetting(&base)
	for categoryIndex := range base.Categories {
		category := &base.Categories[categoryIndex]
		if err := entryStore.SeedTemplates(category.Key, templateModelsToStore(category.Key, category.Templates)); err != nil {
			log.Printf("template seed failed[category=%s]: %v", category.Key, err)
		}

		templates, err := entryStore.FindTemplates(category.Key)
		if err != nil {
			log.Printf("template load failed[category=%s]: %v", category.Key, err)
			continue
		}
		category.Templates = storeTemplatesToModels(templates)
	}
	return base
}

func templateModelsToStore(category string, templates []TemplateModel) []store.Template {
	items := make([]store.Template, 0, len(templates))
	for index, template := range templates {
		items = append(items, store.Template{
			Category:     category,
			Name:         template.Name,
			Value:        template.Value,
			Tags:         strings.Join(template.Tags, "#"),
			DisplayOrder: index + 1,
		})
	}
	return items
}

func storeTemplatesToModels(templates []store.Template) []TemplateModel {
	items := make([]TemplateModel, 0, len(templates))
	for _, template := range templates {
		items = append(items, TemplateModel{
			Name:  template.Name,
			Value: template.Value,
			Tags:  splitTags(template.Tags),
		})
	}
	return items
}

// SeedDefaultTemplates は初期テンプレートをストアへ登録します。
func SeedDefaultTemplates(defaults SettingModel) error {
	normalizeSetting(&defaults)
	for _, category := range defaults.Categories {
		if err := entryStore.SeedTemplates(category.Key, templateModelsToStore(category.Key, category.Templates)); err != nil {
			return err
		}
	}
	setting = hydrateSetting(setting)
	return nil
}

type templatePayload struct {
	Name  string   `json:"Name"`
	Value string   `json:"Value"`
	Tags  []string `json:"Tags"`
}

func saveTemplate(w http.ResponseWriter, r *http.Request) {
	category := r.PathValue("category")
	if !hasCategory(setting, category) {
		writeErrorJSON(w, http.StatusBadRequest, errors.New("カテゴリが見つかりません"))
		return
	}

	var payload templatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	payload.Name = strings.TrimSpace(payload.Name)
	payload.Value = strings.TrimSpace(payload.Value)
	if payload.Name == "" {
		writeErrorJSON(w, http.StatusBadRequest, errors.New("テンプレート名を入力してください"))
		return
	}
	if payload.Value == "" {
		writeErrorJSON(w, http.StatusBadRequest, errors.New("本文が空のためテンプレート保存できません"))
		return
	}

	_, err := entryStore.SaveTemplate(category, store.Template{
		Name:  payload.Name,
		Value: payload.Value,
		Tags:  strings.Join(payload.Tags, "#"),
	})
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	latestSetting, err := loadSettingFromFile()
	if err == nil {
		setting = hydrateSetting(latestSetting)
	}
	writeJSON(w, http.StatusOK, setting)
}

func deleteTemplate(w http.ResponseWriter, r *http.Request) {
	category := r.PathValue("category")
	if !hasCategory(setting, category) {
		writeErrorJSON(w, http.StatusBadRequest, errors.New("カテゴリが見つかりません"))
		return
	}

	name := strings.TrimSpace(r.PathValue("name"))
	if name == "" {
		writeErrorJSON(w, http.StatusBadRequest, errors.New("テンプレート名を入力してください"))
		return
	}

	if err := entryStore.DeleteTemplate(category, name); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	latestSetting, err := loadSettingFromFile()
	if err == nil {
		setting = hydrateSetting(latestSetting)
	}
	writeJSON(w, http.StatusOK, setting)
}

func hasCategory(setting SettingModel, category string) bool {
	for _, categorySetting := range setting.Categories {
		if categorySetting.Key == category {
			return true
		}
	}
	return false
}

func getEntries(w http.ResponseWriter, r *http.Request) {
	category := r.PathValue("category")
	if r.URL.Query().Get("months") == "1" {
		months, err := entryStore.FindMonths(category)
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
		entries, err = entryStore.FindEntriesByMonth(category, month)
	} else {
		entries = entryStore.FindEntries(category)
	}
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	response := make([]entryRequest, 0, len(entries))
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
	payload := entryToPayload(entry)
	images, err := findEntryImages(r.PathValue("category"), r.PathValue("id"))
	if err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}
	payload.Images = images
	writeJSON(w, http.StatusOK, payload)
}

func putEntry(w http.ResponseWriter, r *http.Request) {
	saveEntry(w, r, "")
}

func postEntry(w http.ResponseWriter, r *http.Request) {
	saveEntry(w, r, r.PathValue("id"))
}

func saveEntry(w http.ResponseWriter, r *http.Request, pathID string) {
	category := r.PathValue("category")

	var request entryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	entry := payloadToEntry(category, request)
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
		entry = entryStore.UpsertEntry(entry)
	}

	value, err := Move(entry.Value, category, strconv.Itoa(entry.Id))
	if err != nil {
		if isNewEntry {
			entryStore.DeleteEntry(category, strconv.Itoa(entry.Id))
		}
		writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	entry.Value = value
	entry = entryStore.UpsertEntry(entry)
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

		current, err := entryStore.FindEntryByDate(category, entry.Date)
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
	return entryStore.FindEntryByID(category, id)
}

func deleteEntry(w http.ResponseWriter, r *http.Request) {
	category := r.PathValue("category")
	id := r.PathValue("id")
	var legacyDiaryDate string

	if category == "diary" {
		entry, err := entryStore.FindEntryByID(category, id)
		if err == nil {
			legacyDiaryDate = entry.Date
		}
	}

	entryStore.DeleteEntry(category, id)
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
	storedImage, err := entryStore.FindImage(imagePath)
	if errors.Is(err, gorm.ErrRecordNotFound) && r.PathValue("category") == "diary" {
		if legacyPath, legacyErr := legacyDiaryImagePath(r.PathValue("id"), r.PathValue("file")); legacyErr == nil {
			storedImage, err = entryStore.FindImage(legacyPath)
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

func deleteImage(w http.ResponseWriter, r *http.Request) {
	imagePath := imagePathForRoute(r.PathValue("category"), r.PathValue("id"), r.PathValue("file"))
	if err := entryStore.DeleteImage(imagePath); err != nil {
		writeErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	if r.PathValue("category") == "diary" {
		if legacyPath, err := legacyDiaryImagePath(r.PathValue("id"), r.PathValue("file")); err == nil {
			_ = entryStore.DeleteImage(legacyPath)
		}
	}

	w.WriteHeader(http.StatusOK)
}

func findEntryImages(category string, id string) ([]imageRequest, error) {
	prefixes := entryImagePrefixes(category, id)

	seen := make(map[string]struct{})
	items := make([]imageRequest, 0)
	for _, prefix := range prefixes {
		images, err := entryStore.FindImagesByPrefix(prefix)
		if err != nil {
			return nil, err
		}
		for _, image := range images {
			src := imagePathToURL(image.Path, category, id)
			if src == "" {
				continue
			}
			if _, ok := seen[src]; ok {
				continue
			}
			seen[src] = struct{}{}
			alt := path.Base(image.Path)
			items = append(items, imageRequest{
				Id:       src,
				Src:      src,
				Alt:      alt,
				Markdown: fmt.Sprintf("![%s](%s)", alt, src),
			})
		}
	}

	return items, nil
}

func entryImagePrefixes(category string, id string) []string {
	prefixes := []string{path.Join(category, fmt.Sprintf("%05d", mustAtoi(id))) + "/"}
	if category != "diary" {
		return prefixes
	}

	if entry, err := entryStore.FindEntryByID(category, id); err == nil && strings.TrimSpace(entry.Date) != "" {
		prefixes = append(prefixes, path.Join(category, strings.ReplaceAll(entry.Date, "-", ""))+"/")
	}
	return prefixes
}

func mustAtoi(v string) int {
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return n
}

func imagePathToURL(imagePath string, category string, id string) string {
	n, err := strconv.Atoi(id)
	if err == nil {
		padded := fmt.Sprintf("%05d", n)
		if strings.HasPrefix(imagePath, path.Join(category, padded)+"/") {
			return fmt.Sprintf("/api/%s/%s/images/%s", category, id, path.Base(imagePath))
		}
	}
	if category == "diary" {
		for _, prefix := range entryImagePrefixes(category, id) {
			if strings.HasPrefix(imagePath, prefix) {
				return fmt.Sprintf("/api/%s/%s/images/%s", category, id, path.Base(imagePath))
			}
		}
	}
	return ""
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
	entry, err := entryStore.FindEntryByID("diary", id)
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
