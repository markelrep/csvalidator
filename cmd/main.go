package main

import (
	"flag"
	"log"

	"github.com/markelrep/csvalidator"
)

func main() {
	path := flag.String("p", "", "path to csv files")
	schema := flag.String("s", "", "path to schema")
	flag.Parse()
	if *schema == "" || *path == "" {
		log.Fatalln("schema and path flag are required")
	}

	validator, err := csvalidator.NewValidatorWithConfig(csvalidator.Config{
		FilePath:       *path,
		FirstIsHeader:  true,
		SchemaPath:     *schema,
		WorkerPoolSize: 0,
	})
	if err != nil {
		log.Fatalln("create validator ", err)
	}
	err = validator.Validate()
	if err != nil {
		log.Fatalln("validation error ", err)
	}
}
