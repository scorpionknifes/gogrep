package main

import (
	"context"
	"testing"
	"time"
)

func Test_Queue(t *testing.T) {

	t.Run("Stop", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		queue := newJobQueue(1)
		queue.start(ctx)
		queue.submit(ctx, &mockJob{})
		cancel()
		queue.submit(ctx, &mockJob{})
	})

	t.Run("Single task single goroutine", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		queue := newJobQueue(1)
		queue.start(ctx)
		job := &mockJob{}
		queue.submit(ctx, job)
		queue.wg.Wait()
		if !job.done {
			t.Error("Job not complete", job.done)
		}
		cancel()
	})

	t.Run("Two task single goroutine", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		queue := newJobQueue(1)
		queue.start(ctx)
		job1 := &mockJob{}
		job2 := &mockJob{}
		queue.submit(ctx, job1)
		queue.submit(ctx, job2)
		queue.wg.Wait()
		if !job1.done || !job2.done {
			t.Error("Jobs not complete", job1.done, job2.done)
		}
		cancel()
	})

	t.Run("Two task two goroutine", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		queue := newJobQueue(2)
		queue.start(ctx)
		job1 := &mockJob{}
		job2 := &mockJob{}
		queue.submit(ctx, job1)
		queue.submit(ctx, job2)
		queue.wg.Wait()
		if !job1.done || !job2.done {
			t.Error("Jobs not complete", job1.done, job2.done)
		}
		cancel()
	})
}

func Test_QueueInternalQueue(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("Add nothing to internal queue", func(t *testing.T) {
		queue := newJobQueue(0)
		go queue.dispatch(ctx)

		time.Sleep(1 * time.Second)

		if len(queue.jobsIn) != 0 {
			t.Error("queue should be empty")
		}
	})

	t.Run("Add to internal Queue", func(t *testing.T) {
		queue := newJobQueue(0)
		go queue.dispatch(ctx)

		queue.submit(ctx, &mockJob{})

		time.Sleep(1 * time.Second)

		if len(queue.jobsIn) != 1 {
			t.Error("not added to internal queue")
		}
	})

	t.Run("Add multiple to internal Queue", func(t *testing.T) {
		queue := newJobQueue(0)
		go queue.dispatch(ctx)

		queue.submit(ctx, &mockJob{})
		queue.submit(ctx, &mockJob{})
		queue.submit(ctx, &mockJob{})

		time.Sleep(1 * time.Second)

		if len(queue.jobsIn) != 3 {
			t.Error("not added to internal queue")
		}
	})
}
