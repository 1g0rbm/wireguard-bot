package utils

import "github.com/Masterminds/squirrel"

type Filter struct {
	Eq   squirrel.Eq
	Like squirrel.Like
}

type FilterOption func(Filter)
