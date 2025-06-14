package service

import (
	"context"
	"log"

	"github.com/berezovskyivalerii/testtre/model"
)

type Provider interface {
	FetchUsers(ctx context.Context, url string) ([]model.User, error)
}

type Sender interface {
	SendUser(ctx context.Context, url string, u model.User) error
}

type Dispatcher struct {
	srcURL string
	dstURL string
	p      Provider
	s      Sender
	l      *log.Logger
}

func NewDispatcher(srcURL, dstURL string, p Provider, s Sender, l *log.Logger) *Dispatcher {
	return &Dispatcher{srcURL: srcURL, dstURL: dstURL, p: p, s: s, l: l}
}

func (d *Dispatcher) Run(ctx context.Context) error {
	users, err := d.p.FetchUsers(ctx, d.srcURL)
	if err != nil {
		return err
	}

	for _, u := range users {
		if u.IsBiz() {
			if err := d.s.SendUser(ctx, d.dstURL, u); err != nil {
				d.l.Printf("[ERR] %s -> retry failed: %v", u.Email, err)
			} else {
				d.l.Printf("[OK] sent %s to API B", u.Email)
			}
		} else {
			d.l.Printf("[SKIP] %s", u.Email)
		}
	}
	return nil
}
