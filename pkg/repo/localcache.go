package repo

import (
	"encoding/json"
	"io"
	"os"
	"strings"
)

type LocalCache interface {
	Search(term string) ([]RecipeLink, error)
	Add(RecipeLink) error
	Close() error
	Cached(RecipeLink) bool
}

type RecipeLink interface {
	Title() string
	Description() string
	URL() string
}

type jsonRecipeEntry struct {
	RecipeTitle       string `json:"title"`
	RecipeDescription string `json:"description"`
	RecipeURL         string `json:"url"`
}

func (entry *jsonRecipeEntry) Title() string {
	return entry.RecipeTitle
}

func (entry *jsonRecipeEntry) Description() string {
	return entry.RecipeDescription
}

func (entry *jsonRecipeEntry) URL() string {
	return entry.RecipeURL
}

type JsonFileCache struct {
	Entries  map[string]*jsonRecipeEntry `json:"entries"`
	fileName string                      `json:"-"`
}

func (cache *JsonFileCache) Close() error {
	file, err := os.Create(cache.fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	return enc.Encode(cache)
}

func (cache *JsonFileCache) Cached(entry RecipeLink) bool {
	id := entry.URL()
	_, exists := cache.Entries[id]
	return exists
}

func (cache *JsonFileCache) Add(entry RecipeLink) error {
	id := entry.URL()
	switch jEntry := entry.(type) {
	case *jsonRecipeEntry:
		cache.Entries[id] = jEntry
		return nil
	default:
	}
	jEntry := &jsonRecipeEntry{
		RecipeDescription: entry.Description(),
		RecipeTitle:       entry.Title(),
		RecipeURL:         entry.URL(),
	}
	cache.Entries[id] = jEntry
	return nil
}

func (cache *JsonFileCache) Search(term string) ([]RecipeLink, error) {
	if entry, ok := cache.Entries[term]; ok {
		return []RecipeLink{entry}, nil
	}
	result := make([]RecipeLink, 0)
	for _, value := range cache.Entries {
		if strings.Contains(value.Description(), term) || strings.Contains(value.Title(), term) || strings.Contains(value.URL(), term) {
			result = append(result, value)
		}
	}
	return result, nil
}

func OpenJsonFileCache(name string) (LocalCache, error) {
	file, err := os.Open(name)
	if err != nil {
		file, err = os.Create(name)
		if err != nil {
			return nil, err
		}
	}
	defer file.Close()
	cache := &JsonFileCache{
		fileName: name,
	}
	dec := json.NewDecoder(file)
	err = dec.Decode(cache)
	if err == io.EOF {
		cache.Entries = make(map[string]*jsonRecipeEntry)
		return cache, nil
	}
	return cache, err
}
