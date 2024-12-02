package models

// Time (duration) will always be in seconds
// Distance will always be in meters
// Speed will always be in kilometers per hour
// Power will always be in watts
// Heart rate will always be in beats per minute
// Cadence will always be in revolutions per minute
// Calories will always be in kilocalories

type TargetType string

type IntervalGroupType string
type DurationType string
type BodyPosition string
type HandPosition string

const (
	Power3        TargetType = "power3s"
	Power5        TargetType = "power5s"
	PowerZone     TargetType = "powerZone"
	PowerFtp      TargetType = "powerFtp"
	HeartRateZone TargetType = "heartRateZone"
	Cadence       TargetType = "cadence"
	Calories      TargetType = "calories"
	Speed         TargetType = "speed"
	Distance      TargetType = "distance"

	Warmup         IntervalGroupType = "warmup"
	Workout        IntervalGroupType = "workout"
	Cooldown       IntervalGroupType = "cooldown"
	Rest           IntervalGroupType = "rest"
	Recovery       IntervalGroupType = "recovery"
	FTPTest        IntervalGroupType = "ftp_test"
	Competition    IntervalGroupType = "competition"
	GroupChallenge IntervalGroupType = "groupchallenge"

	DurationTime      DurationType = "time"
	DurationDistance  DurationType = "distance"
	DurationCalories  DurationType = "calories"
	DurationHeartRate DurationType = "heartRate"
	DurationPower     DurationType = "power"

	Seated   BodyPosition = "seated"
	Standing BodyPosition = "standing"
	Hover    BodyPosition = "hover"
	Sprint   BodyPosition = "sprint"
	Climb    BodyPosition = "climb"
	Attack   BodyPosition = "attack"

	Hoods      HandPosition = "hoods"
	Drops      HandPosition = "drops"
	Tops       HandPosition = "tops"
	Aerobars   HandPosition = "aerobars"
	Hooks      HandPosition = "hooks"
	Extensions HandPosition = "extensions"
)

type Session struct {
	DbId        int          `json:"dbId"`
	Id          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	RoomId      string       `json:"roomId"`
	Plan        *SessionPlan `json:"plan"`
	Timeline    *Timeline    `json:"timeline"`
}

type SessionEvent struct {
	RoomId                string
	SessionId             string
	Command               string
	Type                  string
	CurrentInterval       *SessionInterval
	CurrentIntervalGroup  *SessionIntervalGroup
	NextInterval          *SessionInterval
	NextIntervalGroup     *SessionIntervalGroup
	PreviousInterval      *SessionInterval
	PreviousIntervalGroup *SessionIntervalGroup
	TimeEvent             *TimeTick
	TimelineEvent         *TimelineEvent
	// Event                *Event
}

type SessionEventMap struct {
	Groups    []TimeEvent[SessionIntervalGroup]
	Intervals []TimeEvent[SessionInterval]
	Timeline  []TimeEvent[TimelineEvent]
	OffsetMap map[int][]Event
}

type TimeEvent[T any] struct {
	Data     T
	Duration int
	Offset   int
	Command  string
	Type     string
	Index    int
	// OriginalIndex int
}

type Event struct {
	RoomId        string `json:"roomId"`
	SessionId     string `json:"sessionId"`
	Command       string `json:"command"`
	Type          string `json:"type"`
	EventData     any    `json:"eventData"` // TimelineEvent, TimeEvent, IntervalGroup, Interval
	GroupIndex    int    `json:"groupIndex"`
	IntervalIndex int    `json:"intervalIndex"`
}

type Timeline struct {
	DbId        int              `json:"dbId"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Tracks      []*TimelineTrack `json:"tracks"`
}

type TimelineTrack struct {
	Id    string           `json:"id"`
	Items []*TimelineEvent `json:"items"`
	Title string           `json:"title"`
	Type  string           `json:"type"`
	Color string           `json:"color"`
	Icon  string           `json:"icon"`
}

type TimelineEvent struct {
	Id           string                 `json:"id"`
	Offset       int                    `json:"offset"`
	Duration     int                    `json:"duration"`
	EnterCommand string                 `json:"enterCommand"`
	ExitCommand  string                 `json:"exitCommand"`
	ItemId       string                 `json:"itemId"`
	TrasitionIn  *Transition            `json:"transitionIn"`
	TrasitionOut *Transition            `json:"transitionOut"`
	Position     *ItemPosition          `json:"position"`
	Properties   map[string]interface{} `json:"properties"`
}

type Transition struct {
	Type     string `json:"type"`
	Duration int    `json:"duration"` // duration in seconds
}

type ItemPosition struct {
	ScreenSection string `json:"screenSection"`
	Position      int    `json:"position"`
	Top           int    `json:"top"`
	Left          int    `json:"left"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
}

type SessionPlan struct {
	Id             string                  `json:"id"`
	DbId           int                     `json:"dbId"`
	Title          string                  `json:"title"`
	Description    string                  `json:"description"`
	Notes          string                  `json:"notes"`
	Duration       int                     `json:"duration"`
	Groups         []*SessionIntervalGroup `json:"groups"`
	TimelineTracks []*TimelineTrack        `json:"timelineTracks"`
}

type SessionInterval struct {
	Index         int           `json:"index"`
	Id            string        `json:"id"`
	Title         string        `json:"title"`
	Notes         string        `json:"notes"`
	Primary       *Target       `json:"primary"`
	Secondary     *Target       `json:"secondary"`
	DurationType  DurationType  `json:"durationType"`
	DurationValue float64       `json:"durationValue"`
	BodyPosition  *BodyPosition `json:"bodyPosition"`
	HandPosition  *HandPosition `json:"handPosition"`
	Resistance    int           `json:"resistance"`
}

