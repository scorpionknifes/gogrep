package main

import (
	"context"
	"testing"
)

func Test_grepJob_Process(t *testing.T) {
	tests := []struct {
		name string
		j    *grepJob
	}{
		{"empty", &grepJob{}},
		{"basic", &grepJob{"", "empty", "match"}},
		{"error", &grepJob{"", "empty", "["}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.j.Process(context.Background())
		})
	}
}
