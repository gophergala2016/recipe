package recipe

import (
	"github.com/gophergala2016/recipe/pkg/repo"
	"github.com/gophergala2016/recipe/pkg/repo/allrecipes"
	"golang.org/x/net/context"
)

func Search(cacheFileName, term string, options repo.SearchOptions) ([]repo.RecipeLink, error) {
	cache, err := repo.OpenJsonFileCache(cacheFileName)
	if err != nil {
		return nil, err
	}
	defer cache.Close()

	ctx := context.Background()
	repo := allrecipes.New(cache)

	return repo.Search(ctx, term, options)
}
