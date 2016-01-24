package recipe

import (
	"encoding/json"
	"fmt"
	"github.com/gophergala2016/recipe/pkg/repo"
	"github.com/gophergala2016/recipe/pkg/repo/allrecipes"
	"github.com/gophergala2016/recipe/pkg/schema"
	"golang.org/x/net/context"
	"os"
)

func Get(cacheFileName, term string) {
	fmt.Println("Loading database...")
	cache, err := repo.OpenJsonFileCache(cacheFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cache.Close()

	ctx := context.Background()
	r := allrecipes.New(cache)

	fmt.Println("Searching...")
	options := repo.SearchOptions{
		Title: true,
		URL:   true,
	}
	recipeLinks, err := r.Search(ctx, term, options)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, link := range recipeLinks {
		fmt.Printf("Downloading %d/%d...\n”", i, len(recipeLinks))
		url := "http://allrecipes.com" + link.URL() // TODO This should not contain indexer specific code
		recipes, err := r.Get(ctx, url)
		if err != nil {
			fmt.Println(err)
			return
		}
		// r.Get should currently only return one recipe but just in case go over all
		for _, recipe := range recipes {
			if err != nil {
				fmt.Printf("%s skipping...\n", err)
				continue
			}
			// TODO Don't use the recipe name, in future use an ID
			fileName := recipe.Name
			err = saveRecipe(fileName, recipe)
			if err != nil {
				fmt.Printf("%s skipping...\n", err)
				continue
			}
		}
	}
}

func saveRecipe(fileName string, recipe *schema.Recipe) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	return enc.Encode(recipe)
}
