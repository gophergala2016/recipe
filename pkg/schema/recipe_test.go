package schema

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"testing"
	"time"
)

type testCase struct {
	document       []byte
	expectedRecipe *Recipe
}

var testCases = []*testCase{
	&testCase{
		document: schemaDotOrgExample,
		expectedRecipe: &Recipe{
			CreativeWork: CreativeWork{
				Thing: Thing{
					Name:        "Mom's World Famous Banana Bread",
					Image:       "bananabread.jpg",
					Description: "This classic banana bread recipe comes\n  from my mom -- the walnuts add a nice texture and flavor to the banana\n  bread.",
				},
				Author:        "John Smith",
				DatePublished: time.Date(2009, 5, 8, 0, 0, 0, 0, time.UTC),
			},
			Nutrition: &NutritionInformation{
				Calories:   "240 calories",
				FatContent: "9 grams fat",
			},
		},
	},
	&testCase{
		document: allrecipeDotComPage,
		expectedRecipe: &Recipe{
			CreativeWork: CreativeWork{
				Thing: Thing{
					Name:        "Zucchini Lasagna With Beef and Sausage",
					Image:       "http://images.media-allrecipes.com/userphotos/720x405/1104956.jpg",
					Description: `"This recipe is perfect if you have extra zucchini from the garden and/or you are looking for a great lasagna while on the South Beach or Atkins diets. It replaces lasagna noodles with slices of zucchini, but still tastes like the lasagna you love!"`,
				},
				Author: "Jeff B.",
			},
			Nutrition: &NutritionInformation{
				Calories:   "471 kcal",
				FatContent: "25.3 g",
			},
		},
	},
}

func TestParseRecipes(t *testing.T) {
	for _, tc := range testCases {
		testParseRecipes(t, tc.document, tc.expectedRecipe)
	}
}

func testParseRecipes(t *testing.T, b []byte, expectedRecipe *Recipe) {
	r := bytes.NewReader(b)
	doc, err := goquery.NewDocumentFromReader(r)
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

	if err != nil {
		t.Fatal(err)
	}

	if recipe.Name != expectedRecipe.Name {
		t.Errorf("Expected name to be %q, got %q", expectedRecipe.Name, recipe.Name)
	}
	if recipe.Author != expectedRecipe.Author {
		t.Errorf("Expected author to be %q, got %q", expectedRecipe.Author, recipe.Author)
	}
	if recipe.DatePublished != expectedRecipe.DatePublished {
		t.Errorf("Expected datePublished to be %q, got %q", expectedRecipe.DatePublished, recipe.DatePublished)
	}
	if recipe.Image != expectedRecipe.Image {
		t.Errorf("Expected image to be %q, got %q", expectedRecipe.Image, recipe.Image)
	}
	if recipe.Description != expectedRecipe.Description {
		t.Errorf("Expected description to be %q, got %q", expectedRecipe.Description, recipe.Description)
	}
	if recipe.PrepTime != expectedRecipe.PrepTime {
		t.Errorf("Expected prepTime to be %q, got %q", expectedRecipe.PrepTime, recipe.PrepTime)
	}
	if recipe.CookTime != expectedRecipe.CookTime {
		t.Errorf("Expected cookTime to be %q, got %q", expectedRecipe.CookTime, recipe.CookTime)
	}
	if recipe.Yield != expectedRecipe.Yield {
		t.Errorf("Expected recipeYield to be %q, got %q", expectedRecipe.Yield, recipe.Yield)
	}
	if recipe.Nutrition != nil {
		if recipe.Nutrition.Calories != expectedRecipe.Nutrition.Calories {
			t.Errorf("Expected nutrition.calories to be %q, got %q", expectedRecipe.Nutrition.Calories, recipe.Nutrition.Calories)
		}
		if recipe.Nutrition.FatContent != expectedRecipe.Nutrition.FatContent {
			t.Errorf("Expected nutrition.fatContent to be %q, got %q", expectedRecipe.Nutrition.FatContent, recipe.Nutrition.FatContent)
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
