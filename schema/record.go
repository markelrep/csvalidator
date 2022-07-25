package schema

import (
	"bytes"
	"regexp"
)

// recordRegexp is field in schema which uses to define validation type and validation rule for data in cell
type recordRegexp struct {
	pattern *regexp.Regexp
}

// UnmarshalJSON custom unmarshaler
// This filed can contain list of accepted values and regexp in schema
// So we need to have custom unmarshaler to define what type is used and set it into struct
func (r *recordRegexp) UnmarshalJSON(data []byte) error {
	if bytes.EqualFold(data, []byte("null")) {
		return nil
	}

	exp, err := regexp.Compile(string(bytes.Trim(data, `"`)))
	if err != nil {
		return err
	}
	r.pattern = exp
	return nil
}

// Match is checking that the value from cell is meet validation rules from schema
func (r recordRegexp) Match(value string) bool {
	return r.pattern.MatchString(value)
}

// IsNoOp returns true in case this check was empty in schema, and we don't need to process this check
func (r recordRegexp) IsNoOp() bool {
	return r.pattern == nil
}
