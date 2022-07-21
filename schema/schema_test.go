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
				RecordRegexp: contains{
					pattern: regexp.MustCompile(`^([0-9]{1})$`),
				},
			},
			{
				Name: "comment",
				RecordRegexp: contains{
					pattern: regexp.MustCompile(`^comment$`),
				},
			},
		},
	}
	assert.Equal(t, expected, s)
}
