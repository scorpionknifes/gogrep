package main

import (
	"context"
	"io/ioutil"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
)

type grepJob struct {
	path  string
	data  string
	match string
}

func (j *grepJob) Process(ctx context.Context) {
	f := finder{}
	err := f.Find(os.Stdout, ctx, j.path, j.data, j.match)
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	queue := newJobQueue(runtime.NumCPU() - 1)
	queue.start(ctx)

	match := os.Args[1]
	err := filepath.Walk(os.Args[2],
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				panic(err)
			}
			contentType, err := GetFileContentType(file)
			if err != nil {
				panic(err)
			}
			types := strings.Split(contentType, "/")
			if types[0] != "text" {
				return nil
			}

			data, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}

			queue.submit(ctx, &grepJob{path, string(data), match})
			return nil
		})
	if err != nil {
		panic(err)
	}

	queue.wg.Wait()
}
