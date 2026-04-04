package entryapi

import "github.com/multios12/mmemo/pkg/store"

type Store interface {
	Open(dataPath string) error
	FindEntries(category string) []store.Entry
	FindEntryByID(category string, id string) (store.Entry, error)
	FindEntriesByMonth(category string, month string) ([]store.Entry, error)
	FindEntryByDate(category string, day string) (store.Entry, error)
	FindMonths(category string) ([]string, error)
	UpsertEntry(entry store.Entry) store.Entry
	DeleteEntry(category string, id string)
	FindImage(path string) (store.Image, error)
	FindTemplates(category string) ([]store.Template, error)
	SaveTemplate(category string, template store.Template) (store.Template, error)
	SeedTemplates(category string, templates []store.Template) error
}

type sqliteStore struct{}

func SetStore(next Store) {
	if next == nil {
		entryStore = sqliteStore{}
		return
	}
	entryStore = next
}

func (sqliteStore) Open(dataPath string) error {
	return store.Open(dataPath)
}

func (sqliteStore) FindEntries(category string) []store.Entry {
	return store.FindEntries(category)
}

func (sqliteStore) FindEntryByID(category string, id string) (store.Entry, error) {
	return store.FindEntryByID(category, id)
}

func (sqliteStore) FindEntriesByMonth(category string, month string) ([]store.Entry, error) {
	return store.FindEntriesByMonth(category, month)
}

func (sqliteStore) FindEntryByDate(category string, day string) (store.Entry, error) {
	return store.FindEntryByDate(category, day)
}

func (sqliteStore) FindMonths(category string) ([]string, error) {
	return store.FindMonths(category)
}

func (sqliteStore) UpsertEntry(entry store.Entry) store.Entry {
	return store.UpsertEntry(entry)
}

func (sqliteStore) DeleteEntry(category string, id string) {
	store.DeleteEntry(category, id)
}

func (sqliteStore) FindImage(path string) (store.Image, error) {
	return store.FindImage(path)
}

func (sqliteStore) FindTemplates(category string) ([]store.Template, error) {
	return store.FindTemplates(category)
}

func (sqliteStore) SaveTemplate(category string, template store.Template) (store.Template, error) {
	return store.SaveTemplate(category, template)
}

func (sqliteStore) SeedTemplates(category string, templates []store.Template) error {
	return store.SeedTemplates(category, templates)
}
