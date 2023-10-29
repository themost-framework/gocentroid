package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUseEscapeName(t *testing.T) {
	var dialect = &SqlDialect{
		NameFormat: "`$1`",
	}
	name, err := dialect.EscapeName("ProductData.price")
	assert.Nil(t, err)
	assert.Equal(t, "`ProductData`.`price`", name)
}

func TestUseEscape(t *testing.T) {
	var dialect = &SqlDialect{
		NameFormat: "`$1`",
	}
	dialect.Init(DefaultSqlDialectOptions())
	sql, err := dialect.Escape(QueryElement{
		"$eq": []any{
			QueryElement{
				"$getField": "price",
			},
			"Apple",
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, "`price` = 'Apple'", sql)
}
