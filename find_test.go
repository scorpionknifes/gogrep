package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"testing"
)

func Benchmark_Random(b *testing.B) {
	data, err := ioutil.ReadFile("data/lorem.txt")
	if err != nil {
		panic(err)
	}
	match := "lorem" //os.Args[1]

	// file, err := os.Open("file.go")
	// if err != nil {
	// 	panic(err)
	// }
	// getFileContentType(file)

	f := finder{}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err = f.Find(context.Background(), io.Discard, "", string(data), match)
		if err != nil {
			b.Fail()
		}
	}
}

func Test_FinderFind(t *testing.T) {
	type args struct {
		path  string
		text  string
		regex string
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{"empty", args{}, "", false},
		{"no match", args{"", "empty", "match"}, "", false},
		{"bad regex", args{"", "empty", "[match"}, "", true},
		{"match 1:0", args{"", "match", "match"}, "1:0: match\n", false},
		{"path match 1:0", args{"test.go", "match", "match"}, "test.go:1:0: match\n", false},
		{"match 1:0 long", args{"", "aaaaaaaaaaaaaaaaaaaaaaa match aaaaaaaaaaaaaaaaaaaaaa", "match"}, "1:24: aaaaaaaaaaaaaaaaaaa match aaaaaaaaaaaaaaaaaaa\n", false},
		{"match 1:4", args{"", "test match", "match"}, "1:5: test match\n", false},
		{"match 1:4", args{"", "test match test", "match"}, "1:5: test match test\n", false},
		{"match 2:4", args{"", "\ntest match test\n", "match"}, "2:5: test match test\n", false},
		{"match 5:6", args{"", "\n\n\n\ntest tmatch test\n\n\n", "match"}, "5:6: test tmatch test\n", false},
		{"match 1:0, 2:0", args{"", "match\nmatch", "match"}, "1:0: match\n2:0: match\n", false},
		{"match 1:4, 2:4", args{"", "test match test\ntest match test", "match"}, "1:5: test match test\n2:5: test match test\n", false},
		{"multiline 1:4, 2:4", args{"", "test match\ntest match test", "match\n"}, "1:5: test match\\ntest match test\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &finder{}
			w := &bytes.Buffer{}
			// ctx, cancel := context.WithCancel(context.Background())
			// cancel()
			if err := f.Find(context.Background(), w, tt.args.path, tt.args.text, tt.args.regex); (err != nil) != tt.wantErr {
				t.Errorf("finder.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("finder.Find() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_FinderFindContext(t *testing.T) {
	type args struct {
		path  string
		text  string
		regex string
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{"empty", args{}, "", false},
		{"match 1:0", args{"", "match", "match"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &finder{}
			w := &bytes.Buffer{}
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			if err := f.Find(ctx, w, tt.args.path, tt.args.text, tt.args.regex); (err != nil) != tt.wantErr {
				t.Errorf("finder.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("finder.Find() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
