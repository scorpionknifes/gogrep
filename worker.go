package main

import (
	"context"
	"sync"
)

type job interface {
	Process(ctx context.Context)
}

func newWorker(readyPool chan chan job, wg *sync.WaitGroup) *worker {
	assignedJobQueue := make(chan job)
	return &worker{
		assignedJobQueue: assignedJobQueue,
		wg:               wg,
		readyPool:        readyPool,
	}
}

type worker struct {
	assignedJobQueue chan job
	wg               *sync.WaitGroup
	readyPool        chan chan job
}

func (w *worker) start(ctx context.Context) {
	go func() {
		for {
			w.readyPool <- w.assignedJobQueue
			select {
			case job := <-w.assignedJobQueue:
				job.Process(ctx)
				w.wg.Done()
			case <-ctx.Done():
				return
			}
		}
	}()
}
