package memo

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/multios12/mmemo/pkg/store"
)

var (
	oldDiaryImagePattern = regexp.MustCompile(`/api/diary/(\d{4})/(\d{2})/(\d{2})/images/([^)"]+)`)
	oldMemoImagePattern  = regexp.MustCompile(`/api/memos/([^/]+)/(\d{5})/([^)"]+)`)
)

type entryPayload struct {
	Id        int      `json:"Id,omitempty"`
	Title     string   `json:"Title"`
	Date      string   `json:"Date"`
	Category  string   `json:"Category,omitempty"`
	Tags      []string `json:"Tags"`
	Value     string   `json:"Value"`
	HasDetail bool     `json:"HasDetail"`
}

func entryToPayload(entry store.Entry) entryPayload {
	value := normalizeEntryValue(entry)
	return entryPayload{
		Id:        entry.Id,
		Title:     entry.Title,
		Date:      entry.Date,
		Category:  entry.Category,
		Tags:      splitTags(entry.Tags),
		Value:     value,
		HasDetail: strings.TrimSpace(value) != "",
	}
}

func payloadToEntry(category string, payload entryPayload) store.Entry {
	return store.Entry{
		Id:       payload.Id,
		Title:    strings.TrimSpace(payload.Title),
		Date:     strings.TrimSpace(payload.Date),
		Category: category,
		Tags:     strings.Join(payload.Tags, "#"),
		Value:    payload.Value,
	}
}

func splitTags(v string) []string {
	if strings.TrimSpace(v) == "" {
		return []string{}
	}
	items := strings.Split(v, "#")
	tags := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			tags = append(tags, item)
		}
	}
	return tags
}

func normalizeEntryValue(entry store.Entry) string {
	value := entry.Value
	if entry.Category == "diary" {
		return oldDiaryImagePattern.ReplaceAllStringFunc(value, func(match string) string {
			subMatches := oldDiaryImagePattern.FindStringSubmatch(match)
			if len(subMatches) != 5 {
				return match
			}
			return fmt.Sprintf("/api/diary/%d/images/%s", entry.Id, path.Base(subMatches[4]))
		})
	}

	return oldMemoImagePattern.ReplaceAllStringFunc(value, func(match string) string {
		subMatches := oldMemoImagePattern.FindStringSubmatch(match)
		if len(subMatches) != 4 {
			return match
		}
		return fmt.Sprintf("/api/%s/%s/images/%s", subMatches[1], subMatches[2], path.Base(subMatches[3]))
	})
}
