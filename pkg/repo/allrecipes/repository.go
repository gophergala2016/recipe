package allrecipes

import (
	"github.com/gophergala2016/recipe/pkg/repo"
	"github.com/gophergala2016/recipe/pkg/schema"
	"golang.org/x/net/context"
)

type Repository struct {
	cache repo.LocalCache
}

func New(cache repo.LocalCache) repo.Repository {
	r := &Repository{
		cache: cache,
	}
	return r
}

func (r *Repository) Refresh(ctx context.Context) error {
	// TODO Use cache, only index until last entry is reached
	entries := index(ctx)
	for entry := range entries {
		r.cache.Add(entry)
	}
	return nil
}

func (r *Repository) Search(ctx context.Context, term string) ([]repo.RecipeLink, error) {
	return r.cache.Search(term)
}

func (r *Repository) Get(ctx context.Context, url string) (*schema.Recipe, error) {
	// TODO Get recipe
	return nil, nil
}

type Entry struct {
	Title       string
	Description string
	URL         string
}
