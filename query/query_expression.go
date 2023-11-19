package query

import (
	"fmt"

	"github.com/Goldziher/go-utils/maputils"
	"golang.org/x/exp/maps"
)

type QueryElement map[string]interface{}

type QueryDocument struct {
	Collection QueryElement
	Select     QueryElement
	OrderBy    []QueryElement
	GroupBy    []QueryElement
	Expand     []QueryDocument
	Top        int
	Skip       int
	Distinct   bool
}

type QueryExpression struct {
	Query *QueryDocument
}

func (query *QueryExpression) From(collection string) *QueryExpression {
	if query.Query == nil {
		query.Query = &QueryDocument{}
	}
	query.Query.Collection = QueryElement{
		collection: 1,
	}
	return query
}

func (query *QueryExpression) Select(args ...any) *QueryExpression {
	if query.Query == nil {
		query.Query = &QueryDocument{}
	}
	attributes := QueryElement{}
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			attributes = maputils.Merge(attributes, QueryElement{
				v: 1,
			})
		case QueryElement:
			// get first key
			keys := maps.Keys(v)
			if len(keys) == 0 {
				fmt.Println("the given query element is empty and it's going to be omitted")
			} else {
				key0 := keys[0]
				if key0[:1] == "$" {
					// query element should have an alias
					alias := fmt.Sprintf("field%v", len(maps.Keys(attributes))) // e.g. field0
					attributes = maputils.Merge(attributes, QueryElement{
						alias: v,
					})

				} else {
					attributes = maputils.Merge(attributes, v)
				}
			}
		}
	}
	query.Query.Select = attributes
	return query
}

func (query *QueryExpression) Distinct() *QueryExpression {
	if query.Query == nil {
		query.Query = &QueryDocument{}
	}
	query.Query.Distinct = true
	return query
}

func (query *QueryExpression) OrderBy(attr string) *QueryExpression {
	if query.Query == nil {
		query.Query = &QueryDocument{}
	}
	query.Query.OrderBy = []QueryElement{
		{
			"$asc": attr,
		},
	}
	return query
}

func (query *QueryExpression) OrderByDescending(attr string) *QueryExpression {
	if query.Query == nil {
		query.Query = &QueryDocument{}
	}
	query.Query.OrderBy = []QueryElement{
		{
			"$desc": attr,
		},
	}
	return query
}

func (query *QueryExpression) ThenBy(attr string) *QueryExpression {
	if query.Query == nil {
		query.Query = &QueryDocument{}
	}
	query.Query.OrderBy = append(query.Query.OrderBy, QueryElement{
		"$asc": attr,
	})
	return query
}

func (query *QueryExpression) ThenByDescending(attr string) *QueryExpression {
	if query.Query == nil {
		query.Query = &QueryDocument{}
	}
	if query.Query.OrderBy == nil {
		query.Query.OrderBy = make([]QueryElement, 0)
	}
	query.Query.OrderBy = append(query.Query.OrderBy, QueryElement{
		"$desc": attr,
	})
	return query
}

func (query *QueryExpression) GroupBy(args ...any) *QueryExpression {
	if query.Query == nil {
		query.Query = &QueryDocument{}
	}
	groupBy := make([]QueryElement, 0)
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			groupBy = append(groupBy, QueryElement{
				v: 1,
			})
		case QueryElement:
			groupBy = append(groupBy, v)
		}
	}
	query.Query.GroupBy = groupBy
	return query
}
