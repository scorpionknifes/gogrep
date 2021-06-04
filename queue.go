package main

import (
	"context"
	"sync"
)

type jobQueue struct {
	internalQueue chan job
	jobChan       chan job
	workers       []*worker
	jobsIn        []job
	readyPool     chan chan job
	workerChanIn  chan chan job
	wg            *sync.WaitGroup
	started       bool
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

func (q *jobQueue) dispatch(ctx context.Context) {
	for {
		select {
		case job := <-q.internalQueue:
			q.jobsIn = append(q.jobsIn, job)
			q.workerChanIn = q.readyPool
		case workerChannel := <-q.workerChanIn:
			q.jobChan = workerChannel
			q.workerChanIn = nil
		case <-ctx.Done():
			return
		default:
			if len(q.jobsIn) == 0 || q.jobChan == nil {
				break
			}
			q.jobChan <- q.jobsIn[0]
			q.jobsIn = q.jobsIn[1:]
		}
		q.started = true
	}
}

func (q *jobQueue) submit(ctx context.Context, job job) {
	for {
		if q.started {
			break
		}
	}
	q.wg.Add(1)
	select {
	case q.internalQueue <- job:
	case <-ctx.Done():
		q.wg.Done()
		return
	}
}
