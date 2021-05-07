package main

import "sync"

type jobQueue struct {
	internalQueue     chan job
	readyPool         chan chan job
	workers           []*worker
	dispatcherStopped *sync.WaitGroup
	workersStopped    *sync.WaitGroup
	quit              chan bool
}

func newJobQueue(maxWorkers int) *jobQueue {
	workersStopped := sync.WaitGroup{}
	readyPool := make(chan chan job, maxWorkers)
	workers := make([]*worker, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		workers[i] = newWorker(readyPool, &workersStopped)
	}
	return &jobQueue{
		internalQueue:     make(chan job),
		readyPool:         readyPool,
		workers:           workers,
		dispatcherStopped: &sync.WaitGroup{},
		workersStopped:    &workersStopped,
		quit:              make(chan bool),
	}
}

func (q *jobQueue) start() {
	for i := 0; i < len(q.workers); i++ {
		q.workers[i].start()
	}
	go q.dispatch()
}

func (q *jobQueue) dispatch() {
	q.dispatcherStopped.Add(1)
	for {
		select {
		case job := <-q.internalQueue: // We got something in on our queue
			workerChannel := <-q.readyPool // Check out an available worker
			workerChannel <- job           // Send the request to the channel
		case <-q.quit:
			for i := 0; i < len(q.workers); i++ {
				q.workers[i].stop()
			}
			q.workersStopped.Wait()
			q.dispatcherStopped.Done()
			return
		}
	}
}

func (q *jobQueue) submit(job job) {
	q.internalQueue <- job
}

func (q *jobQueue) stop() {
	q.quit <- true
	q.dispatcherStopped.Wait()
}
