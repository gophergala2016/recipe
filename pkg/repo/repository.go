package repo

import (
	"github.com/gophergala2016/recipe/pkg/schema"
	"golang.org/x/net/context"
)

//go:generate stringer -type=SearchMode
type SearchMode byte

const (
	Contains SearchMode = iota
	BeginsWith
	ExactMatch
	WildCards
	RegularExpression
)

type SearchOptions struct {
	Title       bool
	Description bool
	URL         bool
}

type Repository interface {
	Refresh(context.Context) error
	Get(context.Context, string) (*schema.Recipe, error)
	Search(context.Context, string, SearchOptions) ([]RecipeLink, error)
}
