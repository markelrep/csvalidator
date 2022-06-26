package checklist

import (
	"fmt"
	"regexp"

	"github.com/markelrep/csvalidator/schema"

	"github.com/markelrep/csvalidator/files"
)

// MissingColumn checks that all appropriate columns are presented in a file
type MissingColumn struct {
	schema schema.Schema
}

// NewMissingColumn creates new MissingColumn check
func NewMissingColumn(schema schema.Schema) MissingColumn {
	return MissingColumn{schema: schema}
}

// Do is doing the check of MissingColumn
func (mc MissingColumn) Do(f files.File) error {
	missing := make([]string, 0, f.HeadersCount())
	for _, header := range mc.schema.Columns {
		if !f.HasHeader(header.Name.String()) && header.Required {
			missing = append(missing, header.Name.String())
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("%s required headers are missing: %v", f.Path(), missing)
	}
	return nil
}

// ColumnName checks that all columns have a name which was defined in schema
type ColumnName struct {
	schema schema.Schema
}

// NewColumnName creates a new ColumnName check
func NewColumnName(schema schema.Schema) ColumnName {
	return ColumnName{schema: schema}
}

// Do is doing the check of ColumnName
func (cn ColumnName) Do(f files.File) error {
	for i, row := range f.Records {
		if i == 0 && f.FirstIsHeader() {
			for i := range row {
				expected := cn.schema.Columns[i].Name
				got := row[i]
				if expected.IsRegexp() {
					if r, err := regexp.Compile(expected.String()); err == nil {
						if !r.MatchString(got) {
							return fmt.Errorf("%s regexp pattern %s doesn't find in %s", f.Path(), r.String(), got)
						}
					}
					continue
				}
				if expected.String() != got {
					return fmt.Errorf("%s column name is wrong, expected: %v, got: %v", f.Path(), expected, got)
				}
			}
		}
	}
	return nil
}
