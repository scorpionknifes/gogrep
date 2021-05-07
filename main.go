package main

import (
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	data, err := ioutil.ReadFile("data/lorem.txt")
	if err != nil {
		panic(err)
	}
	match := os.Args[1]

	// file, err := os.Open("file.go")
	// if err != nil {
	// 	panic(err)
	// }
	// GetFileContentType(file)

	f := finder{}

	for {
		err = f.Find(os.Stdout, string(data), match)
		if err != nil {
			panic(err)
		}
	}
}
