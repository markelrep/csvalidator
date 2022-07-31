package checklist

import (
	"errors"
	"fmt"
	"testing"

	"github.com/markelrep/csvalidator/config"

	"github.com/stretchr/testify/require"

	"github.com/hashicorp/go-multierror"

	"github.com/markelrep/csvalidator/schema"

	"github.com/markelrep/csvalidator/files"

	"github.com/stretchr/testify/assert"
)

func TestMissingColumns_Do(t *testing.T) {
	cases := []struct {
		name        string
		filePath    string
		schemaPath  string
		expectedErr error
	}{
		{
			name:        "success",
			filePath:    "../samples/file.csv",
			schemaPath:  "../samples/schema.json",
			expectedErr: nil,
		},
		{
			name:        "error",
			filePath:    "../samples/fileOneColumn.csv",
			schemaPath:  "../samples/schema.json",
			expectedErr: errors.New("../samples/fileOneColumn.csv required headers are missing: [id]"),
		},
		{
			name:        "error",
			filePath:    "../samples/file.csv",
			schemaPath:  "../samples/schema_type.json",
			expectedErr: errors.New("../samples/file.csv required headers are missing: [regexp|^type$]"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := files.NewFiles(config.Config{FilePath: tc.filePath, FirstIsHeader: true})
			assert.NoError(t, err)
			f := file[0]
			s, err := schema.Parse(tc.schemaPath)
			assert.NoError(t, err)
			check := NewMissingColumn(s)
			err = check.Do(f)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestColumnName_Do(t *testing.T) {
	cases := []struct {
		name        string
		filePath    string
		schemaPath  string
		expectedErr func() (err error)
	}{
		{
			name:       "success case common name",
			filePath:   "../samples/file.csv",
			schemaPath: "../samples/schema.json",
			expectedErr: func() (err error) {
				return nil
			},
		},
		{
			name:       "bad common column name",
			filePath:   "../samples/fileBadColumnName.csv",
			schemaPath: "../samples/schema.json",
			expectedErr: func() (err error) {
				err = multierror.Append(err, errors.New("../samples/fileBadColumnName.csv column name is wrong, expected: comment, got: comments"))
				return err
			},
		},
		{
			name:       "success case regexp name",
			filePath:   "../samples/file.csv",
			schemaPath: "../samples/schema_regexp.json",
			expectedErr: func() (err error) {
				return nil
			},
		},
		{
			name:       "bad regexp column name",
			filePath:   "../samples/fileBadColumnName.csv",
			schemaPath: "../samples/schema_regexp.json",
			expectedErr: func() (err error) {
				err = multierror.Append(err, errors.New("../samples/fileBadColumnName.csv regexp pattern ^comment$ doesn't find in comments"))
				return err
			},
		},
		{
			name:       "absent in schema",
			filePath:   "../samples/file_redundant_column.csv",
			schemaPath: "../samples/schema.json",
			expectedErr: func() (err error) {
				return nil
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := files.NewFiles(config.Config{FilePath: tc.filePath, FirstIsHeader: true})
			assert.NoError(t, err)
			f := file[0]
			s, err := schema.Parse(tc.schemaPath)
			assert.NoError(t, err)
			check := NewColumnName(s)
			err = check.Do(f)
			assert.Equal(t, tc.expectedErr(), err)
		})
	}
}

func TestColumnRegexpMatch_Do(t *testing.T) {
	cases := []struct {
		name        string
		filePath    string
		schemaPath  string
		expectedErr func() error
	}{
		{
			name:       "success case",
			filePath:   "../samples/file.csv",
			schemaPath: "../samples/schema.json",
			expectedErr: func() error {
				return nil
			},
		},
		{
			name:       "with content error",
			filePath:   "../samples/fileContainWrongData.csv",
			schemaPath: "../samples/schema.json",
			expectedErr: func() error {
				var err error
				err = multierror.Append(err, fmt.Errorf(ErrUnexpectedDataInCellTmpl, "../samples/fileContainWrongData.csv", 5, 1, "id", 12))
				err = multierror.Append(err, fmt.Errorf(ErrUnexpectedDataInCellTmpl, "../samples/fileContainWrongData.csv", 5, 2, "comment", "blah"))
				return err
			},
		},
		{
			name:       "with schema error",
			filePath:   "../samples/file_redundant_column.csv",
			schemaPath: "../samples/schema.json",
			expectedErr: func() (err error) {
				return err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := files.NewFiles(config.Config{FilePath: tc.filePath, FirstIsHeader: true})
			assert.NoError(t, err)
			f := file[0]
			s, err := schema.Parse(tc.schemaPath)
			assert.NoError(t, err)
			check := NewColumnRegexpMatch(s)
			err = check.Do(f)
			assert.Equal(t, tc.expectedErr(), err)
		})
	}
}

func TestColumnExactContain(t *testing.T) {
	// TODO: this test can suddenly fail, because of map is using under `contains` therefore error order is not constant
	cases := []struct {
		name       string
		filePath   string
		schemaPath string
		expected   func() error
	}{
		{
			name:       "success case",
			filePath:   "../samples/file_contains.csv",
			schemaPath: "../samples/schema_contains.json",
			expected: func() error {
				return nil
			},
		},
		{
			name:       "with error",
			filePath:   "../samples/file_contains_err.csv",
			schemaPath: "../samples/schema_contains.json",
			expected: func() (err error) {
				err = multierror.Append(err, multierror.Prefix(fmt.Errorf("some value is defined in schema, but absent in column"), "../samples/file_contains_err.csv"))
				err = multierror.Append(err, multierror.Prefix(fmt.Errorf("value4 is defined in schema, but absent in column"), "../samples/file_contains_err.csv"))
				return err
			},
		},
		{
			name:       "with schema error",
			filePath:   "../samples/file_contains_redundant.csv",
			schemaPath: "../samples/schema_contains.json",
			expected: func() (err error) {
				return err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := files.NewFiles(config.Config{FilePath: tc.filePath, FirstIsHeader: true})
			require.NoError(t, err)
			f := file[0]
			s, err := schema.Parse(tc.schemaPath)
			require.NoError(t, err)
			check := NewColumnExactContain(s)
			err = check.Do(f)
			assert.Equal(t, tc.expected(), err)
		})
	}
}
