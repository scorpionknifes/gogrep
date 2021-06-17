package main

import (
	"bytes"
	"embed"
	_ "embed"
	"io"
	"strings"
	"testing"
)

//go:embed data/example.exe
var exeContent []byte

//go:embed data/lorem0.txt
var txtContent []byte

//go:embed data
var data embed.FS

func Test_getFileContentType(t *testing.T) {
	type args struct {
		file io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"load empty file", args{strings.NewReader("")}, "", true},
		{"load text file", args{bytes.NewBuffer(txtContent)}, "text/plain; charset=utf-8", false},
		{"load exe file", args{bytes.NewBuffer(exeContent)}, "application/octet-stream", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := getFileContentType(tt.args.file)
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
