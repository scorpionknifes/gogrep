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

	match, dirname := cli()

	err := godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: func(filePath string, de *godirwalk.Dirent) error {
			if isIgnore(filePath) {
				return godirwalk.SkipThis
			}

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
func cli() (string, string) {
	const relativePath = "."
	name := os.Args[0]

	switch len(os.Args) {
	case 2:
		if os.Args[1] == "--help" {
			fmt.Printf("Usage: %s PATTERN [PATH]\n", name)
			fmt.Printf("Search for PATTERN in each FILE in PATH")
			fmt.Printf("Example: %s 'hello world' ./folder", name)
			os.Exit(0)
		}
		return os.Args[1], relativePath
	case 3:
		return os.Args[1], os.Args[2]
	default:
		fmt.Printf("usage: %s PATTERN [PATH]\n", name)
		fmt.Printf("Try '%s --help' for more information\n", name)
		os.Exit(1)
		return "", ""
	}
}

// isText checks if the file content is text
func isText(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	contentType, err := GetFileContentType(file)
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
