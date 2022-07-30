package checklist

var ErrHeaderIsMissingTmpl = "%s required headers are missing: %v"
var ErrWrongColumnNameRegexpTmpl = "%s regexp pattern %s doesn't find in %s"
var ErrColumnNameAbsentSchemaTmpl = "%s column [%s] not found in schema"
var ErrWrongColumnNameExactTmpl = "%s column name is wrong, expected: %v, got: %v"
var ErrColumnIndexAbsentTmpl = "%s column # %d is absent in schema"
var ErrUnexpectedDataInCellTmpl = "%s line: %d, column number: %d, column name: %s. \"%v\" value is unexpected in this cell"
