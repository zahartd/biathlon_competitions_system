package engine

import (
	"time"
)

type penaltyInterval struct {
	Start time.Time
	End   time.Time
}

type competitorState struct {
	CompetitorID     int
	RegisteredTime   time.Time
	ScheduledStart   time.Time
	ActualStart      time.Time
	NotStarted       bool
	NotFinished      bool
	NotFinishedMsg   string
	LapEndTimes      []time.Time
	PenaltyIntervals []penaltyInterval
	Shots            int
	Hits             int
	lineHits         int
	FinishTime       time.Time
}
