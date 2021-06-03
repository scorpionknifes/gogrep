package main

import (
	"context"
)

type mockJob struct {
	done bool
}

func (m *mockJob) Process(ctx context.Context) {
	m.done = true
}
