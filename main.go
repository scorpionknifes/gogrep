package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/fatih/color"
)

func main() {
	file, err := os.Open("data/lorem.txt")
	if err != nil {
		panic(err)
	}

	match := "lorem"

	r, err := regexp.Compile(match)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if r.MatchString(scanner.Text()) {
			fmt.Println(r.ReplaceAllString(scanner.Text(), color.YellowString(match)))
		}
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
