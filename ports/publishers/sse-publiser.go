package publishers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"studio-service/core/models"
)

type SsePublisher struct {
	clients map[string]chan models.Event
}

func NewSsePublisher() *SsePublisher {
	return &SsePublisher{
		clients: make(map[string]chan models.Event),
	}
}

func (sp *SsePublisher) RegisterClient(sessionId string) chan models.Event {
	sp.clients[sessionId] = make(chan models.Event, 10)
	return sp.clients[sessionId]
}

func (sp *SsePublisher) RemoveClient(sessionId string) {
	if ch, exists := sp.clients[sessionId]; exists {
		close(ch)
		delete(sp.clients, sessionId)
	}
}

func (sp *SsePublisher) Publish(ctx context.Context, event models.Event) error {
	if ch, exists := sp.clients[event.SessionId]; exists {
		select {
		case ch <- event:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}

func (sp *SsePublisher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("sessionId")
	if sessionId == "" {
		http.Error(w, "Missing sessionId", http.StatusBadRequest)
		return
	}

	clientChan := sp.RegisterClient(sessionId)
	defer sp.RemoveClient(sessionId)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case event := <-clientChan:
			data, err := json.Marshal(event)
			if err != nil {
				log.Printf("Failed to marshal event: %v", err)
				continue
			}
			_, err = w.Write([]byte("data: " + string(data) + "\n\n"))
			if err != nil {
				log.Printf("Failed to write to client: %v", err)
				return
			}
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			return
		}
	}
}
