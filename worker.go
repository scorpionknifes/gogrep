package main

import "sync"

type job interface {
	Process()
}

func newWorker(readyPool chan chan job, done *sync.WaitGroup) *worker {
	return &worker{
		done:             done,
		readyPool:        readyPool,
		assignedJobQueue: make(chan job),
		quit:             make(chan bool),
	}
}

type worker struct {
	done             *sync.WaitGroup
	readyPool        chan chan job
	assignedJobQueue chan job

	quit chan bool
}

func (w *worker) start() {
	w.done.Add(1)
	go func() {
		for {
			w.readyPool <- w.assignedJobQueue
			select {
			case job := <-w.assignedJobQueue:
				job.Process()
			case <-w.quit:
				w.done.Done()
				return
			}
		}
	}()
}

func (w *worker) stop() {
	w.quit <- true
}
