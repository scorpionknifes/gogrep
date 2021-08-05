package main

import (
	"bufio"
	"io"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/karrick/godirwalk"
	"github.com/scorpionknifes/go-pcre"
)

func main() {
	f, perr := os.Create("cpu.pprof")
	if perr != nil {
		os.Exit(1)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	run(os.Stdout, runtime.NumCPU())
}

func run(w io.Writer, size int) {
	regexString := "\\b([A-Z]+_SUSPEND)\\b"

	regex := pcre.MustCompileJIT(regexString, 0, 0)

	queues := make(chan string, size)

	bw := bufio.NewWriter(w)
	defer bw.Flush()

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for filePath := range queues {
				file, err := os.Open(filePath)
				if err != nil {
				}

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					runRegex(bw, scanner.Bytes(), &regex)
				}
				file.Close()
			}
		}()
	}

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
}
