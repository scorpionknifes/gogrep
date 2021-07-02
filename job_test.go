package main

import (
	"context"
	"regexp"
	"testing"
)

func Test_grepJob_Process(t *testing.T) {
	tests := []struct {
		name string
		j    *grepJob
	}{
		{"empty", &grepJob{}},
		{"basic", &grepJob{"", compileRegex("empty"), 1, &lineFinder{}}},
		{"error", &grepJob{"", compileRegex("empty"), 1, &lineFinder{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.j.Process(context.Background())
		})
	}
}

func compileRegex(regex string) *regexp.Regexp {
	r, err := regexp.Compile(regex)
	if err != nil {
		return nil
	}
	return r
}
