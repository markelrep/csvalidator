package schema

import "strings"

// name is column name in json schema
type name string

func (n name) String() string {
	return string(n)
}

func (n name) Regexp() string {
	if n.IsRegexp() {
		return strings.TrimPrefix(n.String(), regexpPref)
	}
	return ""
}

// IsRegexp returns true if name is regexp and false if just name
func (n name) IsRegexp() bool {
	return strings.HasPrefix(n.String(), regexpPref)
}

// IsNoOp returns true in case this check was empty in schema, and we don't need to process this check
func (n name) IsNoOp() bool {
	return n.String() == ""
}
