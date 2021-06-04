package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"

	"github.com/karrick/godirwalk"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	queue := newJobQueue(runtime.NumCPU() - 1)
	queue.start(ctx)

	match, dirname, exitcode := cli(os.Args)
	if exitcode != nil {
		os.Exit(*exitcode)
	}

	err := godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(filePath string, de *godirwalk.Dirent) error {
			if isIgnore(filePath) {
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
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
	if err != nil {
		panic(err)
	}

	queue.wg.Wait()
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

// isText checks if the file or path is gitignored using git check-ignore
// only works if git is install on PC.
func isIgnore(filePath string) bool {
	cmd := exec.Command("git", "check-ignore", filePath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return false
	}
	if err := cmd.Start(); err != nil {
		return false
	}
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(stdout); err != nil {
		return false
	}
	if buf.String() == "" {
		return false
	}
	return true
}
