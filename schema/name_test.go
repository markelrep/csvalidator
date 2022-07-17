package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestName_Regexp(t *testing.T) {
	cases := []struct {
		name     name
		expected string
	}{
		{
			name:     regexpPref + `^pattern$`,
			expected: `^pattern$`,
		},
		{
			name:     "blah",
			expected: "",
		},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, tc.name.Regexp())
	}
}

func TestName_IsRegexp(t *testing.T) {
	cases := []struct {
		name     name
		expected bool
	}{
		{
			name:     regexpPref + `^pattern$`,
			expected: true,
		},
		{
			name:     "blah",
			expected: false,
		},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, tc.name.IsRegexp())
	}
}
