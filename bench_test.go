package main

import (
	"os/exec"
	"testing"
)

func BenchmarkGoGrep(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		Cmd("./gogrep.exe")
	}
}

func BenchmarkRipGrep(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		Cmd("rg", "-n", "-w", "[A-Z]+_SUSPEND")
	}
}

func Cmd(name string, args ...string) []byte {
	command := exec.Command(name, args...)
	command.Dir = "linux-master"
	out, err := command.Output()
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(out))
	return out
}
