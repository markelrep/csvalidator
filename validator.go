package csvalidator

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/markelrep/csvalidator/schema"
	"github.com/markelrep/csvalidator/worker"

	"github.com/markelrep/csvalidator/checklist"
	"github.com/markelrep/csvalidator/files"
)

// Validator stores csv files which to be validated and validate rules
type Validator struct {
	schema schema.Schema
	files  []files.File
	config Config
}

// NewValidator creates a new Validator
func NewValidator(pathFiles, pathSchema string, firstIsHeader bool) (Validator, error) {
	s, err := schema.Parse(pathSchema)
	if err != nil {
		return Validator{}, fmt.Errorf("failed create validator: %w", err)
	}

	file, err := files.NewFiles(pathFiles, firstIsHeader)
	if err != nil {
		return Validator{}, err
	}

	return Validator{
		schema: s,
		files:  file,
	}, nil
}

// NewValidatorWithConfig creates a new Validator with Config
func NewValidatorWithConfig(config Config) (Validator, error) {
	validator, err := NewValidator(config.FilePath, config.SchemaPath, config.FirstIsHeader)
	if err != nil {
		return Validator{}, err
	}
	validator.config = config
	return validator, nil
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
	file   files.File
	checks checklist.Checklist
}

func newJob(file files.File, checks checklist.Checklist) job {
	return job{
		file:   file,
		checks: checks,
	}
}

func (j job) Do() error {
	var errs error
	for _, check := range j.checks.List {
		err := check.Do(j.file)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
	}
	return errs
}
