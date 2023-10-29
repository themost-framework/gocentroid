package query

import (
	"encoding/json"
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

func TestUseEqual(t *testing.T) {
	var dialect = &SqlDialect{
		NameFormat: "`$1`",
	}
	dialect.Init(DefaultSqlDialectOptions())
	sql, err := dialect.Escape(QueryElement{
		"$eq": []any{
			QueryElement{
				"$name": "name",
			},
			"Macbook Air",
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, "`name` = 'Macbook Air'", sql)
}

func TestUseAnd(t *testing.T) {
	var dialect = &SqlDialect{
		NameFormat: "`$1`",
	}
	dialect.Init(DefaultSqlDialectOptions())
	sql, err := dialect.Escape(QueryElement{
		"$and": []any{
			QueryElement{
				"$eq": []any{
					QueryElement{
						"$name": "category",
					},
					"Laptops",
				},
			},
			QueryElement{
				"$lt": []any{
					QueryElement{
						"$name": "price",
					},
					900,
				},
			},
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, "`category` = 'Laptops' AND `price` < 900", sql)
}

func TestUnmarshalAndUseAnd(t *testing.T) {
	var dialect = &SqlDialect{
		NameFormat: "`$1`",
	}
	dialect.Init(DefaultSqlDialectOptions())
	var element QueryElement
	json.Unmarshal([]byte(`
	{
		"$and": [
			{
				"$eq": [
					{
						"$getField": "category"
					},
					"Laptops"
				]
			},
			{
				"$lt": [
					{
						"$getField": "price"
					},
					900
				]
			}
		]
	}
	`), &element)
	sql, err := dialect.Escape(element)
	if assert.NoError(t, err) {
		assert.Equal(t, "`category` = 'Laptops' AND `price` < 900", sql)
	}
}
