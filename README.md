# CSValidator

CSValidator is a tool for validation of .csv files using JSON schema

## Quick Start

### Install
`go get -u github.com/markelrep/csvalidator`

### Schema

```json
{
  "columns":[
    {
      "name": "id",
      "required": true
    },
    {
      "name": "comment",
      "required": false
    }
  ]
}
```
`columns` is array of objects with validation rules for each columns in .csv file

`name` of column, which should be the same as in .csv file otherwise validation will be failed. This field also supports regexp. [Example](https://github.com/markelrep/csvalidator/blob/master/samples/schema_regexp.json)

`required` true means that this column is required to exist in file, false that isn't required

### Usage
```go
package main
import "github.com/markelrep/csvalidator"

func main() {
	validator, err := csvalidator.NewValidator(csvalidator.Config{
		FilePath:       "./path/to/csv/files",
		FirstIsHeader:  true,
		SchemaPath:     "./path/to/json/schema",
		WorkerPoolSize: 0,
	})
	if err != nil {
		// handle error
    }
	if err := validator.Validate(); err != nil {
		// handle error
	}
}
```