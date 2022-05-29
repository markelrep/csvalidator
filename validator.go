package csvalidator

import (
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/markelrep/csvalidator/schema"
	"github.com/markelrep/csvalidator/worker"

	"github.com/markelrep/csvalidator/checklist"
	"github.com/markelrep/csvalidator/files"
)

type Validator struct {
	schema schema.Schema
	files  []files.File
}

func NewValidator(path, schemaPath string, firstHeader bool) (Validator, error) {
	s, err := schema.Parse(schemaPath)
	if err != nil {
		return Validator{}, fmt.Errorf("failed create validator: %w", err)
	}

	file, err := files.NewFiles(path, firstHeader)
	if err != nil {
		return Validator{}, err
	}

	return Validator{
		schema: s,
		files:  file,
	}, nil
}

func (v *Validator) Validate() error {
	wp := worker.NewPool(0) // TODO: configurable
	checks := checklist.NewChecklist(v.schema)

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
