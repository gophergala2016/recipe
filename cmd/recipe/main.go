package main

import (
	"github.com/gophergala2016/recipe/pkg/repo"
	"github.com/gophergala2016/recipe/pkg/repo/allrecipes"
	"golang.org/x/net/context"
	"log"
)

var dbFile = "test.json"

func main() {
	cache, err := repo.OpenJsonFileCache(dbFile)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	ctx := context.Background()

	repo := allrecipes.New(cache)
	err = repo.Refresh(ctx)
	if err != nil {
		log.Printf("[ERROR]Refresh:%s\n", err)
		return
	}

	err = cache.Close()
	if err != nil {
		log.Printf("[ERROR]Cache.Close:%s\n", err)
		return
	}
}
