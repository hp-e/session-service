package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"studio-service/core/models"
	"studio-service/ports/publishers"
	"studio-service/ports/sessions"
	"sync"
	"syscall"
)

func handlePanic() {
	if err := recover(); err != nil {
		// if logger != nil {
		// 	_ = logger.LogError("context", "main", "message", "Caught panic", "error", err)
		// 	logger.Error(string(debug.Stack()))
		// }
	}
}

func main() {
	var (
	// err  error
	// stop context.CancelFunc
	)

	defer handlePanic()
	// mainCtx, stop = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	// defer stop()

	ssePublisher := publishers.NewSsePublisher()
	publisher := publishers.NewScreenPublisher()

	sessionManager := sessions.NewSessionManager(publisher)

	go func() {
		log.Println("Starting SSE server on :8080")
		log.Fatal(http.ListenAndServe(":8080", ssePublisher))
	}()

	// Create sample sessions
	wg := &sync.WaitGroup{}

	// Create a sample session and run it
	sessionPlan := models.CreateSampleSessions()[0]
	session := &models.Session{
		RoomId:   "room-1",
		Id:       "session-1",
		Plan:     sessionPlan,
		Timeline: nil,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go sessionManager.Run(ctx, session)

	// Handle termination signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	cancel()
	log.Println("Shutting down...")

	wg.Wait()
	log.Println("All sessions completed")
}
