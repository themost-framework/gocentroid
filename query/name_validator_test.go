package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewObjectNameValidator(t *testing.T) {
	var validator = &ObjectNameValidator{}
	validator.Init(DefaultNamePattern)
	assert.NotNil(t, validator.Pattern)
	var valid, _ = validator.Test("name", false, false)
	assert.True(t, valid)
	valid, _ = validator.Test("name", true, false)
	assert.True(t, valid)
	valid, _ = validator.Test("*", true, false)
	assert.True(t, valid)
	valid, _ = validator.Test("4*5", true, false)
	assert.False(t, valid)
	valid, _ = validator.Test("n*ame", false, false)
	assert.False(t, valid)
	valid, _ = validator.Test("Products.name", true, false)
	assert.True(t, valid)
	valid, _ = validator.Test("Products.n*ame", true, false)
	assert.False(t, valid)
	valid, _ = validator.Test("Products.*", true, false)
	assert.True(t, valid)
}

func TestDefaultObjectNameValidator(t *testing.T) {
	var validator = DefaultNameValidator()
	var valid, _ = validator.Test("name", false, false)
	assert.True(t, valid)
}
