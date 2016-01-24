package schema

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"time"
)

type Thing struct {
	Image       string `json:"image"`
	Description string `json:"description"`
	Name        string `json:"name"`
	// TODO Missing properties https://schema.org/Thing
}

type CreativeWork struct {
	Thing
	Author        string    `json:"author"`
	DatePublished time.Time `json:"datePublished"`
	// TODO Missing a lot of properties https://schema.org/CreativeWork
}

// Recipe for cooking.
type Recipe struct {
	CreativeWork
	CookTime      time.Time             `json:"cookTime"`
	CookingMethod string                `json:"cookingMethod"`
	Nutrition     *NutritionInformation `json:"nutrition"`
	PrepTime      time.Time             `json:"prepTime"`
	Category      string                `json:"recipeCategory"`
	Cuisine       string                `json:"recipeCuisine"`
	Ingredients   []string              `json:"recipeIngredient"`
	Instructions  string                `json:"recipeInstructions"`
	Yield         string                `json:"recipeYield"`
	TotalTime     time.Time             `json:"totalTime"`
}

// NutritionInformation for a recipe.
type NutritionInformation struct {
	Calories   string `json:"calories"`
	FatContent string `json:"fatContent"`
	// TODO Missing a lot of properties https://schema.org/NutritionInformation
}

func ParseNutritionInformation(sel *goquery.Selection) (*NutritionInformation, error) {
	if _, exists := sel.Attr("itemscope"); !exists {
		return nil, ErrMissingItemScope
	}
	itemtype, exists := sel.Attr("itemtype")
	if !exists {
		return nil, ErrMissingItemType
	}
	if itemtype != NutritionInformationSchemaURL {
		return nil, ErrWrongItemType
	}

	caloriesSel := sel.Find("[itemprop='calories']")
	calories := strings.TrimSpace(caloriesSel.Text())

	fatContentSel := sel.Find("[itemprop='fatContent']")
	fatContent := strings.TrimSpace(fatContentSel.Text())

	return &NutritionInformation{
		Calories:   calories,
		FatContent: fatContent,
	}, nil
}

var (
	ErrMissingItemScope = errors.New("Couldn't find itemscope attribute on node")
	ErrMissingItemType  = errors.New("Couldn't find itemtype attribute on node")
	ErrWrongItemType    = errors.New("Wrong itemtype attribute value")
)

const (
	RecipeSchemaURL               = "http://schema.org/Recipe"
	NutritionInformationSchemaURL = "http://schema.org/NutritionInformation"
)

func ParseRecipes(doc *goquery.Document) ([]*Recipe, error) {
	recipesSel := doc.Find("[itemscope=''][itemtype='http://schema.org/Recipe']")
	var (
		err     error
		recipes = make([]*Recipe, 0)
	)
	recipesSel.EachWithBreak(func(i int, article *goquery.Selection) bool {
		recipe, e := ParseRecipe(article)
		if e != nil {
			err = e
			return true
		}
		recipes = append(recipes, recipe)
		return false
	})
	return recipes, err
}

func ParseRecipe(sel *goquery.Selection) (*Recipe, error) {
	if _, exists := sel.Attr("itemscope"); !exists {
		return nil, ErrMissingItemScope
	}
	itemtype, exists := sel.Attr("itemtype")
	if !exists {
		return nil, ErrMissingItemType
	}
	if itemtype != RecipeSchemaURL {
		return nil, ErrWrongItemType
	}
	recipe := &Recipe{
		CreativeWork: CreativeWork{
			Thing: Thing{},
		},
	}

	nameSel := sel.Find("[itemprop='name']")
	recipe.Name = strings.TrimSpace(nameSel.Text())

	authorSel := sel.Find("[itemprop='author']").First()
	recipe.Author = strings.TrimSpace(authorSel.Text())

	datePublishedSel := sel.Find("[itemprop='datePublished']")
	datePublishedText, exists := datePublishedSel.Attr("content")
	if !exists {
		datePublishedText = datePublishedSel.Text()
	}
	var err error
	if len(datePublishedText) != 0 {
		recipe.DatePublished, err = time.Parse("2006-01-02", datePublishedText)
		if err != nil {
			return nil, err
		}
	}

	nutritionInformationSel := sel.Find(fmt.Sprintf("[itemscope=''][itemtype='%s']", NutritionInformationSchemaURL))
	if nutritionInformationSel.Size() > 0 {
		recipe.Nutrition, err = ParseNutritionInformation(nutritionInformationSel)
		if err != nil {
			return nil, err
		}
	}

	imageSel := sel.Find("[itemprop='image']")
	recipe.Image, _ = imageSel.Attr("src")

	descriptionSel := sel.Find("[itemprop='description']")
	recipe.Description = strings.TrimSpace(descriptionSel.Text())

	return recipe, nil
}
