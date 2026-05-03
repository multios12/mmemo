package web

import (
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/multios12/mmemo/pkg/markdown"
	"github.com/multios12/mmemo/pkg/store"
)

var (
	oldDiaryImagePattern = regexp.MustCompile(`/api/diary/(\d{4})/(\d{2})/(\d{2})/images/([^)"]+)`)
	oldEntryImagePattern = regexp.MustCompile(`/api/memos/([^/]+)/(\d{5})/([^)"]+)`)
)

type entryRequest struct {
	Id        int      `json:"Id,omitempty"`
	Outline   string   `json:"Outline"`
	Date      string   `json:"Date"`
	Category  string   `json:"Category,omitempty"`
	Tags      []string `json:"Tags"`
	Value     string   `json:"Value"`
	HTML      string   `json:"HTML,omitempty"`
	HasDetail bool     `json:"HasDetail"`
	CreatedAt string   `json:"CreatedAt,omitempty"`
	UpdatedAt string   `json:"UpdatedAt,omitempty"`
}

func entryToPayload(entry store.Entry) entryRequest {
	value := normalizeEntryValue(entry)
	return entryRequest{
		Id:        entry.Id,
		Outline:   entry.Outline,
		Date:      entry.Date,
		Category:  entry.Category,
		Tags:      splitTags(entry.Tags),
		Value:     value,
		HTML:      markdown.ToHTML(value),
		HasDetail: strings.TrimSpace(value) != "",
		CreatedAt: entry.CreatedAt.Format(time.RFC3339),
		UpdatedAt: entry.UpdatedAt.Format(time.RFC3339),
	}
}

func payloadToEntry(category string, payload entryRequest) store.Entry {
	return store.Entry{
		Id:       payload.Id,
		Outline:  strings.TrimSpace(payload.Outline),
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

	return oldEntryImagePattern.ReplaceAllStringFunc(value, func(match string) string {
		subMatches := oldEntryImagePattern.FindStringSubmatch(match)
		if len(subMatches) != 4 {
			return match
		}
		return fmt.Sprintf("/api/%s/%s/images/%s", subMatches[1], subMatches[2], path.Base(subMatches[3]))
	})
}