type SessionIntervalGroup struct {
	Index       int                `json:"index"`
	Id          string             `json:"id"`
	RepeatCount int                `json:"repeatCount"`
	Intervals   []*SessionInterval `json:"intervals"`
	Title       string             `json:"title"`
	Notes       string             `json:"notes"`
	Type        TargetType         `json:"type"`
}

type IntervalTarget struct {
	Primary   *Target `json:"primary"`
	Secondary *Target `json:"secondary"`
}

type Target struct {
	Type       string  `json:"type"`
	StartValue float64 `json:"startValue"`
	EndValue   float64 `json:"endValue"`

	PercentageBaseType string `json:"percentageBaseType"`
}

type Time struct {
	Duration  int `json:"duration"`  // total duration for session in seconds
	Elapsed   int `json:"elapsed"`   // elapsed time in seconds
	Remaining int `json:"remaining"` // remaining time in seconds

}
type TimeTick struct {
	Total                 Time  `json:"total"`
	CurrentInterval       Time  `json:"currentInterval"`
	CurrentIntervalGroup  Time  `json:"currentIntervalGroup"`
	PreviousInterval      *Time `json:"previousInterval"`
	PreviousIntervalGroup *Time `json:"previousIntervalGroup"`
	NextInterval          *Time `json:"nextInterval"`
	NextIntervalGroup     *Time `json:"nextIntervalGroup"`
}

func CreateSessionPlan() *SessionPlan {
	return &SessionPlan{
		Id:             "session.default",
		DbId:           0,
		Title:          "Default Session",
		Description:    "Default Session Description",
		Notes:          "Default Session Notes",
		Duration:       2700,
		Groups:         make([]*SessionIntervalGroup, 0),
		TimelineTracks: make([]*TimelineTrack, 0),
	}
}

func CreateSampleSessions() []*SessionPlan {
	// Create a warmup group
	warmupGroup := &SessionIntervalGroup{
		Index:       1,
		Id:          "warmup-group",
		RepeatCount: 1,
		Type:        "warmup",
		Title:       "Warmup",
		Notes:       "Start at a moderate effort",
		Intervals: []*SessionInterval{
			{
				Index:         1,
				Id:            "warmup-interval",
				DurationType:  DurationTime,
				DurationValue: 60,
				// BodyPosition:  stringPointer(Seated),
				// HandPosition:  Hoods,

				Primary: &Target{
					Type:       string(Power3),
					StartValue: 100,
					EndValue:   150,
				},
			},
		},
	}

	// Create a workout group
	workoutGroup := &SessionIntervalGroup{
		Index:       2,
		Id:          "workout-group",
		RepeatCount: 1,
		Type:        "workout",
		Title:       "Workout",
		Notes:       "Alternate high and low intensity",
		Intervals: []*SessionInterval{
			{
				Index:         1,
				Id:            "workout-interval-1",
				DurationType:  DurationTime,
				DurationValue: 60,
				// BodyPosition:  &Standing,
				// HandPosition:  &Drops,

				Primary: &Target{
					Type:       string(PowerZone),
					StartValue: 200,
					EndValue:   250,
				},
			},
			{
				Index:         2,
				Id:            "workout-interval-2",
				DurationType:  DurationTime,
				DurationValue: 60,
				// BodyPosition:  &Seated,
				// HandPosition:  &Tops,

				Primary: &Target{
					Type:       string(PowerZone),
					StartValue: 100,
					EndValue:   150,
				},
			},
		},
	}

	// Create a cooldown group
	cooldownGroup := &SessionIntervalGroup{
		Index:       3,
		Id:          "cooldown-group",
		RepeatCount: 1,
		Type:        "Cooldown",
		Title:       "Cooldown",
		Notes:       "Relax and recover",
		Intervals: []*SessionInterval{
			{
				Index:         1,
				Id:            "cooldown-interval",
				DurationType:  DurationTime,
				DurationValue: 120,
				// BodyPosition:  &Seated,
				// HandPosition:  &Hoods,

				Primary: &Target{
					Type:       string(Power3),
					StartValue: 50,
					EndValue:   100,
				},
			},
		},
	}

	// Create timeline tracks and events
	timelineTrack := &TimelineTrack{
		Id:    "track-1",
		Title: "Workout Plan Display",
		Type:  "plan",
		Color: "blue",
		Items: []*TimelineEvent{
			{
				Id:           "show-plan",
				EnterCommand: "show_plan",
				ExitCommand:  "hide_plan",
				Position: &ItemPosition{
					ScreenSection: "main",
					Position:      1,
				},
				Properties: map[string]interface{}{
					"text": "Workout Plan",
				},
			},
		},
	}

	// Assemble the session plans
	session1 := &SessionPlan{
		Id:             "session-1",
		DbId:           1,
		Title:          "Session 1",
		Description:    "A simple warmup, workout, and cooldown session",
		Duration:       900, // 15 minutes
		Groups:         []*SessionIntervalGroup{warmupGroup, workoutGroup, cooldownGroup},
		TimelineTracks: []*TimelineTrack{timelineTrack},
	}

	session2 := &SessionPlan{
		Id:             "session-2",
		DbId:           2,
		Title:          "Session 2",
		Description:    "Another cycling session with the same structure",
		Duration:       900, // 15 minutes
		Groups:         []*SessionIntervalGroup{warmupGroup, workoutGroup, cooldownGroup},
		TimelineTracks: []*TimelineTrack{timelineTrack},
	}

	return []*SessionPlan{session1, session2}
}

func stringPointer(s string) *string {
	return &s
}
