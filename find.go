package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/fatih/color"
)

const (
	defaultHead = 20
	defaultTail = 20
)

type finder interface {
	Find(ctx context.Context, w io.Writer, path string, regex *regexp.Regexp, nlines int) error
}

type lineFinder struct {
	path  string
	w     io.Writer
	text  *string
	regex *regexp.Regexp
}

func (f *lineFinder) Find(ctx context.Context, w io.Writer, path string, regex *regexp.Regexp, _ int) error {

	f.path = path
	f.w = w
	f.regex = regex

	return f.find(ctx)
}

func (f *lineFinder) find(ctx context.Context) error {
	lineNumber := 0

	file, err := os.Open(f.path)
	if err != nil {
		return err
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		lineNumber++
		line := sc.Bytes()
		if f.regex.Match(line) {
			allIndex := f.regex.FindAllIndex(line, -1)

			for _, index := range allIndex {
				if err := ctx.Err(); err != nil {
					return err
				}
				charNumber := index[0]

				// Slice of bytes
				match := line[index[0]:index[1]]

				headNumber := 0
				if index[0] > 20 {
					headNumber = index[0] - 20
				}
				tailNumber := len(line)
				if tailNumber > index[1]+20 {
					tailNumber = index[1] + 20
				}

				head := line[headNumber:index[0]]

				tail := line[index[1]:tailNumber]

				f.print(lineNumber, charNumber, head, match, tail)
			}

		}

	}
	return nil
}

func (f *lineFinder) print(lineNumber, charNumber int, head, match, tail []byte) {
	path := ""
	if f.path != "" {
		path = color.RedString("%s", f.path) + ":"
	}

	fmt.Fprintf(
		f.w,
		"%s%s:%s: %s%s%s\n",
		path,
		color.GreenString("%d", lineNumber),
		color.GreenString("%d", charNumber),
		head,
		color.YellowString("%s", match),
		tail,
	)
}
