package checklist

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/go-multierror"

	"github.com/markelrep/csvalidator/schema"

	"github.com/markelrep/csvalidator/files"

	"github.com/stretchr/testify/assert"
)

func TestMissingColumns_Do(t *testing.T) {
	cases := []struct {
		filePath    string
		schemaPath  string
		expectedErr error
	}{
		{
			filePath:    "../samples/file.csv",
			schemaPath:  "../samples/schema.json",
			expectedErr: nil,
		},
		{
			filePath:    "../samples/fileOneColumn.csv",
			schemaPath:  "../samples/schema.json",
			expectedErr: errors.New("../samples/fileOneColumn.csv required headers are missing: [id]"),
		},
	}

	for _, tc := range cases {
		file, err := files.NewFiles(tc.filePath, true)
		assert.NoError(t, err)
		f := file[0]
		s, err := schema.Parse(tc.schemaPath)
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
			filePath:    "../samples/file.csv",
			schemaPath:  "../samples/schema.json",
			expectedErr: nil,
		},
		{
			filePath:    "../samples/fileBadColumnName.csv",
			schemaPath:  "../samples/schema.json",
			expectedErr: errors.New("../samples/fileBadColumnName.csv column name is wrong, expected: comment, got: comments"),
		},
		{
			filePath:    "../samples/file.csv",
			schemaPath:  "../samples/schema_regexp.json",
			expectedErr: nil,
		},
		{
			filePath:    "../samples/fileBadColumnName.csv",
			schemaPath:  "../samples/schema_regexp.json",
			expectedErr: errors.New("../samples/fileBadColumnName.csv regexp pattern ^comment$ doesn't find in comments"),
		},
	}

	for _, tc := range cases {
		file, err := files.NewFiles(tc.filePath, true)
		assert.NoError(t, err)
		f := file[0]
		s, err := schema.Parse(tc.schemaPath)
		assert.NoError(t, err)
		check := NewColumnName(s)
		err = check.Do(f)
		assert.Equal(t, tc.expectedErr, err)
	}
}

func TestNewColumnContains_Do(t *testing.T) {
	cases := []struct {
		filePath    string
		schemaPath  string
		expectedErr func() error
	}{
		{
			filePath:   "../samples/file.csv",
			schemaPath: "../samples/schema.json",
			expectedErr: func() error {
				return nil
			},
		},
		{
			filePath:   "../samples/fileContainWrongData.csv",
			schemaPath: "../samples/schema.json",
			expectedErr: func() error {
				var err error
				err = multierror.Append(err, fmt.Errorf(ErrUnexpectedDataInCellTmpl, "../samples/fileContainWrongData.csv", 5, 1, "id", 12))
				err = multierror.Append(err, fmt.Errorf(ErrUnexpectedDataInCellTmpl, "../samples/fileContainWrongData.csv", 5, 2, "comment", "blah"))
				return err
			},
		},
	}

	for _, tc := range cases {
		file, err := files.NewFiles(tc.filePath, true)
		assert.NoError(t, err)
		f := file[0]
		s, err := schema.Parse(tc.schemaPath)
		assert.NoError(t, err)
		check := NewColumnContains(s)
		err = check.Do(f)
		assert.Equal(t, tc.expectedErr(), err)
	}
}
