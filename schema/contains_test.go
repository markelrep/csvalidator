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
			field: `["data","foo","bar"]`,
			expected: contains{
				kind: list,
				list: map[any]struct{}{"data": {}, "foo": {}, "bar": {}},
			},
		},
		{
			field: `[1,2,3]`,
			expected: contains{
				kind: list,
				list: map[any]struct{}{1.0: {}, 2.0: {}, 3.0: {}},
			},
		},
		{
			field: `[]`,
			expected: contains{
				kind: list,
				list: map[any]struct{}{},
			},
		},
		{
			field: `"regexp|^pattern$"`,
			expected: contains{
				pattern: regexp.MustCompile(`^pattern$`),
				kind:    regexpPattern,
			},
		},
		{
			field:    `""`,
			expected: contains{},
		},
	}

	for _, tc := range cases {
		var contains contains
		err := json.Unmarshal([]byte(tc.field), &contains)
		require.NoError(t, err)
		assert.Equal(t, tc.expected, contains)
	}
}

func TestContains_Contain(t *testing.T) {
	cases := []struct {
		value    string
		contains contains
		expected bool
	}{
		{
			value: "value2",
			contains: contains{
				kind: list,
				list: map[any]struct{}{"value1": {}, "value2": {}},
			},
			expected: true,
		},
		{
			value: "value1",
			contains: contains{
				kind: list,
				list: map[any]struct{}{"value1": {}, "value2": {}},
			},
			expected: true,
		},
		{
			value: "value3",
			contains: contains{
				kind: list,
				list: map[any]struct{}{"value1": {}, "value2": {}},
			},
			expected: false,
		},
		{
			value: "value1",
			contains: contains{
				kind:    regexpPattern,
				pattern: regexp.MustCompile(`^value1$`),
			},
			expected: true,
		},
		{
			value: "value2",
			contains: contains{
				kind:    regexpPattern,
				pattern: regexp.MustCompile(`^value1$`),
			},
			expected: false,
		},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, tc.contains.Contain(tc.value))
	}
}
