package main

import (
	"context"
	"sync"
)

func newWorker(readyPool chan chan job, wg *sync.WaitGroup, number int) *worker {
	assignedJobQueue := make(chan job)
	return &worker{
		number:           number,
		assignedJobQueue: assignedJobQueue,
		wg:               wg,
		readyPool:        readyPool,
	}
}

type worker struct {
	number           int
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
				// fmt.Println("worker number:",w.number)
				job.Process(ctx)
				w.wg.Done()
			case <-ctx.Done():
				return
			}
		}
	}()
}
