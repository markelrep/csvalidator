package schema

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type contains []string

func (c contains) Contain(values map[string]struct{}) (err error) {
	containMap := make(map[string]struct{})

	for _, contain := range c {
		containMap[contain] = struct{}{}
		_, ok := values[contain]
		if !ok {
			err = multierror.Append(err, fmt.Errorf("%s is defined in schema, but absent in column", contain))
		}
	}

	for v, _ := range values {
		_, ok := containMap[v]
		if !ok {
			err = multierror.Append(err, fmt.Errorf("%s is not defined in schema, but exist in column", v))
		}
	}
	return err
}

func (c contains) IsNoOp() bool {
	return c == nil || len(c) == 0
}
