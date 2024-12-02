package sessions

import (
	"context"
	"log"
	"log/slog"
	"studio-service/core"
	"studio-service/core/models"
	"studio-service/ports/helpers"
	"sync"
	"time"
)

type SessionServiceImpl struct {
	mu        sync.Mutex
	sessions  map[string]*models.Session
	publisher core.EventPublisher
}

func NewSessionManager(publisher core.EventPublisher) core.SessionManager {
	return &SessionServiceImpl{
		sessions:  make(map[string]*models.Session),
		publisher: publisher,
	}
}

func (s *SessionServiceImpl) Run(ctx context.Context, session *models.Session) {
	s.mu.Lock()
	s.sessions[session.Id] = session
	s.mu.Unlock()

	log.Printf("Starting session: %s", session.Id)

	go s.runSession(ctx, session)
}

// , eventTimeline map[int][]models.Event
func (s *SessionServiceImpl) runSession(ctx context.Context, session *models.Session) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	mapper, _ := helpers.BuildEventMapplan(session.Plan, session.Timeline)
	// for k, v := range mapper.OffsetMap {
	// 	slog.Info("Event Timeline", "offset", k, "commands", len(v))
	// }

	for _, grp := range mapper.Groups {
		slog.Info("Group", "Duration", grp.Duration, "Offset", grp.Offset, "Type", grp.Type)
	}

	for x, i := range mapper.Intervals {
		slog.Info("Interval", "Idx", x, "Duration", i.Duration, "Offset", i.Offset, "Type", i.Type)
	}
	totalDuration := session.Plan.Duration
	elapsed := 0
	timerEvent := models.TimeTick{
		Total: models.Time{
			Duration:  totalDuration,
			Elapsed:   elapsed,
			Remaining: totalDuration - elapsed,
		},
		CurrentInterval: models.Time{
			Duration:  0,
			Elapsed:   0,
			Remaining: 0,
		},
		CurrentIntervalGroup: models.Time{
			Duration:  0,
			Elapsed:   0,
			Remaining: 0,
		},
		NextInterval: &models.Time{
			Duration:  0,
			Elapsed:   0,
			Remaining: 0,
		},
		NextIntervalGroup: &models.Time{
			Duration:  0,
			Elapsed:   0,
			Remaining: 0,
		},
		PreviousInterval: &models.Time{
			Duration:  0,
			Elapsed:   0,
			Remaining: 0,
		},
		PreviousIntervalGroup: &models.Time{
			Duration:  0,
			Elapsed:   0,
			Remaining: 0,
		},
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("Session %s stopped", session.Id)
			s.cleanupSession(session.Id)
			return
		case <-ticker.C:

			if elapsed > totalDuration {
				log.Printf("Session %s completed", session.Id)
				s.cleanupSession(session.Id)
				return
			}

			timerEvent.Total.Elapsed = elapsed
			timerEvent.Total.Remaining = totalDuration - elapsed
			timerEvent.CurrentInterval.Elapsed += 1
			timerEvent.CurrentInterval.Remaining -= 1
			timerEvent.CurrentIntervalGroup.Elapsed += 1
			timerEvent.CurrentIntervalGroup.Remaining -= 1

			// slog.Info("Session running", "session_id", session.Id, "elapsed", elapsed, "remaining", totalDuration-elapsed)
			if events, exists := mapper.OffsetMap[elapsed]; exists {
				for _, event := range events {
					err := s.publisher.Publish(ctx, event)

					switch event.Command {
					case "interval_start":

						timerEvent.CurrentInterval.Duration = mapper.Intervals[event.IntervalIndex].Duration //int(event.EventData.(models.SessionInterval).DurationValue)
						timerEvent.CurrentInterval.Elapsed = 0
						timerEvent.CurrentInterval.Remaining = timerEvent.CurrentInterval.Duration //int(event.EventData.(models.SessionInterval).DurationValue)

						if event.IntervalIndex < len(mapper.Intervals)-1 {
							timerEvent.NextInterval.Duration = mapper.Intervals[event.IntervalIndex+1].Duration //int(mapper.Intervals[event.IntervalIndex+1].DurationValue),
							timerEvent.NextInterval.Elapsed = 0
							timerEvent.NextInterval.Remaining = mapper.Intervals[event.IntervalIndex+1].Duration //int(mapper.Intervals[event.IntervalIndex+1].DurationValue),
						}

					case "group_interval_start":
						// timerEvent.CurrentIntervalGroup.Duration=  int(event.EventData.(models.SessionIntervalGroup).DurationValue)
						timerEvent.CurrentIntervalGroup.Elapsed = 0
						// timerEvent.CurrentIntervalGroup.Remaining= int(event.EventData.(models.SessionIntervalGroup).DurationValue)
						if event.GroupIndex < len(mapper.Groups)-1 {
							timerEvent.NextIntervalGroup.Duration = mapper.Groups[event.GroupIndex+1].Duration //int(mapper.Groups[event.GroupIndex+1].DurationValue)
							timerEvent.NextIntervalGroup.Elapsed = 0
							timerEvent.NextIntervalGroup.Remaining = mapper.Groups[event.GroupIndex+1].Duration //int(mapper.Groups[event.GroupIndex+1].DurationValue)
						}
					}

					if err != nil {
						log.Printf("Failed to publish event: %v", err)
					}
				}
			}

			// err := s.publisher.Publish(ctx, models.Event{
			// 	RoomId:    session.RoomId,
			// 	SessionId: session.Id,
			// 	Command:   fmt.Sprintf("time_tick_%d", timerEvent.Total.Elapsed),
			// 	EventData: timerEvent,
			// })

			// if err != nil {
			// 	log.Printf("Failed to publish time update event: %v", err)
			// }
			slog.Info("Timer Tick", "elapsed", timerEvent.Total.Elapsed,
				"remaining", timerEvent.Total.Remaining,
				"CI e", timerEvent.CurrentInterval.Elapsed,
				"CI r", timerEvent.CurrentInterval.Remaining,
				"CIG e", timerEvent.CurrentIntervalGroup.Elapsed,
				"NI", timerEvent.NextInterval.Duration,
				"NG", timerEvent.NextIntervalGroup.Duration,
			// "N I Grp",
			// timerEvent.NextIntervalGroup.Elapsed, "P I",
			// timerEvent.PreviousInterval.Elapsed, "P I Grp",
			// timerEvent.PreviousIntervalGroup.Elapsed,
			)

			elapsed++
		}
	}
}

func (s *SessionServiceImpl) cleanupSession(sessionId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionId)
}
