package web

import (
	"testing"

	"github.com/multios12/mmemo/pkg/store"
)

func TestNormalizeEntryValue_RewritesRelativeEntryImagePaths(t *testing.T) {
	entry := store.Entry{
		Category: "memo",
		Value:    "![image](./memo/123/sample.png)",
	}

	got := normalizeEntryValue(entry)
	want := "![image](./memo/123/sample.png)"

	if got != want {
		t.Fatalf("normalizeEntryValue() = %q, want %q", got, want)
	}
}

func TestNormalizeEntryValue_RewritesLegacyEntryImagePaths(t *testing.T) {
	entry := store.Entry{
		Category: "memo",
		Value:    "![image](/api/memos/memo/00012/sample.png)",
	}

	got := normalizeEntryValue(entry)
	want := "![image](./memo/00012/sample.png)"

	if got != want {
		t.Fatalf("normalizeEntryValue() = %q, want %q", got, want)
	}
}
