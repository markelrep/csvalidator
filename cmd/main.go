package main

import (
	"flag"
	"log"

	"github.com/markelrep/csvalidator"
)

func main() {
	path := flag.String("p", "", "path to csv files")
	schema := flag.String("s", "", "path to schema")
	lazyQuotes := flag.Bool("l", false, "lazy quotes")
	comma := flag.String("c", ",", "separator")
	flag.Parse()
	if *schema == "" || *path == "" {
		log.Fatalln("schema and path flag are required")
	}

	commaRune := []rune(*comma)
	if len(commaRune) > 1 {
		log.Fatalln("separator is wrong", *comma)
	}

	validator, err := csvalidator.NewValidator(csvalidator.Config{
		FilePath:       *path,
		FirstIsHeader:  true,
		SchemaPath:     *schema,
		LazyQuotes:     *lazyQuotes,
		Comma:          commaRune[0],
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
