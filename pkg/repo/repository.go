package repo

import (
	"github.com/gophergala2016/recipe/pkg/schema"
	"golang.org/x/net/context"
)

type Repository interface {
	Refresh(context.Context) error
	Get(context.Context, string) (*schema.Recipe, error)
	Search(context.Context, string) ([]RecipeLink, error)
}
