package csvalidator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingColumns_Do(t *testing.T) {
	cases := []struct {
		filePath    string
		schemaPath  string
		expectedErr error
	}{
		{
			filePath:    "./samples/file.csv",
			schemaPath:  "./samples/schema.json",
			expectedErr: nil,
		},
		{
			filePath:    "./samples/fileOneColumn.csv",
			schemaPath:  "./samples/schema.json",
			expectedErr: errors.New("required headers are missing: [id]"),
		},
	}

	for _, tc := range cases {
		files, err := NewFiles(tc.filePath, true)
		assert.NoError(t, err)
		f := files[0]
		s, err := parseSchema(tc.schemaPath)
		assert.NoError(t, err)
		check := NewMissingColumn(s)
		err = check.Do(f)
		assert.Equal(t, tc.expectedErr, err)
	}
}

func TestColumnName_Do(t *testing.T) {
	cases := []struct {
		filePath    string
		schemaPath  string
		expectedErr error
	}{
		{
			filePath:    "./samples/file.csv",
			schemaPath:  "./samples/schema.json",
			expectedErr: nil,
		},
		{
			filePath:    "./samples/fileBadColumnName.csv",
			schemaPath:  "./samples/schema.json",
			expectedErr: errors.New("validation failed, column name is wrong, expected: comment, got: comments"),
		},
	}

	for _, tc := range cases {
		files, err := NewFiles(tc.filePath, true)
		assert.NoError(t, err)
		f := files[0]
		s, err := parseSchema(tc.schemaPath)
		assert.NoError(t, err)
		check := NewColumnName(s)
		err = check.Do(f)
		assert.Equal(t, tc.expectedErr, err)
	}
}
