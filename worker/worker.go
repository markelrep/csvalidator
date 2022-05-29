package worker

import (
	"runtime"
	"sync/atomic"

	"github.com/hashicorp/go-multierror"
)

type Job interface {
	Do() error
}

var defaultPoolSize = runtime.NumCPU()

// Pool of workers which pick up some job and Do it
type Pool struct {
	poolSize       int
	jobQueue       chan Job
	runningWorkers int32
	errChan        chan error
}

// NewPool creates new Pool
func NewPool(size int) *Pool {
	if size == 0 {
		size = defaultPoolSize
	}
	p := Pool{
		poolSize:       size,
		jobQueue:       make(chan Job),
		runningWorkers: 0,
		errChan:        make(chan error, size),
	}
	go p.run()
	return &p
}

// run starts workers
func (p *Pool) run() {
	atomic.AddInt32(&p.runningWorkers, int32(p.poolSize))
	for i := 0; i < p.poolSize; i++ {
		go func(id int) {
			for j := range p.jobQueue {
				err := j.Do()
				if err != nil {
					p.errChan <- err
				}
			}
			atomic.AddInt32(&p.runningWorkers, -1)
		}(i + 1)
	}
}

// done returns chan with empty struct when all workers are done
func (p *Pool) done() <-chan struct{} {
	v := atomic.LoadInt32(&p.runningWorkers)
	if v == 0 {
		d := make(chan struct{}, 1)
		d <- struct{}{}
		return d
	}
	return nil
}

func (p *Pool) Enqueue(j Job) {
	p.jobQueue <- j
}

// StopQueueingJob closes job Queue channel, used to notify workers there are no any job, so they can shutdown
func (p *Pool) StopQueueingJob() {
	close(p.jobQueue)
}

// Wait blocks further execution while all jobs are not done, collects errors and then return these errors
func (p *Pool) Wait() error {
	var errs error
	for {
		select {
		case <-p.done():
			if errs != nil {
				return errs
			}
			return nil
		case err := <-p.errChan:
			if err != nil {
				errs = multierror.Append(errs, err)
			}
		default:
			continue
		}
	}
}
