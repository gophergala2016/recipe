package recipe

import (
	"github.com/gophergala2016/recipe/pkg/repo"
	"github.com/gophergala2016/recipe/pkg/repo/allrecipes"
	"golang.org/x/net/context"
)

func Refresh(cacheFileName string) error {
	cache, err := repo.OpenJsonFileCache(cacheFileName)
	if err != nil {
		return err
	}
	defer cache.Close()

	ctx := context.Background()
	repo := allrecipes.New(cache)
	return repo.Refresh(ctx)
}
