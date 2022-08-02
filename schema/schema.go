package schema

import (
	"encoding/json"
	"fmt"
	"os"
)

const regexpPref = "regexp|"

type ColumnsMap map[int]Column

// Column recordRegexp all appropriate data which will be needed to validate column in csv
type Column struct {
	Name         name         `json:"name"`
	Required     bool         `json:"required"`
	RecordRegexp recordRegexp `json:"record_regexp"`
	ExactContain exactContain `json:"exact_contain"`
}

// Schema recordRegexp suite of information by which file validates
type Schema struct {
	Columns    []Column   `json:"columns"`
	ColumnsMap ColumnsMap `json:"-"`
}

func Parse(schemaPath string) (Schema, error) {
	file, err := os.ReadFile(schemaPath)
	if err != nil {
		return Schema{}, fmt.Errorf("failed read Schema by path %v: %w", schemaPath, err)
	}
	var s Schema
	err = json.Unmarshal(file, &s)
	if err != nil {
		return Schema{}, fmt.Errorf("failed unmarshal Schema from %v: %w", schemaPath, err)
	}
	s.ColumnsMap = make(ColumnsMap, len(s.Columns))
	for i, v := range s.Columns {
		s.ColumnsMap[i] = v
	}
	return s, nil
}

func (s Schema) GetColumn(index int) (Column, bool) {
	column, ok := s.ColumnsMap[index]
	return column, ok
}
