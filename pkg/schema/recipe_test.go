package schema

import (
	"github.com/PuerkitoBio/goquery"
	"os"
	"testing"
	"time"
)

func TestParseRecipes(t *testing.T) {
	file, err := os.Open("example_recipe.html")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		t.Fatal(err)
	}
	recipes, err := ParseRecipes(doc)
	if err != nil {
		t.Error(err)
	}

	if len(recipes) != 1 {
		t.Errorf("Expected 1 recipe to be parsed, got %d", len(recipes))
		return
	}
	recipe := recipes[0]

	expectedDatePublished, err := time.Parse("Jan 2, 2006", "May 8, 2009")
	if err != nil {
		t.Fatal(err)
	}
	expectedRecipe := &Recipe{
		CreativeWork: CreativeWork{
			Thing: Thing{
				Name:        "Mom's World Famous Banana Bread",
				Image:       "bananabread.jpg",
				Description: "This classic banana bread recipe comes\n  from my mom -- the walnuts add a nice texture and flavor to the banana\n  bread.",
			},
			Author:        "John Smith",
			DatePublished: expectedDatePublished,
		},
		Nutrition: &NutritionInformation{
			Calories:   "240 calories",
			FatContent: "9 grams fat",
		},
	}

	if recipe.Name != expectedRecipe.Name {
		t.Errorf("Expected name to be %s, got %s", expectedRecipe.Name, recipe.Name)
	}
	if recipe.Author != expectedRecipe.Author {
		t.Errorf("Expected author to be %s, got %s", expectedRecipe.Author, recipe.Author)
	}
	if recipe.DatePublished != expectedRecipe.DatePublished {
		t.Errorf("Expected datePublished to be %s, got %s", expectedRecipe.DatePublished, recipe.DatePublished)
	}
	if recipe.Image != expectedRecipe.Image {
		t.Errorf("Expected image to be %s, got %s", expectedRecipe.Image, recipe.Image)
	}
	if recipe.Description != expectedRecipe.Description {
		t.Errorf("Expected description to be %s, got %s", expectedRecipe.Description, recipe.Description)
	}
	if recipe.PrepTime != expectedRecipe.PrepTime {
		t.Errorf("Expected prepTime to be %s, got %s", expectedRecipe.PrepTime, recipe.PrepTime)
	}
	if recipe.CookTime != expectedRecipe.CookTime {
		t.Errorf("Expected cookTime to be %s, got %s", expectedRecipe.CookTime, recipe.CookTime)
	}
	if recipe.Yield != expectedRecipe.Yield {
		t.Errorf("Expected recipeYield to be %s, got %s", expectedRecipe.Yield, recipe.Yield)
	}
	if recipe.Nutrition != nil {
		if recipe.Nutrition.Calories != expectedRecipe.Nutrition.Calories {
			t.Errorf("Expected nutrition.calories to be %s, got %s", expectedRecipe.Nutrition.Calories, recipe.Nutrition.Calories)
		}
		if recipe.Nutrition.FatContent != expectedRecipe.Nutrition.FatContent {
			t.Errorf("Expected nutrition.fatContent to be %s, got %s", expectedRecipe.Nutrition.FatContent, recipe.Nutrition.FatContent)
		}
	} else {
		t.Error("Missing nutrition data")
	}
}

/*
TODO Test missing properties

Ingredients:
- <span itemprop="recipeIngredient">3 or 4 ripe bananas, smashed</span>
- <span itemprop="recipeIngredient">1 egg</span>
- <span itemprop="recipeIngredient">3/4 cup of sugar</span>
...
Instructions:
<span itemprop="recipeInstructions">
Preheat the oven to 350 degrees. Mix in the ingredients in a bowl. Add
the flour last. Pour the mixture into a loaf pan and bake for one hour.
</span>
140 comments:
<div itemprop="interactionStatistic" itemscope itemtype="http://schema.org/InteractionCounter">
  <meta itemprop="interactionType" content="http://schema.org/CommentAction" />
  <meta itemprop="userInteractionCount" content="140" />
</div>
From Janel, May 5 -- thank you, great recipe!
...
*/
