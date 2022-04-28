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

func TestValidator_NewValidator(t *testing.T) {
	cases := []struct {
		path        string
		schemaPath  string
		firstHeader bool
	}{
		{
			path:        "./samples",
			schemaPath:  "./samples/schema.json",
			firstHeader: true,
		},
		{
			path:        "./samples/file.csv",
			schemaPath:  "./samples/schema.json",
			firstHeader: true,
		},
	}
	for _, tc := range cases {
		_, err := NewValidator(tc.path, tc.schemaPath, tc.firstHeader)
		assert.NoError(t, err)
	}
}
