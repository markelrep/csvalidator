package schema

import (
	"errors"
	"strings"
)

// exactContain represent slice of values which should be contained in column
type exactContain []string

// Contain checks that column contains values from schema and vice versa
func (c exactContain) Contain(values map[string]struct{}) error {
	containMap := make(map[string]struct{})
	var strBuilder strings.Builder

	for _, contain := range c {
		containMap[contain] = struct{}{}
		_, ok := values[contain]
		if !ok {
			errStr := contain + " is defined in schema, but absent in column"
			strBuilder.WriteString(errStr + "\n")
		}
	}

	for v := range values {
		_, ok := containMap[v]
		if !ok {
			errStr := v + " is not defined in schema, but exist in column"
			strBuilder.WriteString(errStr + "\n")
		}
	}

	str := strBuilder.String()
	if str != "" {
		return errors.New(str)
	}
	return nil
}

// IsNoOp return true if check is absent in schema
func (c exactContain) IsNoOp() bool {
	return len(c) == 0
}
