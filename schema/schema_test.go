package schema_test

import (
	"testing"

	"github.com/markelrep/csvalidator/schema"

	"github.com/stretchr/testify/assert"
)

func TestSchema_Parse(t *testing.T) {
	s, err := schema.Parse("../samples/schema.json")
	assert.NoError(t, err)
	expected := schema.Schema{
		Columns: []schema.Column{
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
