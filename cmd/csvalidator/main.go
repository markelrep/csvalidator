package main

import (
	"fmt"
	"log"

	"github.com/markelrep/csvalidator/config"

	"github.com/markelrep/csvalidator"

	flags "github.com/umputun/go-flags"
)

func main() {
	var cfg config.Config
	_, err := flags.Parse(&cfg)
	if err != nil {
		log.Fatalln(err)
	}

	if cfg.FilePath == "" || cfg.SchemaPath == "" {
		log.Fatalln("schema and path flag are required")
	}

	commaRune := []rune(cfg.CommaString)
	if len(commaRune) > 1 {
		log.Fatalln("separator is wrong", cfg.CommaString)
	}
	if cfg.CommaString != "" {
		cfg.Comma = commaRune[0]
	}

	fmt.Printf("ATTENTION! First is header was set as %v\n", cfg.FirstIsHeader)

	validator, err := csvalidator.NewValidator(cfg)
	if err != nil {
		log.Fatalln("create validator ", err)
	}
	err = validator.Validate()
	if err != nil {
		log.Fatalln("validation error ", err)
	}
}
