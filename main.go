package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/berezovskyivalerii/testtre/client"
	"github.com/berezovskyivalerii/testtre/config"
	"github.com/berezovskyivalerii/testtre/service"
)

func main() {
	cfg := config.Load()

	if cfg.DstURL == "" {
		log.Fatal("укажите DST_URL в .env/ENV или флаг --dst")
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cl := client.New()
	d  := service.NewDispatcher(cfg.SrcURL, cfg.DstURL, cl, cl, logger)

	if err := d.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		logger.Fatalf("pipeline error: %v", err)
	}

	logger.Println("graceful shutdown complete")
}
