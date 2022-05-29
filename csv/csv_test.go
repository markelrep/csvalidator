package csv

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadCSV(t *testing.T) {
	expected := [][]string{
		{"id", "comment"},
		{"1", "comment"},
		{"2", "comment"},
		{"3", "comment"},
		{"4", "comment"},
		{"5", "comment"},
	}

	cases := []struct {
		filePath string
		expected [][]string
	}{
		{
			filePath: "../samples/file.csv",
			expected: expected,
		},
		{
			filePath: "../samples/fileWithBOM.csv",
			expected: expected,
		},
	}

	for _, tc := range cases {
		records, err := ReadCSV(tc.filePath)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, records)
	}
}

func TestGetHeaders(t *testing.T) {
	r, err := ReadCSV("../samples/file.csv")
	assert.NoError(t, err)
	headers := GetHeaders(r)

	expected := map[string]struct{}{
		"id":      {},
		"comment": {},
	}

	headersJson, err := json.Marshal(headers)
	assert.NoError(t, err)
	expectedJson, err := json.Marshal(expected)
	assert.NoError(t, err)
	assert.JSONEq(t, string(expectedJson), string(headersJson))
}

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
		assert.Equal(t, tc.expected, IsCSV(tc.path))
	}
}
