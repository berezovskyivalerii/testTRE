package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/berezovskyivalerii/testtre/model"
)

func serverWithStatuses(statuses ...int) (*httptest.Server, *int) {
	var i int
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statuses[i])
		i++
	}))
	return s, &i
}

func TestSendUser_RetrySucceeds(t *testing.T) {
	srv, cnt := serverWithStatuses(500, 500, 201)
	defer srv.Close()

	cl := New()
	user := model.User{Name: "Polina", Email: "polina@foo.biz"}
	err := cl.SendUser(context.Background(), srv.URL, user)
	if err != nil {
		t.Fatalf("SendUser returned error: %v", err)
	}
	if *cnt != 3 {
		t.Fatalf("expected 3 attempts, got %d", *cnt)
	}
}

func TestSendUser_StopsOn400(t *testing.T) {
	srv, cnt := serverWithStatuses(400)
	defer srv.Close()

	cl := New()
	user := model.User{Name: "Alexandra", Email: "alexandra_m_02@foo.biz"}
	err := cl.SendUser(context.Background(), srv.URL, user)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if *cnt != 1 {
		t.Fatalf("expected 1 attempt, got %d", *cnt)
	}
}

func TestSendUser_ContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	cl := New()
	user := model.User{Name: "Valera", Email: "valera@foo.biz"}
	err := cl.SendUser(ctx, srv.URL, user)
	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}
}
