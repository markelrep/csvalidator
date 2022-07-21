package checklist

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-multierror"

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
		return fmt.Errorf(ErrHeaderIsMissingTmpl, f.Path(), missing)
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
				if expected.IsNoOp() {
					continue
				}
				got := row[i]
				if expected.IsRegexp() {
					if r, err := regexp.Compile(expected.Regexp()); err == nil {
						if !r.MatchString(got) {
							return fmt.Errorf(ErrWrongColumnNameRegexpTmpl, f.Path(), r.String(), got)
						}
					}
					continue
				}
				if expected.String() != got {
					return fmt.Errorf(ErrWrongColumnNameExactTmpl, f.Path(), expected, got)
				}
			}
		}
	}
	return nil
}

type ColumnContains struct {
	schema schema.Schema
}

func NewColumnContains(schema schema.Schema) ColumnContains {
	return ColumnContains{schema: schema}
}

func (cc ColumnContains) Do(f files.File) (err error) {
	for i, row := range f.Records {
		if i == 0 && f.FirstIsHeader() {
			continue
		}
		for j, record := range row {
			contains := cc.schema.Columns[j].RecordRegexp
			if contains.IsNoOp() {
				continue
			}
			if !contains.Contain(record) { // TODO: maybe need to return error here instead of bool
				err = multierror.Append(err,
					fmt.Errorf(ErrUnexpectedDataInCellTmpl, f.Path(), i+1, j+1, cc.schema.Columns[j].Name, record),
				)
			}
		}
	}
	return err
}
