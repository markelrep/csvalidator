package csvalidator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator_Validate(t *testing.T) {
	validator, err := NewValidator("./samples/file.csv", "./samples/schema.json", true)
	assert.NoError(t, err)
	err = validator.Validate()
	assert.NoError(t, err)
}
