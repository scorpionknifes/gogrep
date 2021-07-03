package main

import (
	"os/exec"
	"testing"
)

func BenchmarkGoGrep(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		Cmd("./gogrep -w '[A-Z]+_SUSPEND'")
	}
}

func BenchmarkRipGrep(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		Cmd("rg -n -w '[A-Z]+_SUSPEND'")
	}
}

func BenchmarkGitGrep(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		Cmd("git grep -P -n -w '[A-Z]+_SUSPEND'")
	}
}

func Cmd(cmd string) []byte {
	command := exec.Command("bash", "-c", cmd)
	command.Dir = "linux-master"
	out, err := command.Output()
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(out))
	return out
}
