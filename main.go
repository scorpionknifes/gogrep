package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/karrick/godirwalk"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	name := os.Args[0]
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s PATTERN [PATH]\n", name)
		fmt.Fprintf(flag.CommandLine.Output(), "Search for PATTERN in each FILE in PATH\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Example: %s 'hello world' ./folder\n", name)
		flag.PrintDefaults()
	}
	flag.Parse()

	if err := start(ctx, flag.Args()); err != nil {
		fmt.Printf("usage: %s PATTERN [PATH]\n", name)
		fmt.Printf("Try '%s --help' for more information\n", name)
		os.Exit(1)
	}
}

func start(ctx context.Context, args []string) error {
	log.Println(args)
	queue := newJobQueue(runtime.NumCPU() - 1)
	queue.start(ctx)

	if len(args) == 0 {
		return errors.New("Invalid Arguments")
	}

	dirname := "."
	if len(args) > 1 {
		dirname = args[1]
	}

	matcher, _ := newMatcher(dirname)

	err := godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: godirCallBack(ctx, queue, args[0], matcher),
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
	if err != nil {
		return err
	}

	queue.wg.Wait()
	return nil
}

func godirCallBack(ctx context.Context, queue *jobQueue, match string, matcher gitignore.Matcher) func(filePath string, _ *godirwalk.Dirent) error {
	return func(filePath string, _ *godirwalk.Dirent) error {
		if matcher != nil && isIgnore(filePath, matcher) {
			return godirwalk.SkipThis
		}

		// don't check .git folder but make sure to scan .gitignore
		if strings.Contains(filePath, ".git") && !strings.Contains(filePath, ".gitignore") {
			return godirwalk.SkipThis
		}

		if !isText(filePath) {
			return nil
		}

		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil
		}

		queue.submit(ctx, &grepJob{filePath, string(data), match})
		return nil
	}

}

// isText checks if the file content is text
func isText(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	contentType, err := getFileContentType(file)
	if err != nil {
		return false
	}
	types := strings.Split(contentType, "/")
	return types[0] == "text"
}

func newMatcher(dirname string) (gitignore.Matcher, error) {
	fs := osfs.New(".")
	ps, err := gitignore.ReadPatterns(fs, strings.Split(dirname, "/"))
	if err != nil {
		return nil, err
	}
	return gitignore.NewMatcher(ps), nil
}

// isIgnore checks if the file or path is gitignored using git check-ignore
// only works if git is install on PC.
func isIgnore(filePath string, matcher gitignore.Matcher) bool {
	return matcher.Match(strings.Split(filePath, "/"), true)
}
