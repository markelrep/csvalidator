package csvalidator

import (
	"errors"
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/markelrep/csvalidator/config"

	"github.com/markelrep/csvalidator/schema"
	"github.com/markelrep/csvalidator/worker"

	"github.com/markelrep/csvalidator/checklist"
	"github.com/markelrep/csvalidator/files"
)

// Validator stores csv files which to be validated and validate rules
type Validator struct {
	schema schema.Schema
	files  files.Files
	config config.Config
}

// NewValidator creates a new Validator
func NewValidator(config config.Config) (Validator, error) {
	s, err := schema.Parse(config.SchemaPath)
	if err != nil {
		return Validator{}, fmt.Errorf("failed create validator: %w", err)
	}

	file, err := files.NewFiles(config)
	if err != nil {
		return Validator{}, err
	}

	return Validator{
		config: config,
		schema: s,
		files:  file,
	}, nil
}

// Validate checks files and expose errors.
// each file check runs concurrently
// errors return after all files are validated
func (v *Validator) Validate() error {
	wp := worker.NewPool(v.config.WorkerPoolSize)
	checks := checklist.NewChecklist(v.schema)
	if len(checks.List) == 0 {
		return errors.New("nothing to check")
	}

	for _, f := range v.files {
		j := newJob(f, checks)
		wp.Enqueue(j)
	}
	wp.StopQueueingJob()
	err := wp.Wait()
	if err != nil {
		return err
	}
	return nil
}

type job struct {
	file   *files.File
	checks checklist.Checklist
}

func newJob(file *files.File, checks checklist.Checklist) job {
	return job{
		file:   file,
		checks: checks,
	}
}

func (j job) Do() error {
	var errs error
	var wg sync.WaitGroup
	for _, check := range j.checks.List {
		wg.Add(1)
		go func(check checklist.Checker) {
			err := check.Do(j.file)
			if err != nil {
				errs = multierror.Append(errs, err)
			}
			wg.Done()
		}(check)
	}
	wg.Wait()
	return errs
}
