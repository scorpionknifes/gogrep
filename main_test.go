package main

import (
	"reflect"
	"testing"
)

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
			if got := isIgnore(tt.args.filePath); got != tt.want {
				t.Errorf("isIgnore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cli(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
		want2 *int
	}{
		{"empty", args{[]string{"gogrep"}}, "", "", intPtr(1)},
		{"help", args{[]string{"gogrep", "--help"}}, "", "", intPtr(0)},
		{"1 args", args{[]string{"gogrep", "test"}}, "test", ".", nil},
		{"2 args", args{[]string{"gogrep", "test", "./test"}}, "test", "./test", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := cli(tt.args.args)
			if got != tt.want {
				t.Errorf("cli() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("cli() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("cli() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}
