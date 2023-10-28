package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQueryExpression(t *testing.T) {
	expr := &QueryExpression{}
	expr.from("ProductData").pick("id", "name")
	assert.NotEmpty(t, expr.Query)
	assert.Equal(t, expr.Query["$collection"], QueryElement{
		"ProductData": 1,
	})
}
