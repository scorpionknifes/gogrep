package main

import (
	"context"
	"sync"
)

type jobQueue struct {
	internalQueue chan job
	jobChans      []chan job
	jobChan       chan job
	workers       []*worker
	jobsIn        []job
	jobIn         job
	readyPool     chan chan job
	wg            *sync.WaitGroup
	started       bool
}

func newJobQueue(maxWorkers int) *jobQueue {
	workers := make([]*worker, maxWorkers)
	readyPool := make(chan chan job, maxWorkers)
	wg := sync.WaitGroup{}

	for i := 0; i < maxWorkers; i++ {
		workers[i] = newWorker(readyPool, &wg, i)
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
	q.started = true
	for {
		select {
		case job := <-q.internalQueue:
			q.jobsIn = append(q.jobsIn, job)
			if len(q.jobsIn) != 0 && len(q.jobChans) != 0 && q.jobChan == nil {
				q.jobChan = q.jobChans[0]
				q.jobChans = q.jobChans[1:]
				q.jobIn = q.jobsIn[0]
			}
		case workerChannel := <-q.readyPool:
			q.jobChans = append(q.jobChans, workerChannel)
			if len(q.jobsIn) != 0 && len(q.jobChans) != 0 && q.jobChan == nil {
				q.jobChan = q.jobChans[0]
				q.jobChans = q.jobChans[1:]
				q.jobIn = q.jobsIn[0]
			}
		case <-ctx.Done():
			return
		case q.jobChan <- q.jobIn:
			q.jobsIn = q.jobsIn[1:]
			if len(q.jobsIn) != 0 && len(q.jobChans) != 0 {
				q.jobChan = q.jobChans[0]
				q.jobChans = q.jobChans[1:]
				q.jobIn = q.jobsIn[0]
			} else {
				q.jobChan = nil
				q.jobIn = nil
			}
		}
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
