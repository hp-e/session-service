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
