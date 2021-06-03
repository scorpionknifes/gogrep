package main

import (
	"context"
	"sync"
)

type jobQueue struct {
	internalQueue chan job
	workers       []*worker
	readyPool     chan chan job
	wg            *sync.WaitGroup
}

func newJobQueue(maxWorkers int) *jobQueue {
	workers := make([]*worker, maxWorkers)
	readyPool := make(chan chan job, maxWorkers)
	wg := sync.WaitGroup{}

	for i := 0; i < maxWorkers; i++ {
		workers[i] = newWorker(readyPool, &wg)
	}

	return &jobQueue{
		internalQueue: make(chan job),
		workers:       workers,
		readyPool:     readyPool,
		wg:            &wg,
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

// func (q *jobQueue) dispatch(ctx context.Context) {
// 	for {
// 		select {
// 		case job := <-q.internalQueue:
// 			workerChannel := <-q.readyPool
// 			workerChannel <- job
// 		case <-ctx.Done():
// 			return
// 		}
// 	}
// }

func (q *jobQueue) dispatch(ctx context.Context) {
	var jobChan chan job
	var jobsIn []job
	var workerChanIn chan chan job
	for {
		select {
		case job := <-q.internalQueue:
			jobsIn = append(jobsIn, job)
			workerChanIn = q.readyPool
		case workerChannel := <-workerChanIn:
			jobChan = workerChannel
			workerChanIn = nil
		case (len(jobsIn) != 0) && (jobChan <- jobsIn[0]):
			jobsIn = jobsIn[1:]
		case <-ctx.Done():
			return
		}
	}
}

func (q *jobQueue) submit(ctx context.Context, job job) {
	select {
	case q.internalQueue <- job:
		q.wg.Add(1)
	case <-ctx.Done():
		return
	}
}
