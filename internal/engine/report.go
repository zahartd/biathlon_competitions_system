package engine

import (
	"fmt"
	"strings"
	"time"
)

type ReportRow struct {
	CompetitorID   int
	Status         string
	LapTimes       []time.Duration
	LapSpeeds      []float64
	PenaltyTime    time.Duration
	PenaltySpeed   float64
	Hits           int
	Shots          int
	ScheduledStart time.Time // aux info for sorting, not for report
}

func (r ReportRow) Format() string {
	var lapStrs []string
	for i, d := range r.LapTimes {
		lapStrs = append(lapStrs, fmt.Sprintf("{%s, %.3f}", formatDuration(d), r.LapSpeeds[i]))
	}
	laps := strings.Join(lapStrs, ", ")
	penStr := fmt.Sprintf("{%s, %.3f}", formatDuration(r.PenaltyTime), r.PenaltySpeed)
	return fmt.Sprintf("[%s] %d [%s] %s %d/%d\n",
		r.Status, r.CompetitorID, laps, penStr, r.Hits, r.Shots)
}

func formatDuration(d time.Duration) string {
	ms := d.Milliseconds() % 1000
	s := int(d.Seconds()) % 60
	m := int(d.Minutes())
	return fmt.Sprintf("%02d:%02d.%03d", m, s, ms)
}
