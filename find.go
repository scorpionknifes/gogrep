package main

import (
	"fmt"
	"io"
	"log"
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
	w     io.Writer
	text  *string
	regex *regexp.Regexp
}

func (f *finder) Find(w io.Writer, text string, regex string) error {
	r, err := regexp.Compile(regex)
	if err != nil {
		return err
	}
	f.w = w
	f.regex = r
	f.text = &text

	return f.find()
}

func (f *finder) find() error {
	if *f.text == "" {
		return nil
	}
	text := *f.text + "\n"
	newlines := newlineRegex.FindAllStringIndex(text, -1)
	newline := 0

	allStr := f.regex.FindAllStringIndex(text, -1)

	for _, str := range allStr {
		for len(newlines) != newline && str[0] > newlines[newline][0] {
			newline++
		}
		match := color.YellowString(text[str[0]:str[1]])
		lineNumber := newline + 1
		charNumber := str[0]
		headNumber := 0
		head := ""
		if newline > 0 {
			headNumber = defaultHead
			charNumber = str[0] - newlines[newline-1][1]
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
		if newline < len(newlines) {
			lastNumber := -str[1] + newlines[newline][0]
			tailNumber = defaultHead + defaultTail - headNumber
			log.Println(tailNumber)
			if lastNumber < tailNumber && lastNumber > 0 {
				tailNumber = lastNumber
			}
		}
		tail := ""
		if tailNumber+str[1] > len(allStr) {
			tail = text[str[1] : str[1]+len(allStr)-1]
		} else {
			tail = text[str[1] : str[1]+tailNumber]
		}

		f.print(lineNumber, charNumber, head, match, tail)
	}
	return nil
}

func (f *finder) print(lineNumber, charNumber int, head, match, tail string) {
	fmt.Fprintf(
		f.w,
		"%s:%s: %s%s%s\n",
		color.GreenString("%d", lineNumber),
		color.GreenString("%d", charNumber),
		newlineRegex.ReplaceAllString(head, "\\n"),
		newlineRegex.ReplaceAllString(match, "\\n"),
		newlineRegex.ReplaceAllString(tail, "\\n"),
	)
}
