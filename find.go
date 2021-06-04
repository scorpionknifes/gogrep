package main

import (
	"context"
	"fmt"
	"io"
	"regexp"

	"github.com/fatih/color"
)

const (
	defaultHead = 20
	defaultTail = 20
)

var (
	newlineRegex = regexp.MustCompile("\n")
)

type finder struct {
	path   string
	w      io.Writer
	text   *string
	regex  *regexp.Regexp
	nlines int
}

func (f *finder) Find(ctx context.Context, w io.Writer, path string, text string, regex string) error {
	r, err := regexp.Compile(regex)
	if err != nil {
		return err
	}
	f.path = path
	f.w = w
	f.regex = r
	f.text = &text
	f.nlines = len(newlineRegex.FindAllStringIndex(regex, -1))

	return f.find(ctx)
}

func (f *finder) find(ctx context.Context) error {
	if *f.text == "" {
		return nil
	}
	text := *f.text + "\n"
	newlines := newlineRegex.FindAllStringIndex(text, -1)
	newline := 0

	allStr := f.regex.FindAllStringIndex(text, -1)

	for _, str := range allStr {
		if err := ctx.Err(); err != nil {
			return err
		}
		for len(newlines) != newline && str[0] > newlines[newline][0] {
			newline++
		}
		match := color.YellowString(text[str[0]:str[1]])
		lineNumber := newline + 1

		charNumber := str[0]

		headNumber := 0
		head := ""
		if newline > 0 {
			//steps
			charNumber = str[0] - newlines[newline-1][1]

			//steps
			headNumber = defaultHead
			if charNumber < headNumber {
				headNumber = charNumber
			}
			head = text[str[0]-headNumber : str[0]]
		} else {
			if str[0] > defaultHead {
				head = text[str[0]-defaultHead : str[0]]
				headNumber = defaultHead
			} else {
				head = text[0:str[0]]
				headNumber = str[0]
			}
		}
		tailNumber := 0
		if newline+f.nlines < len(newlines) {
			lastNumber := -str[1] + newlines[newline+f.nlines][0]
			tailNumber = defaultHead + defaultTail - headNumber
			if lastNumber < tailNumber {
				tailNumber = lastNumber
			}
		}
		tail := text[str[1] : str[1]+tailNumber]
		f.print(lineNumber, charNumber, head, match, tail)
	}
	return nil
}

func (f *finder) print(lineNumber, charNumber int, head, match, tail string) {
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
		newlineRegex.ReplaceAllString(head, "\\n"),
		newlineRegex.ReplaceAllString(match, "\\n"),
		newlineRegex.ReplaceAllString(tail, "\\n"),
	)
}
