package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sayan9168/sayanox-chain/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a, err := app.New()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	go func() {
		if err := a.Start(ctx); err != nil {
			log.Printf("failed to start app: %v", err)
			cancel()
		}
	}()

	log.Println("Sayanox Chain started")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	select {
	case <-sigCh:
		log.Println("Shutting down Sayanox Chain...")
	case <-ctx.Done():
	}

	if err := a.Stop(context.Background()); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	log.Println("Node stopped")
}
