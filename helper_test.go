package main

import (
	"os"
	"testing"
)

func Test_getFileContentType(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"load text file", args{"./data/lorem0.txt"}, "text/plain; charset=utf-8", false},
		{"load exe file", args{"./data/example.exe"}, "application/octet-stream", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileContentType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got, err := getFileContentType(file)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileContentType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getFileContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}
