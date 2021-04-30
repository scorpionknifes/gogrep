package main

import (
	"bytes"
	"testing"
)

func Test_finder_Find(t *testing.T) {
	type args struct {
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
		{"no match", args{"empty", "match"}, "", false},
		{"match 1:0", args{"match", "match"}, "1:0: match\n", false},
		{"match 1:4", args{"test match", "match"}, "1:5: test match\n", false},
		{"match 1:4", args{"test match test", "match"}, "1:5: test match test\n", false},
		{"match 2:4", args{"\ntest match test\n", "match"}, "2:5: test match test\n", false},
		{"match 5:6", args{"\n\n\n\ntest tmatch test\n\n\n", "match"}, "5:6: test tmatch test\n", false},
		{"match 1:0, 2:0", args{"match\nmatch", "match"}, "1:0: match\n2:0: match\n", false},
		{"match 1:4, 2:4", args{"test match test\ntest match test", "match"}, "1:5: test match test\n2:5: test match test\n", false},
		{"multiline 1:4, 2:4", args{"test match\ntest match test", "match\n"}, "1:5: test match test\n2:5: test match test\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &finder{}
			w := &bytes.Buffer{}
			if err := f.Find(w, tt.args.text, tt.args.regex); (err != nil) != tt.wantErr {
				t.Errorf("finder.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("finder.Find() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
