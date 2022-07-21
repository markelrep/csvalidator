package schema

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestRecordRegexp_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		field    string
		expected recordRegexp
	}{
		{
			field:    `null`,
			expected: recordRegexp{},
		},
		{
			field: `"^pattern$"`,
			expected: recordRegexp{
				pattern: regexp.MustCompile(`^pattern$`),
			},
		},
		{
			field: `""`,
			expected: recordRegexp{
				pattern: regexp.MustCompile(""),
			},
		},
	}

	for _, tc := range cases {
		var c recordRegexp
		err := json.Unmarshal([]byte(tc.field), &c)
		require.NoError(t, err)
		assert.Equal(t, tc.expected, c)
	}
}

func TestRecordRegexp_Contain(t *testing.T) {
	cases := []struct {
		value    string
		record   recordRegexp
		expected bool
	}{
		{
			value: "value1",
			record: recordRegexp{
				pattern: regexp.MustCompile(`^value1$`),
			},
			expected: true,
		},
		{
			value: "value2",
			record: recordRegexp{
				pattern: regexp.MustCompile(`^value1$`),
			},
			expected: false,
		},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, tc.record.Contain(tc.value))
	}
}

func TestRecordRegexp_IsNoOp(t *testing.T) {
	cases := []struct {
		record   recordRegexp
		expected bool
	}{
		{
			record:   recordRegexp{pattern: regexp.MustCompile("")},
			expected: false,
		},
		{
			record:   recordRegexp{},
			expected: true,
		},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, tc.record.IsNoOp())
	}
}
