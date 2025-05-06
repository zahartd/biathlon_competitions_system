package models

import (
	"fmt"
	"time"
)

func (r ReportRow) Format() string {
	laps := ""
	for i, d := range r.LapTimes {
		laps += fmt.Sprintf("{%s, %.3f}", formatDuration(d), r.LapSpeeds[i])
		if i < len(r.LapTimes)-1 {
			laps += ", "
		}
	}
	pen := fmt.Sprintf("{%s, %.3f}", formatDuration(r.PenaltyTime), r.PenaltySpeed)
	return fmt.Sprintf("[%s] %d [%s] %s %d/%d\n",
		r.Status, r.CompetitorID, laps, pen, r.Hits, r.Shots)
}

func formatDuration(d time.Duration) string {
	ms := d.Milliseconds() % 1000
	s := int(d.Seconds()) % 60
	m := int(d.Minutes())
	return fmt.Sprintf("%02d:%02d.%03d", m, s, ms)
}
