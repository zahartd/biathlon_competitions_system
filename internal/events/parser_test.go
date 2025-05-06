package events

import (
	"fmt"
	"testing"
	"time"

	"github.com/zahartd/biathlon_competitions_system/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestParseEvent(t *testing.T) {
	tests := []struct {
		name           string
		line           string
		expected       models.Event
		expectedErrSub string
	}{
		{
			name: "valid",
			line: "[12:00:00.000] 1 1",
			expected: models.Event{
				Time:         time.Date(0, time.January, 1, 12, 0, 0, 0, time.UTC),
				ID:           models.EventID(1),
				CompetitorID: 1,
				ExtraParams:  []string{},
			},
		},
		{
			name: "valid with extras",
			line: "[08:05:30.123] 1 2 foo bar baz",
			expected: models.Event{
				Time:         time.Date(0, time.January, 1, 8, 5, 30, 123000000, time.UTC),
				ID:           models.EventID(1),
				CompetitorID: 2,
				ExtraParams:  []string{"foo", "bar", "baz"},
			},
		},
		{
			name:           "too few fields",
			line:           "[00:00:00.000] 1",
			expectedErrSub: "invalid event line",
		},
		{
			name:           "bad timestamp format",
			line:           "00:00:00.000 1 1",
			expectedErrSub: "invalid timestamp",
		},
		{
			name:           "bad time format",
			line:           "[0:00:00.000] 1 1",
			expectedErrSub: "invalid timestamp",
		},
		{
			name:           "non-integer event id",
			line:           "[12:00:00.000] bad 1",
			expectedErrSub: "invalid event ID",
		},
		{
			name:           "non-integer competitor id",
			line:           "[12:00:00.000] 1 bad",
			expectedErrSub: "invalid competitor ID",
		},
	}

	parser := NewParser()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := parser.ParseEvent(tc.line)

			if tc.expectedErrSub != "" {
				assert.Contains(t, err.Error(), tc.expectedErrSub,
					fmt.Sprintf("expected error containing %s, but got %s", tc.expectedErrSub, err.Error()))
			} else {
				assert.Nil(t, err, fmt.Sprintf("unexpected error: %v", err))
				assert.Equal(t, actual, tc.expected, "Incorrect parsed event")
			}
		})
	}
}
