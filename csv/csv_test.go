package csv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCSV(t *testing.T) {
	cases := []struct {
		path     string
		expected bool
	}{
		{
			path:     "../folder/file.csv",
			expected: true,
		},
		{
			path:     "file.csv",
			expected: true,
		},
		{
			path:     "../folder/file.txt",
			expected: false,
		},
		{
			path:     "e.txt",
			expected: false,
		},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.expected, isCSV(tc.path))
	}
}
