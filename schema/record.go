package schema

import (
	"bytes"
	"regexp"
)

// contains is field in schema which uses to define validation type and validation rule for data in cell
type contains struct {
	pattern *regexp.Regexp
}

// UnmarshalJSON custom unmarshaler
// This filed can contain list of accepted values and regexp in schema
// So we need to have custom unmarshaler to define what type is used and set it into struct
func (c *contains) UnmarshalJSON(data []byte) error {
	if bytes.EqualFold(data, []byte("null")) {
		return nil
	}

	r, err := regexp.Compile(string(bytes.Trim(data, `"`)))
	if err != nil {
		return err
	}
	c.pattern = r
	return nil
}

// Contain is checking that the value from cell is meet validation rules from schema
func (c contains) Contain(value string) bool {
	return c.pattern.MatchString(value)
}

// IsNoOp returns true in case this check was empty in schema, and we don't need to process this check
func (c contains) IsNoOp() bool {
	return c.pattern.String() == ""
}
