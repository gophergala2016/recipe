package recipe

import (
	"fmt"
	"github.com/gophergala2016/recipe/pkg/repo"
	"github.com/gophergala2016/recipe/pkg/repo/allrecipes"
	"golang.org/x/net/context"
)

func Refresh(dbPath string) error {
	// cache, err := repo.OpenJsonFileCache(cacheFileName)
	fmt.Println("Loading database...")
	cache, err := repo.OpenBleveCache(dbPath)
	if err != nil {
		return err
	}
	defer cache.Close()

	ctx := context.Background()
	repo := allrecipes.New(cache)
	fmt.Println("Refreshing...")
	return repo.Refresh(ctx)
}
