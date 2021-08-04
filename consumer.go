package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/scorpionknifes/go-pcre"
)

type routine struct {
	value string
}

func serialConsumer(regex *pcre.Regexp) chan string {
	f := make(chan string)

	go func() {
		w := bufio.NewWriter(os.Stdout)
		defer w.Flush()

		for filePath := range f {
			file, err := os.Open(filePath)
			if err != nil {
			}

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				runRegex(w, scanner.Bytes(), regex)
			}
			file.Close()
		}
	}()

	return f
}

func bufferConsumer(r chan routine, l chan []byte, wg *sync.WaitGroup, regex *pcre.Regexp) {

	w := bufio.NewWriter(io.Discard)
	defer w.Flush()
	defer wg.Done()

	for f := range r {
		wg.Add(1)
		file, err := os.Open(f.value)
		if err != nil {
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			l <- scanner.Bytes()
		}
		wg.Done()

		w.Flush()
		file.Close()
	}
}

func consumer(r chan routine, wg *sync.WaitGroup, regex *pcre.Regexp) {

	w := bufio.NewWriter(io.Discard)
	defer w.Flush()

	queueRegex := NewListMutex()

	l := make(chan []byte)

	var lwg sync.WaitGroup

	go func() {
		for {
			queueRegex.Lock()
			if queueRegex.Len() > 0 {
				e := queueRegex.Front()
				queueRegex.Remove(e)
				l <- e.Value.([]byte)
			}
			queueRegex.Unlock()
		}
	}()

	for i := 0; i < 3; i++ {
		go lineConsumer(w, l, &lwg, regex)
	}

	for f := range r {
		file, err := os.Open(f.value)
		if err != nil {
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lwg.Add(1)
			queueRegex.Lock()
			queueRegex.PushBack(scanner.Bytes())
			queueRegex.Unlock()
		}
		wg.Done()
		w.Flush()
	}
}

func lineConsumer(w io.Writer, l chan []byte, lwg *sync.WaitGroup, regex *pcre.Regexp) {
	for line := range l {
		lwg.Add(1)
		runRegex(w, line, regex)
		lwg.Done()
	}
}

func runRegex(w io.Writer, line []byte, regex *pcre.Regexp) {

	matcher := regex.NewMatcher()
	if matcher.Match(line, 0) {
		fmt.Fprintf(w, "%s\n", line)
	}
	// if len(regex.FindIndex(line, 0)) != 0 {
	// 	fmt.Fprintf(w, "%s\n", line)
	// }
}
