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

func TestValidator_Validate_Contains(t *testing.T) {
	validator, err := NewValidator("./samples/file_2.csv", "./samples/schema_2.json", true)
	assert.NoError(t, err)
	err = validator.Validate()
	assert.EqualError(t, err, "[Row 4 | Column 1] expected one of: 'only, good, words'. Actual is 'bad!'")
}