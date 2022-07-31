package schema

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

// exactContain represent slice of values which should be contained in column
type exactContain []string

// Contain checks that column contains values from schema and vice versa
func (c exactContain) Contain(values map[string]struct{}) (err error) {
	containMap := make(map[string]struct{})

	for _, contain := range c {
		containMap[contain] = struct{}{}
		_, ok := values[contain]
		if !ok {
			err = multierror.Append(err, fmt.Errorf("%s is defined in schema, but absent in column", contain))
		}
	}

	for v := range values {
		_, ok := containMap[v]
		if !ok {
			err = multierror.Append(err, fmt.Errorf("%s is not defined in schema, but exist in column", v))
		}
	}
	return err
}

// IsNoOp return true if check is absent in schema
func (c exactContain) IsNoOp() bool {
	return len(c) == 0
}
