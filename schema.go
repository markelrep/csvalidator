package csvalidator

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type column struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
	// Values validation
	DataType string   `json:"dataType"`
	Contains []string `json:"contains"`
}

func (c *column) Validate(value string) error {
	if err := c.validateContains(value); err != nil {
		return err
	}
	return nil
}

func (c *column) validateContains(value string) error {
	if len(c.Contains) == 0 {
		return nil
	}
	if !stringInSlice(value, c.Contains) {
		return fmt.Errorf("expected one of: '%v'. Actual is '%v'", strings.Join(c.Contains, ", "), value)
	}
	return nil
}

func stringInSlice(value string, values []string) bool {
	for v := range values {
		if value == values[v] {
			return true
		}
	}
	return false
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
