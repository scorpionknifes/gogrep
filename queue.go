package main

import (
	"context"
	"log"
	"sync"
)

type jobQueue struct {
	internalQueue     chan job
	workers           []*worker
	dispatcherStopped *sync.WaitGroup
	workersStopped    *sync.WaitGroup
}

func newJobQueue(maxWorkers int) *jobQueue {
	workersStopped := sync.WaitGroup{}
	workers := make([]*worker, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		workers[i] = newWorker(readyPool, &workersStopped)
	}
	return &jobQueue{
		internalQueue:     make(chan job),
		workers:           workers,
		dispatcherStopped: &sync.WaitGroup{},
		workersStopped:    &workersStopped,
	}
}

func (q *jobQueue) start(ctx context.Context) {
	for i := 0; i < len(q.workers); i++ {
		q.workers[i].start(ctx)
	}
	go q.dispatch(ctx)
}

// append to array of jobs
// receive from internal
//

func (q *jobQueue) dispatch(ctx context.Context) {
	q.dispatcherStopped.Add(1)
	for {
		select {
		case job := <-q.internalQueue:
			log.Println("submit2")
			workerChannel := <-q.readyPool
			log.Println("submit3")
			workerChannel <- job
			log.Println("submit4")
		case <-ctx.Done():
			q.workersStopped.Wait()
			q.dispatcherStopped.Done()
			return
		}
	}
}

func (q *jobQueue) submit(ctx context.Context, job job) {
	select {
	case q.internalQueue <- job:
		log.Println("submit")
	case <-ctx.Done():
		q.stop()
	}
}

func (q *jobQueue) stop() {
	q.dispatcherStopped.Wait()
}
