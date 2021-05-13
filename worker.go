package main

import (
	"context"
	"log"
	"sync"
)

type job interface {
	Process()
}

func newWorker(readyPool chan chan job, done *sync.WaitGroup) *worker {
	return &worker{
		done:             done,
		readyPool:        readyPool,
		assignedJobQueue: make(chan job),
	}
	// TODO create channel and return channel
}

type worker struct {
	done             *sync.WaitGroup
	readyPool        chan chan job
	assignedJobQueue chan job
}

func (w *worker) start(ctx context.Context) {
	w.done.Add(1)
	go func() {
		for {
			select {
			case job := <-w.assignedJobQueue:
				job.Process()
			case <-ctx.Done():
				log.Println("Done worker")
				w.done.Done()
				return
			}
		}
	}()
}
