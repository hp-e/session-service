package publishers

import (
	"context"
	"fmt"
	"log/slog"
	"studio-service/core"
	"studio-service/core/models"
)

type screenPublisher struct {
}

func NewScreenPublisher() core.EventPublisher {
	return &screenPublisher{}
}

func (sp *screenPublisher) Publish(ctx context.Context, event models.Event) error {
	msg := fmt.Sprintf("Publishing event: %s", event.Command)
	// d, _ := json.Marshal(event.EventData)

	slog.Info(msg, "room_id", event.RoomId, "session_id", event.SessionId, "command", event.Command)

	return nil
}
