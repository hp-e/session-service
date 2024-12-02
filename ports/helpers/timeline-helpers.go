package helpers

import (
	"log/slog"
	"studio-service/core/models"
)

const (
	IntervalStartCommand         = "interval_start"
	IntervalCompleteCommand      = "interval_complete"
	GroupIntervalStartCommand    = "group_interval_start"
	GroupIntervalCompleteCommand = "group_interval_complete"
)

func BuildEventMapplan(plan *models.SessionPlan, timeline *models.Timeline) (*models.SessionEventMap, error) {

	r := &models.SessionEventMap{
		Groups:    make([]models.TimeEvent[models.SessionIntervalGroup], len(plan.Groups)),
		Timeline:  make([]models.TimeEvent[models.TimelineEvent], 0),
		Intervals: make([]models.TimeEvent[models.SessionInterval], 0),
		OffsetMap: make(map[int][]models.Event),
	}

	// eventMap := make(map[int][]models.Event)

	offset := 0
	intIndex := 0
	if plan != nil && plan.Groups != nil {
		for grpIdx, group := range plan.Groups {

			for i := range group.RepeatCount {
				slog.Info("Repeat count", "count", i)

				groupStartOffset := offset
				// groupEndOffset := groupStartOffset + calculateGroupDuration(group)
				// offset := 0 //groupStartOffset
				groupDuration := 0

				for intIdx, interval := range group.Intervals {

					// intervalStartOffset := offset
					duration := int(interval.DurationValue)
					groupDuration += duration

					r.Intervals = append(r.Intervals, models.TimeEvent[models.SessionInterval]{
						Data:     *interval,
						Offset:   offset,
						Duration: duration,
						Command:  IntervalStartCommand,
						Type:     "interval",
						Index:    intIdx,
					})

					r.OffsetMap[offset] = append(r.OffsetMap[offset], models.Event{
						SessionId:     plan.Id,
						Command:       IntervalStartCommand,
						EventData:     interval,
						Type:          "interval",
						IntervalIndex: intIndex,
						GroupIndex:    grpIdx,
					})
					r.OffsetMap[offset+duration] = append(r.OffsetMap[offset+duration], models.Event{
						SessionId:     plan.Id,
						Command:       IntervalCompleteCommand,
						EventData:     interval,
						Type:          "interval",
						IntervalIndex: intIndex,
						GroupIndex:    grpIdx,
					})
					offset += duration
					intIndex++
				}

				r.OffsetMap[groupStartOffset] = append(r.OffsetMap[groupStartOffset], models.Event{
					SessionId:  plan.Id,
					Command:    GroupIntervalStartCommand,
					EventData:  group,
					Type:       "group",
					GroupIndex: grpIdx,
				})

				groupEndOffset := groupStartOffset + groupDuration
				r.OffsetMap[groupEndOffset] = append(r.OffsetMap[groupEndOffset], models.Event{
					SessionId:  plan.Id,
					Command:    GroupIntervalCompleteCommand,
					EventData:  group,
					Type:       "group",
					GroupIndex: grpIdx,
				})

				r.Groups = append(r.Groups, models.TimeEvent[models.SessionIntervalGroup]{
					Data:     *group,
					Offset:   groupStartOffset,
					Duration: groupDuration,
					Command:  GroupIntervalStartCommand,
					Type:     "group",
					Index:    i,
				})
			}
		}
	}

	if timeline == nil {
		return r, nil
	}

	for _, track := range timeline.Tracks {
		// zero offset for each track
		timelineOffset := 0

		for _, event := range track.Items {

			offset := event.Offset + timelineOffset
			r.OffsetMap[offset] = append(r.OffsetMap[offset], models.Event{
				SessionId: plan.Id,
				Command:   event.EnterCommand,
				EventData: event,
				Type:      "timeline",
			})

			exitOffset := offset + event.Duration

			r.OffsetMap[exitOffset] = append(r.OffsetMap[exitOffset], models.Event{
				SessionId: plan.Id,
				Command:   event.ExitCommand,
				EventData: event,
				Type:      "timeline",
			})

			r.Timeline = append(r.Timeline, models.TimeEvent[models.TimelineEvent]{
				Data:     *event,
				Offset:   offset,
				Duration: event.Duration,
				Command:  event.EnterCommand,
				Type:     "timeline",
				Index:    0,
			})
		}
	}

	return r, nil
}

func BuildEventTimeline(plan *models.SessionPlan, timeline *models.Timeline) map[int][]models.Event {
	eventMap := make(map[int][]models.Event)

	offset := 0

	if plan != nil && plan.Groups != nil {
		for _, group := range plan.Groups {

			for i := range group.RepeatCount {
				slog.Info("Repeat count", "count", i)

				groupStartOffset := offset
				// groupEndOffset := groupStartOffset + calculateGroupDuration(group)
				// offset := 0 //groupStartOffset
				groupDuration := 0

				for _, interval := range group.Intervals {
					// intervalStartOffset := offset
					duration := int(interval.DurationValue)
					groupDuration += duration

					eventMap[offset] = append(eventMap[offset], models.Event{
						SessionId: plan.Id,
						Command:   IntervalStartCommand,
						EventData: interval,
						Type:      "interval",
					})
					eventMap[offset+duration] = append(eventMap[offset+duration], models.Event{
						SessionId: plan.Id,
						Command:   IntervalCompleteCommand,
						EventData: interval,
						Type:      "interval",
					})
					offset += duration
				}

				eventMap[groupStartOffset] = append(eventMap[groupStartOffset], models.Event{
					SessionId: plan.Id,
					Command:   GroupIntervalStartCommand,
					EventData: group,
					Type:      "group",
				})

				groupEndOffset := groupStartOffset + groupDuration
				eventMap[groupEndOffset] = append(eventMap[groupEndOffset], models.Event{
					SessionId: plan.Id,
					Command:   GroupIntervalCompleteCommand,
					EventData: group,
					Type:      "group",
				})
			}
		}
	}

	if timeline == nil {
		return eventMap
	}

	for _, track := range timeline.Tracks {
		// zero offset for each track
		timelineOffset := 0

		for _, event := range track.Items {

			offset := event.Offset + timelineOffset
			eventMap[offset] = append(eventMap[offset], models.Event{
				SessionId: plan.Id,
				Command:   event.EnterCommand,
				EventData: event,
				Type:      "timeline",
			})

			exitOffset := offset + event.Duration

			eventMap[exitOffset] = append(eventMap[exitOffset], models.Event{
				SessionId: plan.Id,
				Command:   event.ExitCommand,
				EventData: event,
				Type:      "timeline",
			})
		}
	}

	return eventMap
}
