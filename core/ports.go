package core

import (
	"context"
	"studio-service/core/models"
)

type EventPublisher interface {
	Publish(ctx context.Context, event models.Event) error
}

type SessionManager interface {
	Run(ctx context.Context, session *models.Session)
}

type Publisher interface {
	Publish(ctx context.Context, data map[string]any) error
}
