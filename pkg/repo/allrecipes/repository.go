package allrecipes

import (
	"github.com/PuerkitoBio/goquery"
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
	index(ctx, r.cache)
	return nil
}

func (r *Repository) Search(ctx context.Context, term string, opt repo.SearchOptions) ([]repo.RecipeLink, error) {
	return r.cache.Search(term, opt)
}

func (r *Repository) Get(ctx context.Context, url string) ([]*schema.Recipe, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}
	return schema.ParseRecipes(doc)
}

type Entry struct {
	Title       string
	Description string
	URL         string
}
