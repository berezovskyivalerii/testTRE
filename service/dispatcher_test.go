package service

import (
	"context"
	"log"
	"testing"
	"io"

	"github.com/berezovskyivalerii/testtre/model"
)

type stubProv struct{}

func (stubProv) FetchUsers(ctx context.Context, _ string) ([]model.User, error) {
	return []model.User{
		{Name: "A", Email: "a@corp.biz"},
		{Name: "B", Email: "b@corp.com"},
	}, nil
}

type countingSender struct{ n int }

func (s *countingSender) SendUser(ctx context.Context, _ string, _ model.User) error {
	s.n++
	return nil
}

func TestDispatcher_Run(t *testing.T) {
	cs := &countingSender{}

	silent := log.New(io.Discard, "", 0)

	d := NewDispatcher("src", "dst", stubProv{}, cs, silent)

	if err := d.Run(context.Background()); err != nil {
		t.Fatalf("Dispatcher.Run error: %v", err)
	}
	if cs.n != 1 {
		t.Fatalf("expected 1 SendUser call, got %d", cs.n)
	}
}
