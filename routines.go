package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type routine struct {
	path string
}

func consumer(r chan routine, regex string) {

	rexp, err := regexp.Compile(regex)
	if err != nil {
		return
	}

	for f := range r {
		find(f, rexp)
	}
}

func find(f routine, regex *regexp.Regexp) {
	lineNumber := 0
	file, err := os.Open(f.path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buffer := make([]byte, 512)

	_, err = file.Read(buffer)
	if err != nil {
		return
	}

	contentType := http.DetectContentType(buffer)

	types := strings.Split(contentType, "/")
	if types[0] != "text" {
		return
	}

	// NewReaderSize
	sc := bufio.NewScanner(file)
	for sc.Scan() {

		lineNumber++
		line := sc.Bytes()
		if lineNumber == 1 {
			line = append(buffer[:], line[:]...)
			fmt.Printf("%s\n", line)
		}

		if regex.Match(line) {
			allIndex := regex.FindAllIndex(line, -1)

			for _, index := range allIndex {
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

				// f.print(lineNumber, charNumber, head, match, tail)
				fmt.Printf("%d, %d, %s%s%s\n", lineNumber, charNumber, head, match, tail)
			}
		}
	}
}
