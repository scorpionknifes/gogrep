package main

import (
	"context"
	"testing"
)

func Test_start(t *testing.T) {
	type args struct {
		ctx  context.Context
		args []string
	}
	tests := []struct {
		name string
		args args
	}{
		{"empty", args{context.TODO(), []string{"gogrep"}}},
		{"lorem", args{context.TODO(), []string{"gogrep", "lorem."}}},
		{"lorem data", args{context.TODO(), []string{"gogrep", "lorem.", "data"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start(tt.args.ctx, tt.args.args)
		})
	}
}

func Test_isTest(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{}, false},
		{"load empty file", args{"./data/empty.txt"}, false},
		{"load text file", args{"./data/lorem0.txt"}, true},
		{"load exe file", args{"./data/example.exe"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isText(tt.args.filePath); got != tt.want {
				t.Errorf("isIgnore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isIgnore(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{}, false},
		{"load empty file", args{"./data/empty.txt"}, false},
		{"load text file", args{"./data/lorem0.txt"}, false},
		{"load exe file", args{"./data/example.exe"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher, _ := newMatcher("data")
			if got := isIgnore(tt.args.filePath, matcher); got != tt.want {
				t.Errorf("isIgnore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}
