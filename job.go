package main

import (
	"context"
	"os"
	"regexp"
)

type job interface {
	Process(ctx context.Context)
}

type grepJob struct {
	path   string
	regex  *regexp.Regexp
	nlines int
	finder finder
}

func (j *grepJob) Process(ctx context.Context) {
	// Replace os.Stdout with ioutil.Discard for benchmarking
	err := j.finder.Find(ctx, os.Stdout, j.path, j.regex, j.nlines)
	if err != nil {
		return
	}
}
