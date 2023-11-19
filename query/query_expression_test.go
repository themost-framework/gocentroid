package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQueryExpression(t *testing.T) {
	expr := &QueryExpression{}
	expr.From("ProductData").Select("id", "name", "price")
	assert.NotEmpty(t, expr.Query)
	assert.Equal(t, expr.Query.Collection, QueryElement{
		"ProductData": 1,
	})
}

func TestUseSelect(t *testing.T) {
	expr := &QueryExpression{}
	expr.From("ProductData").Select("id", "name", QueryElement{
		"$round": []any{
			QueryElement{
				"$getField": "price",
			},
			2,
		},
	})
	assert.NotEmpty(t, expr.Query)
	assert.Equal(t, expr.Query.Collection, QueryElement{
		"ProductData": 1,
	})
}
