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

// Import your models package here
// import "models"

// type Event struct {
// 	RoomId    string      `json:"roomId"`
// 	SessionId string      `json:"sessionId"`
// 	Command   string      `json:"command"`
// 	EventData interface{} `json:"eventData"`
// }

// Session represents a cycling session
// type Session struct {
// 	RoomId    string
// 	SessionId string
// 	Plan      *SessionPlan
// 	Events    chan Event
// }

// // SessionManager manages active sessions
// type SessionManager struct {
// 	mu       sync.Mutex
// 	sessions map[string]*Session
// }

// func NewSessionManager() *SessionManager {
// 	return &SessionManager{
// 		sessions: make(map[string]*Session),
// 	}
// }

// func (sm *SessionManager) AddSession(session *Session) {
// 	sm.mu.Lock()
// 	defer sm.mu.Unlock()
// 	sm.sessions[session.SessionId] = session
// }

// func (sm *SessionManager) GetSession(sessionId string) (*Session, bool) {
// 	sm.mu.Lock()
// 	defer sm.mu.Unlock()
// 	session, exists := sm.sessions[sessionId]
// 	return session, exists
// }

// func (sm *SessionManager) RemoveSession(sessionId string) {
// 	sm.mu.Lock()
// 	defer sm.mu.Unlock()
// 	delete(sm.sessions, sessionId)
// }

// // RunSession handles the session logic
// func (session *Session) RunSession(wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	ticker := time.NewTicker(1 * time.Second)
// 	defer ticker.Stop()

// 	elapsed := 0
// 	totalDuration := session.Plan.Duration

// 	for {
// 		select {
// 		case <-ticker.C:
// 			elapsed++
// 			remaining := totalDuration - elapsed

// 			timeEvent := TimeEvent{
// 				Total: Time{
// 					Duration:  totalDuration,
// 					Elapsed:   elapsed,
// 					Remaining: remaining,
// 				},
// 			}

// 			event := Event{
// 				RoomId:    session.RoomId,
// 				SessionId: session.SessionId,
// 				Command:   "time_update",
// 				EventData: timeEvent,
// 			}

// 			session.Events <- event

// 			if remaining <= 0 {
// 				close(session.Events)
// 				return
// 			}
// 		}
// 	}
// }

// func serveSSE(sm *SessionManager, w http.ResponseWriter, r *http.Request) {
// 	// roomId := r.URL.Query().Get("roomId")
// 	sessionId := r.URL.Query().Get("sessionId")

// 	session, exists := sm.GetSession(sessionId)
// 	if !exists {
// 		http.Error(w, "Session not found", http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "text/event-stream")
// 	w.Header().Set("Cache-Control", "no-cache")
// 	w.Header().Set("Connection", "keep-alive")

// 	for event := range session.Events {
// 		data, err := json.Marshal(event)
// 		if err != nil {
// 			log.Printf("Error marshalling event: %v", err)
// 			continue
// 		}

//			fmt.Fprintf(w, "data: %s\n\n", data)
//			w.(http.Flusher).Flush()
//		}
//	}
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
