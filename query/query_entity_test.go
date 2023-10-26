package query

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestNewQueryEntity(t *testing.T) {
	collection := &QueryEntity{
		Name: "Products",
	}
	assert.Equal(t, collection.Name, "Products")
}

func TestSerializeQueryEntity(t *testing.T) {
	alias := "products"
	collection := &QueryEntity{
		Name:  "ProductData",
		Alias: &alias,
	}
	str, err := json.Marshal(collection)
	assert.Nil(t, err)
	assert.Equal(t, str, []byte("{\"products\":\"$ProductData\"}"))
}

func TestDeserializeQueryEntity(t *testing.T) {
	alias := "products"
	collection := &QueryEntity{
		Name:  "ProductData",
		Alias: &alias,
	}
	str, err := json.Marshal(collection)
	assert.Nil(t, err)
	var newCollection QueryEntity
	json.Unmarshal(str, &newCollection)
	assert.Equal(t, newCollection.Name, "ProductData")
}

func TestConvertToQueryEntity(t *testing.T) {
	data := map[string]interface{}{
		"products": "$ProductData",
	}
	bytes, _ := json.Marshal(data)
	var newCollection QueryEntity
	json.Unmarshal(bytes, &newCollection)
	assert.Equal(t, newCollection.Name, "ProductData")
	assert.Equal(t, *newCollection.Alias, "products")
	var collection QueryEntity
	mapstructure.Decode(data, collection)
	assert.Equal(t, newCollection.Name, "ProductData")

}
