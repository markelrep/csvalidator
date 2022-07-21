package schema

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestContains_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		field    string
		expected contains
	}{
		{
			field:    `null`,
			expected: contains{},
		},
		{
			field: `"^pattern$"`,
			expected: contains{
				pattern: regexp.MustCompile(`^pattern$`),
			},
		},
		{
			field: `""`,
			expected: contains{
				pattern: regexp.MustCompile(""),
			},
		},
	}

	for _, tc := range cases {
		var c contains
		err := json.Unmarshal([]byte(tc.field), &c)
		require.NoError(t, err)
		assert.Equal(t, tc.expected, c)
	}
}

func TestContains_Contain(t *testing.T) {
	cases := []struct {
		value    string
		contains contains
		expected bool
	}{
		{
			value: "value1",
			contains: contains{
				pattern: regexp.MustCompile(`^value1$`),
			},
			expected: true,
		},
		{
			value: "value2",
			contains: contains{
				pattern: regexp.MustCompile(`^value1$`),
			},
			expected: false,
		},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, tc.contains.Contain(tc.value))
	}
}
