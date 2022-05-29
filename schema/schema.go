package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type name string

func (n name) String() string {
	return string(n)
}

func (n *name) IsRegexp() bool {
	pref := "regexp|"
	isRegexp := strings.HasPrefix(n.String(), pref)
	if isRegexp {
		*n = name(strings.TrimPrefix(n.String(), pref))
	}
	return isRegexp
}

type Column struct {
	Name     name   `json:"name"`
	DataType string `json:"dataType"`
	Required bool   `json:"required"`
}

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
