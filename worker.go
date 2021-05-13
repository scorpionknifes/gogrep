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
	readyPool <- assignedJobQueue
	return &worker{
		assignedJobQueue: assignedJobQueue,
		wg:               wg,
	}
}

type worker struct {
	assignedJobQueue chan job
	wg               *sync.WaitGroup
}

func (w *worker) start(ctx context.Context) {
	go func() {
		for {
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
