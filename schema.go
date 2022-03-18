package csvalidator

import (
	"encoding/json"
	"fmt"
	"os"
)

type column struct {
	Name     string `json:"name"`
	DataType string `json:"dataType"`
	Required bool   `json:"required"`
}

type schema struct {
	Columns []column `json:"columns"`
}

func parseSchema(schemaPath string) (schema, error) {
	file, err := os.ReadFile(schemaPath)
	if err != nil {
		return schema{}, fmt.Errorf("failed read schema by path %v: %w", schemaPath, err)
	}
	var s schema
	err = json.Unmarshal(file, &s)
	if err != nil {
		return schema{}, fmt.Errorf("failed unmarshal schema from %v: %w", schemaPath, err)
	}
	return s, nil
}
