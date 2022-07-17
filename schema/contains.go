package schema

import (
	"bytes"
	"encoding/json"
	"regexp"
)

// containType is a kind of data in schema
type containType int

const (
	unknown       containType = iota
	regexpPattern containType = iota
	list          containType = iota
)

// contains is field in schema which uses to define validation type and validation rule for data in cell
type contains struct {
	kind    containType
	list    map[any]struct{}
	pattern *regexp.Regexp
}

// UnmarshalJSON custom unmarshaler
// This filed can contain list of accepted values and regexp in schema
// So we need to have custom unmarshaler to define what type is used and set it into struct
func (c *contains) UnmarshalJSON(data []byte) error {
	if bytes.EqualFold(data, []byte("null")) {
		return nil
	}
	if len(data) < 2 {
		return nil
	}
	isSlice := data[0] == byte(91) && data[len(data)-1] == byte(93)
	if isSlice {
		c.list = make(map[any]struct{})
		var slice []any
		err := json.Unmarshal(data, &slice)
		if err != nil {
			return err
		}
		for _, l := range slice {
			c.list[l] = struct{}{}
		}
		c.kind = list
		return nil
	}
	if len(data) < 7 {
		return nil
	}
	isRegexp := bytes.EqualFold(bytes.Trim(data, `"`)[:7], []byte("regexp|"))
	if isRegexp {
		c.kind = regexpPattern
		r, err := regexp.Compile(string(bytes.Trim(data[8:], `"`)))
		if err != nil {
			return err
		}
		c.pattern = r
		return nil
	}
	return nil
}

// Contain is checking that the value from cell is meet validation rules from schema
func (c contains) Contain(value any) bool {
	switch c.kind {
	case list:
		_, ok := c.list[value]
		return ok
	case regexpPattern:
		v, ok := value.(string)
		if !ok {
			// TODO: log the error
			return false
		}
		return c.pattern.MatchString(v)
	}
	return false
}

// IsNoOp returns true in case this check was empty in schema, and we don't need to process this check
func (c contains) IsNoOp() bool {
	return c.kind == unknown
}
