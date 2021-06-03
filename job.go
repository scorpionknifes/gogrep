package main

import (
	"context"
	"os"
)

type job interface {
	Process(ctx context.Context)
}

type grepJob struct {
	path  string
	data  string
	match string
}

func (j *grepJob) Process(ctx context.Context) {
	f := finder{}
	// Replace os.Stdout with ioutil.Discard for benchmarking
	err := f.Find(os.Stdout, ctx, j.path, j.data, j.match)
	if err != nil {
		return
	}
}
