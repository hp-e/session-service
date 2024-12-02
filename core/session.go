package core

// Sample sessions creation
// func CreateSampleSessions() []*SessionPlan {
// 	// Create a warmup group
// 	warmupGroup := &SessionIntervalGroup{
// 		Index:       1,
// 		Id:          "warmup-group",
// 		RepeatCount: 1,
// 		Type:        "Warmup",
// 		Title:       stringPointer("Warmup"),
// 		Notes:       stringPointer("Start at a moderate effort"),
// 		Intervals: []*SessionInterval{
// 			{
// 				Index:         1,
// 				Id:            "warmup-interval",
// 				DurationType:  "DurationTime",
// 				DurationValue: 300,
// 				BodyPosition:  &"Seated",
// 				HandPosition:  &Hoods,
// 				Primary: &Target{
// 					Type:       string(Power3),
// 					StartValue: 100,
// 					EndValue:   150,
// 				},
// 			},
// 		},
// 	}

// 	// Create a workout group
// 	workoutGroup := &SessionIntervalGroup{
// 		Index:       2,
// 		Id:          "workout-group",
// 		RepeatCount: 1,
// 		Type:        Workout,
// 		Title:       stringPointer("Workout"),
// 		Notes:       stringPointer("Alternate high and low intensity"),
// 		Intervals: []*SessionInterval{
// 			{
// 				Index:         1,
// 				Id:            "workout-interval-1",
// 				DurationType:  DurationTime,
// 				DurationValue: 120,
// 				BodyPosition:  &Standing,
// 				HandPosition:  &Drops,
// 				Target: []*IntervalTarget{
// 					{
// 						Primary: &Target{
// 							Type:       string(PowerZone),
// 							StartValue: 200,
// 							EndValue:   250,
// 						},
// 					},
// 				},
// 			},
// 			{
// 				Index:         2,
// 				Id:            "workout-interval-2",
// 				DurationType:  DurationTime,
// 				DurationValue: 120,
// 				BodyPosition:  &Seated,
// 				HandPosition:  &Tops,
// 				Target: []*IntervalTarget{
// 					{
// 						Primary: &Target{
// 							Type:       string(PowerZone),
// 							StartValue: 100,
// 							EndValue:   150,
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	// Create a cooldown group
// 	cooldownGroup := &SessionIntervalGroup{
// 		Index:       3,
// 		Id:          "cooldown-group",
// 		RepeatCount: 1,
// 		Type:        Cooldown,
// 		Title:       stringPointer("Cooldown"),
// 		Notes:       stringPointer("Relax and recover"),
// 		Intervals: []*SessionInterval{
// 			{
// 				Index:         1,
// 				Id:            "cooldown-interval",
// 				DurationType:  DurationTime,
// 				DurationValue: 300,
// 				BodyPosition:  &Seated,
// 				HandPosition:  &Hoods,
// 				Target: []*IntervalTarget{
// 					{
// 						Primary: &Target{
// 							Type:       string(Power3),
// 							StartValue: 50,
// 							EndValue:   100,
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}

// 	// Create timeline tracks and events
// 	timelineTrack := &TimelineTrack{
// 		Id:    "track-1",
// 		Title: "Workout Plan Display",
// 		Type:  "plan",
// 		Color: "blue",
// 		Items: []*TimelineEvent{
// 			{
// 				Id:           "show-plan",
// 				EnterCommand: "show_plan",
// 				ExitCommand:  "hide_plan",
// 				Position: &ItemPosition{
// 					ScreenSection: "main",
// 					Position:      1,
// 				},
// 				Properties: map[string]interface{}{
// 					"text": "Workout Plan",
// 				},
// 			},
// 		},
// 	}

// 	// Assemble the session plans
// 	session1 := &SessionPlan{
// 		Id:             "session-1",
// 		DbId:           1,
// 		Title:          "Session 1",
// 		Description:    "A simple warmup, workout, and cooldown session",
// 		Duration:       900, // 15 minutes
// 		Groups:         []*SessionIntervalGroup{warmupGroup, workoutGroup, cooldownGroup},
// 		TimelineTracks: []*TimelineTrack{timelineTrack},
// 	}

// 	session2 := &SessionPlan{
// 		Id:             "session-2",
// 		DbId:           2,
// 		Title:          "Session 2",
// 		Description:    "Another cycling session with the same structure",
// 		Duration:       900, // 15 minutes
// 		Groups:         []*SessionIntervalGroup{warmupGroup, workoutGroup, cooldownGroup},
// 		TimelineTracks: []*TimelineTrack{timelineTrack},
// 	}

// 	return []*SessionPlan{session1, session2}
// }

// func stringPointer(s string) *string {
// 	return &s
// }

// func main() {
// 	// Create sessions
// 	sessions := CreateSampleSessions()

// 	// Print session details for verification
// 	for _, session := range sessions {
// 		fmt.Printf("Session ID: %s\n", session.Id)
// 		fmt.Printf("Title: %s\n", session.Title)
// 		fmt.Printf("Description: %s\n", session.Description)
// 		fmt.Printf("Duration: %d seconds\n", session.Duration)
// 		for _, group := range session.Groups {
// 			fmt.Printf(" Group: %s (%s)\n", *group.Title, group.Type)
// 			for _, interval := range group.Intervals {
// 				fmt.Printf("  Interval: %s (%fs)\n", interval.Id, interval.DurationValue)
// 			}
// 		}
// 		fmt.Println()
// 	}
// }
