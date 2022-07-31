package checklist

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"

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
func (mc MissingColumn) Do(f *files.File) error {
	missing := make([]string, 0, f.HeadersCount())
	for i, column := range mc.schema.Columns {
		if column.Required {
			if !column.Name.IsRegexp() {
				if !f.HasHeader(column.Name.String()) {
					missing = append(missing, column.Name.String())
				}
				continue
			}
			if f.HeadersCount() >= i {
				missing = append(missing, column.Name.String())
				continue
			}
			if ok, _ := regexp.Match(column.Name.Regexp(), []byte(f.Headers()[i])); !ok {
				missing = append(missing, column.Name.String())
			}
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
func (cn ColumnName) Do(f *files.File) (errs error) {
	for i, got := range f.Headers() {
		column, ok := cn.schema.GetColumn(i)
		if !ok {
			continue
		}
		expected := column.Name
		if expected.IsNoOp() {
			continue
		}
		if expected.IsRegexp() {
			if r, err := regexp.Compile(expected.Regexp()); err == nil {
				if !r.MatchString(got) {
					errs = multierror.Append(errs, fmt.Errorf(ErrWrongColumnNameRegexpTmpl, f.Path(), r.String(), got))
					continue
				}
			}
			continue
		}
		if expected.String() != got {
			errs = multierror.Append(errs, fmt.Errorf(ErrWrongColumnNameExactTmpl, f.Path(), expected, got))
			continue
		}
	}
	return errs
}

// ColumnRegexpMatch checks data in a column on regexp match
type ColumnRegexpMatch struct {
	schema schema.Schema
}

// NewColumnRegexpMatch creates new ColumnRegexpMatch
func NewColumnRegexpMatch(schema schema.Schema) ColumnRegexpMatch {
	return ColumnRegexpMatch{schema: schema}
}

// Do check of ColumnRegexpMatch
func (cc ColumnRegexpMatch) Do(f *files.File) (err error) {
	for row := range f.Stream() {
		for j, record := range row.Data {
			column, ok := cc.schema.GetColumn(j)
			if !ok {
				continue
			}
			regexpPattern := column.RecordRegexp
			if regexpPattern.IsNoOp() {
				continue
			}
			if !regexpPattern.Match(record) {
				err = multierror.Append(err,
					fmt.Errorf(ErrUnexpectedDataInCellTmpl, f.Path(), row.Index+1, j+1, cc.schema.Columns[j].Name, record),
				)
			}
		}
	}
	return err
}

// ColumnExactContain checks data in column contains strings which specified in schema
type ColumnExactContain struct {
	schema schema.Schema
}

// NewColumnExactContain creates new ColumnExactContain
func NewColumnExactContain(schema schema.Schema) ColumnExactContain {
	return ColumnExactContain{schema: schema}
}

// Do check of ColumnExactContain
func (c ColumnExactContain) Do(f *files.File) (err error) {
	data := make(map[string]map[string]struct{}) // todo: two same columns
	indexes := make(map[int]string)

	for row := range f.Stream() {
		if row.Index == 1 {
			for i := 0; i < len(row.Data); i++ {
				column, ok := c.schema.GetColumn(i)
				if !ok {
					continue
				}
				if column.ExactContain.IsNoOp() {
					continue
				}
				id := uuid.NewString()
				indexes[i] = id
				data[id] = map[string]struct{}{}
			}
		}
		for j, cell := range row.Data {
			column, ok := c.schema.GetColumn(j)
			if !ok {
				continue
			}
			if column.ExactContain.IsNoOp() {
				continue
			}
			data[indexes[j]][cell] = struct{}{}
		}
	}
	for index, key := range indexes {
		e := c.schema.Columns[index].ExactContain.Contain(data[key])
		if e != nil {
			err = multierror.Append(err, multierror.Prefix(e, f.Path()))
		}
	}
	return err
}
