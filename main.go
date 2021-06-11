package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/karrick/godirwalk"
)

func main() {
	start(os.Args)
}

func start(args []string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	queue := newJobQueue(runtime.NumCPU() - 1)
	queue.start(ctx)

	match, dirname, exitcode := cli(args)
	if exitcode != nil {
		os.Exit(*exitcode)
	}

	err := godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: godirCallBack(ctx, queue, match, newMatcher(dirname)),
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
	if err != nil {
		panic(err)
	}

	queue.wg.Wait()
}

func godirCallBack(ctx context.Context, queue *jobQueue, match string, matcher gitignore.Matcher) func(filePath string, _ *godirwalk.Dirent) error {
	return func(filePath string, _ *godirwalk.Dirent) error {
		if isIgnore(filePath, matcher) {
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

// cli
func cli(args []string) (string, string, *int) {
	const relativePath = "."
	name := args[0]

	switch len(args) {
	case 2:
		if args[1] == "--help" {
			fmt.Printf("Usage: %s PATTERN [PATH]\n", name)
			fmt.Printf("Search for PATTERN in each FILE in PATH")
			fmt.Printf("Example: %s 'hello world' ./folder", name)
			exitcode := 0
			return "", "", &exitcode
		}
		return args[1], relativePath, nil
	case 3:
		return args[1], args[2], nil
	default:
		fmt.Printf("usage: %s PATTERN [PATH]\n", name)
		fmt.Printf("Try '%s --help' for more information\n", name)
		exitcode := 1
		return "", "", &exitcode
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

func newMatcher(dirname string) gitignore.Matcher {
	fs := osfs.New(".")
	ps, err := gitignore.ReadPatterns(fs, strings.Split(dirname, "/"))
	if err != nil {
		panic(err)
	}
	return gitignore.NewMatcher(ps)
}

// isIgnore checks if the file or path is gitignored using git check-ignore
// only works if git is install on PC.
func isIgnore(filePath string, matcher gitignore.Matcher) bool {
	return matcher.Match(strings.Split(filePath, "/"), true)
}
