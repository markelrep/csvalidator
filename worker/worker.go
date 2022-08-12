package worker

import (
	"runtime"
	"sync"
)

type Job interface {
	Do(errCh chan error)
}

var defaultPoolSize = runtime.NumCPU()

// Pool of workers which pick up some job and Do it
type Pool struct {
	poolSize int
	jobQueue chan Job
	errChan  chan error
}

// NewPool creates new Pool
func NewPool(size int) *Pool {
	if size == 0 {
		size = defaultPoolSize
	}
	p := Pool{
		poolSize: size,
		jobQueue: make(chan Job),
		errChan:  make(chan error, size),
	}
	go p.run()
	return &p
}

// run starts workers
func (p *Pool) run() {
	var wg sync.WaitGroup
	wg.Add(p.poolSize)
	for i := 0; i < p.poolSize; i++ {
		go func(id int) {
			for j := range p.jobQueue {
				j.Do(p.errChan)
			}
			wg.Done()
		}(i + 1)
	}
	wg.Wait()
	close(p.errChan)
}

func (p *Pool) Errors() chan error {
	return p.errChan
}

func (p *Pool) Enqueue(j Job) {
	p.jobQueue <- j
}

// StopQueueingJob closes job Queue channel, used to notify workers there are no any job, so they can shutdown
func (p *Pool) StopQueueingJob() {
	close(p.jobQueue)
}
