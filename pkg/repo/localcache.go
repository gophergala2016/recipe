package repo

import (
	"encoding/json"
	"golang.org/x/net/context"
	"io"
	"os"
	"strings"
)

type LocalCache interface {
	Search(context.Context, string, SearchOptions) ([]*RecipeLink, error)
	Add(*RecipeLink) error
	Close() error
	Cached(*RecipeLink) bool
}

type RecipeLink struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type JsonFileCache struct {
	Entries  map[string]*RecipeLink `json:"entries"`
	fileName string                 `json:"-"`
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

func (cache *JsonFileCache) Cached(entry *RecipeLink) bool {
	id := entry.URL
	_, exists := cache.Entries[id]
	return exists
}

func (cache *JsonFileCache) Add(entry *RecipeLink) error {
	id := entry.URL
	cache.Entries[id] = entry
	return nil
}

func (cache *JsonFileCache) Search(ctx context.Context, term string, options SearchOptions) ([]*RecipeLink, error) {
	if entry, ok := cache.Entries[term]; ok {
		return []*RecipeLink{entry}, nil
	}
	result := make([]*RecipeLink, 0)
	for _, value := range cache.Entries {
		if strings.Contains(value.Description, term) && options.Description ||
			strings.Contains(value.Title, term) && options.Title ||
			strings.Contains(value.URL, term) && options.URL {
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
		cache.Entries = make(map[string]*RecipeLink)
		return cache, nil
	}
	return cache, err
}
