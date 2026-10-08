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

	a := app.New()

	go func() {
		if err := a.Start(ctx); err != nil {
			log.Fatalf("failed to start app: %v", err)
		}
	}()

	log.Println("Sayanox Chain started")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	<-sigCh

	log.Println("Shutting down Sayanox Chain...")

	if err := a.Stop(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	log.Println("Node stopped")
}
