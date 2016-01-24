package repo

import (
	"encoding/json"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/language/en"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"path/filepath"
)

type BleveCache struct {
	index  bleve.Index
	dbPath string
}

func OpenBleveCache(dbPath string) (LocalCache, error) {
	index, err := bleve.Open(dbPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		indexMapping := buildIndexMapping()
		index, err = bleve.New(dbPath, indexMapping)
		if err != nil {
			return nil, err
		}

	} else if err != nil {
		return nil, err
	}
	return &BleveCache{
		dbPath: dbPath,
		index:  index,
	}, nil
}

func (cache *BleveCache) Search(ctx context.Context, term string, options SearchOptions) ([]*RecipeLink, error) {
	q := bleve.NewQueryStringQuery(term)
	req := bleve.NewSearchRequest(q)
	res, err := cache.index.Search(req)
	if err != nil {
		return nil, err
	}
	links := make([]*RecipeLink, len(res.Hits))
	for i, hit := range res.Hits {
		filePath := cache.dbPath + "/" + hit.ID + ".json"
		links[i], err = loadRecipeLink(filePath)
		if err != nil {
			return links, err
		}
	}
	return links, err
}

func loadRecipeLink(jsonFilePath string) (*RecipeLink, error) {
	file, err := os.Open(jsonFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var data RecipeLink
	dec := json.NewDecoder(file)
	err = dec.Decode(&data)
	return &data, err
}

func (cache *BleveCache) Add(link *RecipeLink) error {
	id := link.Title
	filePath := cache.dbPath + "/" + id + ".json"
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	err = enc.Encode(link)
	if err != nil {
		return err
	}
	return cache.index.Index(id, link)
}

func (cache *BleveCache) Close() error {
	return cache.index.Close()
}
func (cache *BleveCache) Cached(link *RecipeLink) bool {
	id := link.Title
	doc, err := cache.index.Document(id)
	return err == nil && doc != nil
}

func buildIndexMapping() *bleve.IndexMapping {
	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	// RecipeLink
	linkMapping := bleve.NewDocumentMapping()
	linkMapping.AddFieldMappingsAt("title", englishTextFieldMapping)
	linkMapping.AddFieldMappingsAt("description", englishTextFieldMapping)
	linkMapping.AddFieldMappingsAt("url", englishTextFieldMapping)

	// Recipe
	recipeMapping := bleve.NewDocumentMapping()
	recipeMapping.AddFieldMappingsAt("name", englishTextFieldMapping)
	recipeMapping.AddFieldMappingsAt("description", englishTextFieldMapping)
	recipeMapping.AddFieldMappingsAt("author", englishTextFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultAnalyzer = "en"
	indexMapping.AddDocumentMapping("title", linkMapping)
	//indexMapping.AddDocumentMapping("recipe", recipeMapping)

	return indexMapping
}

func indexRecipeLink(i bleve.Index, jsonFilePath string) error {
	jsonBytes, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return err
	}
	fileName := filepath.Base(jsonFilePath)
	ext := filepath.Ext(fileName)
	docId := fileName[:(len(fileName) - len(ext))]
	return i.Index(docId, jsonBytes)
}
