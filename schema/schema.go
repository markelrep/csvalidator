package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// name is column name in json schema
type name string

func (n name) String() string {
	return string(n)
}

// IsRegexp returns true if name is regexp and false if just name
func (n *name) IsRegexp() bool {
	pref := "regexp|"
	isRegexp := strings.HasPrefix(n.String(), pref)
	if isRegexp {
		*n = name(strings.TrimPrefix(n.String(), pref))
	}
	return isRegexp
}

// Column contains all appropriate data which will be needed to validate column in csv
type Column struct {
	Name     name   `json:"name"`
	DataType string `json:"dataType"`
	Required bool   `json:"required"`
}

// Schema contains suite of information by which file validates
type Schema struct {
	Columns []Column `json:"columns"`
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
	return s, nil
}
