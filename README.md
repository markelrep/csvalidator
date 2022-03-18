# CSValidator

## Overview
The main purpose of project to provide opportunity validate CSV (Comma Separated Values) file by pre-defined schema in JSON formatted file

## Usage

```go
validator := csvalidator.NewValidator("./path/to/csv/file", "./path/to/json/schema", true)
if err := validator.Validate(); err != nil {
	// handle error
}
```