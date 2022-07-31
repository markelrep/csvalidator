package schema

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchema_Parse(t *testing.T) {
	s, err := Parse("../samples/schema_type.json")
	assert.NoError(t, err)
	expected := Schema{
		Columns: []Column{
			{
				Name:     "id",
				Required: true,
				RecordRegexp: recordRegexp{
					pattern: regexp.MustCompile(`^([0-9]{1})$`),
				},
			},
			{
				Name: "comment",
				RecordRegexp: recordRegexp{
					pattern: regexp.MustCompile(`^comment$`),
				},
				ExactContain: exactContain{"comment"},
			},
			{
				Name:     "regexp|^type$",
				Required: true,
			},
		},
		ColumnsMap: map[int]Column{
			0: {
				Name:     "id",
				Required: true,
				RecordRegexp: recordRegexp{
					pattern: regexp.MustCompile(`^([0-9]{1})$`),
				},
			},
			1: {
				Name: "comment",
				RecordRegexp: recordRegexp{
					pattern: regexp.MustCompile(`^comment$`),
				},
				ExactContain: exactContain{"comment"},
			},
			2: {
				Name:     "regexp|^type$",
				Required: true,
			},
		},
	}
	assert.Equal(t, expected, s)
}
