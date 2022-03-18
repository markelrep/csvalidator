package csvalidator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchema_ParseSchema(t *testing.T) {
	s, err := parseSchema("./samples/schema.json")
	assert.NoError(t, err)
	expected := schema{
		Columns: []column{
			{
				Name:     "id",
				Required: true,
			},
			{
				Name: "comment",
			},
		},
	}
	assert.Equal(t, expected, s)
}
