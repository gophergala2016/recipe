package recipe

import (
	"fmt"
	"github.com/gophergala2016/recipe/pkg/repo"
	"github.com/gophergala2016/recipe/pkg/repo/allrecipes"
	"golang.org/x/net/context"
)

func Search(dbPath, term string, options repo.SearchOptions) ([]*repo.RecipeLink, error) {
	fmt.Println("Loading database...")
	cache, err := repo.OpenBleveCache(dbPath)
	if err != nil {
		return nil, err
	}
	defer cache.Close()

	ctx := context.Background()
	repo := allrecipes.New(cache)

	return repo.Search(ctx, term, options)
}
