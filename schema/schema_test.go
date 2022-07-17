package schema

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchema_Parse(t *testing.T) {
	s, err := Parse("../samples/schema.json")
	assert.NoError(t, err)
	expected := Schema{
		Columns: []column{
			{
				Name:     "id",
				Required: true,
				Contains: contains{
					kind:    regexpPattern,
					pattern: regexp.MustCompile(`^([0-9]{1})$`),
				},
			},
			{
				Name: "comment",
				Contains: contains{
					kind: list,
					list: map[any]struct{}{"comment": {}},
				},
			},
		},
	}
	assert.Equal(t, expected, s)
}
