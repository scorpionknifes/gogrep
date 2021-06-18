package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"strings"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/karrick/godirwalk"
)

func main() {

	// f, perr := os.Create("cpu.pprof")
	// if perr != nil {
	// 	os.Exit(1)
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	name := os.Args[0]
	word := flag.Bool("w", false, "word boundary")
	multiline := flag.Bool("U", false, "multiline")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s PATTERN [PATH]\n", name)
		fmt.Fprintf(flag.CommandLine.Output(), "Search for PATTERN in each FILE in PATH\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Example: %s 'hello world' ./folder\n", name)
		flag.PrintDefaults()
	}
	flag.Parse()

	if err := start(ctx, flag.Args(), *word, *multiline); err != nil {
		fmt.Println(err)
		fmt.Printf("usage: %s PATTERN [PATH]\n", name)
		fmt.Printf("Try '%s --help' for more information\n", name)
		os.Exit(1)
	}
}

func start(ctx context.Context, args []string, word bool, multiline bool) error {
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

	regex := args[0]
	if word {
		regex = "\\b(" + regex + ")\\b"
	}

	r, err := regexp.Compile(regex)
	if err != nil {
		return err
	}

	var finder finder
	inlines := 0
	finder = &lineFinder{}
	if multiline {
		inlines = len(newlineRegex.FindAllStringIndex(regex, -1))
		finder = &multilineFinder{}
	}

	matcher, _ := newMatcher(dirname)

	err = godirwalk.Walk(dirname, &godirwalk.Options{
		Callback: godirCallBack(ctx, queue, r, inlines, finder, matcher),
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
	if err != nil {
		return err
	}

	queue.wg.Wait()
	return nil
}

func godirCallBack(ctx context.Context, queue *jobQueue, match *regexp.Regexp, inlines int, finder finder, matcher gitignore.Matcher) func(filePath string, _ *godirwalk.Dirent) error {
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

		// data, err := ioutil.ReadFile(filePath)
		// if err != nil {
		// 	return nil
		// }

		queue.submit(ctx, &grepJob{filePath, match, inlines, finder})
		return nil
	}

}

// isText checks if the file content is text
func isText(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

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
