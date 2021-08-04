package main

import (
	"bufio"
	"os"
	"runtime"

	"github.com/karrick/godirwalk"
	"github.com/scorpionknifes/go-pcre"
)

func main() {
	run(2, 10)
}

func run(readNum, lineNum int) {
	w := bufio.NewWriter(os.Stdout)

	regexString := "\\b([A-Z]+_SUSPEND)\\b"

	regex := pcre.MustCompileJIT(regexString, 0, 0)

	queues := make(chan string, runtime.NumCPU())

	go func() {
		for queue := range queues {
			r := serialConsumer(&regex)
			r <- queue
			close(r)
		}
	}()

	err := godirwalk.Walk(".", &godirwalk.Options{
		Callback: func(filePath string, _ *godirwalk.Dirent) error {
			queues <- filePath
			return nil
		},
		Unsorted: true,
	})

	close(queues)

	if err != nil {
		panic(err)
	}

	err = w.Flush()
	if err != nil {
		panic(err)
	}

}
