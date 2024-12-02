package helpers_test

import (
	"studio-service/core/models"
	"studio-service/ports/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildEventTimeline(t *testing.T) {
	// Sample session plan
	sessionPlan := &models.SessionPlan{
		Id:       "test-session",
		Duration: 600,
		Groups: []*models.SessionIntervalGroup{
			{
				Index: 1,
				Id:    "warmup-group",
				Intervals: []*models.SessionInterval{
					{
						Index:         1,
						Id:            "warmup-interval",
						DurationValue: 120,
					},
					{
						Index:         1,
						Id:            "warmup-interval",
						DurationValue: 120,
					},
				},
			},
			{
				Index: 2,
				Id:    "workout-group",
				Intervals: []*models.SessionInterval{
					{
						Index:         1,
						Id:            "workout-interval-1",
						DurationValue: 60,
					},
					{
						Index:         2,
						Id:            "workout-interval-2",
						DurationValue: 60,
					},
					{
						Index:         1,
						Id:            "workout-interval-1",
						DurationValue: 60,
					},
					{
						Index:         2,
						Id:            "workout-interval-2",
						DurationValue: 60,
					},
					{
						Index:         1,
						Id:            "workout-interval-1",
						DurationValue: 60,
					},
					{
						Index:         2,
						Id:            "workout-interval-2",
						DurationValue: 120,
					},
				},
			},
			{
				Index: 1,
				Id:    "cooldown-group",
				Intervals: []*models.SessionInterval{
					{
						Index:         1,
						Id:            "warmup-interval",
						DurationValue: 240,
					},
				},
			},
		},
	}

	// Sample timeline
	timeline := &models.Timeline{
		Tracks: []*models.TimelineTrack{
			{
				Id: "track-1",
				Items: []*models.TimelineEvent{
					{
						Id:           "timeline-event-1",
						Offset:       0,
						Duration:     900,
						EnterCommand: "show_plan",
						ExitCommand:  "hide_plan",
					},
				},
			},
			{
				Id: "track-2",
				Items: []*models.TimelineEvent{
					{
						Id:           "timeline-event-2",
						Offset:       800,
						Duration:     90,
						EnterCommand: "show_overlay",
						ExitCommand:  "hide_overlay",
						ItemId:       "overlay-1",
					},
				},
			},
		},
	}

	// Build the event timeline
	eventTimeline := helpers.BuildEventTimeline(sessionPlan, timeline)

	// Expected results
	expectedTimeline := map[int][]models.Event{
		0: {
			{SessionId: "test-session", Command: "interval_start", EventData: sessionPlan.Groups[0].Intervals[0]},
			{SessionId: "test-session", Command: "group_interval_start", EventData: sessionPlan.Groups[0]},
			{SessionId: "test-session", Command: "show_plan", EventData: timeline.Tracks[0].Items[0]},
		},
		120: {

			{SessionId: "test-session", Command: "interval_complete", EventData: sessionPlan.Groups[0].Intervals[0]},
			{SessionId: "test-session", Command: "interval_start", EventData: sessionPlan.Groups[0].Intervals[1]},
		},
		240: {
			{SessionId: "test-session", Command: "group_interval_complete", EventData: sessionPlan.Groups[0]},
			{SessionId: "test-session", Command: "group_interval_start", EventData: sessionPlan.Groups[1]},
			{SessionId: "test-session", Command: "interval_start", EventData: sessionPlan.Groups[1].Intervals[0]},
		},
		860: {
			{SessionId: "test-session", Command: "group_interval_complete", EventData: sessionPlan.Groups[1]},

			{SessionId: "test-session", Command: "interval_complete", EventData: sessionPlan.Groups[1].Intervals[0]},
			{SessionId: "test-session", Command: "interval_start", EventData: sessionPlan.Groups[1].Intervals[1]},
		},
		540: {
			{SessionId: "test-session", Command: "interval_complete", EventData: sessionPlan.Groups[1].Intervals[1]},
			{SessionId: "test-session", Command: "group_interval_complete", EventData: sessionPlan.Groups[1]},
		},
		800: {
			{SessionId: "test-session", Command: "show_overlay", EventData: timeline.Tracks[1].Items[0]},
		},
		890: {
			{SessionId: "test-session", Command: "hide_overlay", EventData: timeline.Tracks[1].Items[0]},
		},
		900: {
			{SessionId: "test-session", Command: "hide_plan", EventData: timeline.Tracks[0].Items[0]},
		},
	}

	// Assertions
	assert.Equal(t, len(expectedTimeline), len(eventTimeline), "Event timeline length mismatch")

	for offset, events := range expectedTimeline {
		assert.Contains(t, eventTimeline, offset, "Offset missing in event timeline")
		assert.Equal(t, events, eventTimeline[offset], "Events mismatch at offset %d", offset)
	}
}
