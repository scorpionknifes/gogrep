package main

import (
	"io/ioutil"
	"os"
)

func main() {

	data, err := ioutil.ReadFile("data/lorem.txt")
	if err != nil {
		panic(err)
	}
	match := os.Args[1]

	f := finder{}

	err = f.Find(os.Stdout, string(data), match)
	if err != nil {
		panic(err)
	}
}
