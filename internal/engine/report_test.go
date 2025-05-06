package engine

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		d        time.Duration
		expected string
	}{
		{"zero", 0, "00:00.000"},
		{"seconds+ms", 12*time.Second + 123*time.Millisecond, "00:12.123"},
		{"minutes+seconds+ms", 10*time.Minute + 9*time.Second + 89*time.Millisecond, "10:09.089"},
		{"hours+minutes+second+ms", 2*time.Hour + 1*time.Minute + 2*time.Second + 111*time.Millisecond, "121:02.111"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := formatDuration(tc.d)
			assert.Equal(t, actual, tc.expected, fmt.Sprintf("formatDuration(%v) = %q, but expected %q", tc.d, actual, tc.expected))
		})
	}
}

func TestReportRow_Format(t *testing.T) {
	tests := []struct {
		name     string
		row      ReportRow
		expected string
	}{
		{
			name: "multiple laps and one penalty",
			row: ReportRow{
				CompetitorID:   1,
				Status:         "NotFinished",
				LapTimes:       []time.Duration{10*time.Second + 500*time.Millisecond, 12 * time.Second},
				LapSpeeds:      []float64{1000.0 / 10.5, 1000.0 / 12},
				PenaltyTime:    4*time.Second + 200*time.Millisecond,
				PenaltySpeed:   50.0 / 2.2,
				Hits:           4,
				Shots:          5,
				ScheduledStart: time.Date(0, time.January, 1, 9, 32, 0, 23*1e6, time.UTC),
			},
			expected: fmt.Sprintf(
				"[NotFinished] 1 [{00:10.500, %.3f}, {00:12.000, %.3f}] {00:04.200, %.3f} 4/5\n",
				1000.0/10.5, 1000.0/12, 50.0/2.2,
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.row.Format()
			assert.Equal(t, tc.expected, got)
		})
	}
}
