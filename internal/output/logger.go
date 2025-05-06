package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/zahartd/biathlon_competitions_system/internal/events"
	"github.com/zahartd/biathlon_competitions_system/internal/models"
)

type Logger struct {
	w io.Writer
}

func NewLogger(w io.Writer) *Logger {
	return &Logger{w: w}
}

func (l *Logger) Write(event models.Event) {
	timestamp := event.Time.Format(events.TimeLayoutHMSMilli)
	var line string
	switch event.ID {
	case models.EventRegister:
		line = fmt.Sprintf("[%s] The competitor(%d) registered", timestamp, event.CompetitorID)
	case models.EventDraw:
		line = fmt.Sprintf(
			"[%s] The start time for the competitor(%d) was set by a draw to %s",
			timestamp, event.CompetitorID, event.ExtraParams[0],
		)
	case models.EventOnLine:
		line = fmt.Sprintf("[%s] The competitor(%d) is on the start line", timestamp, event.CompetitorID)
	case models.EventStart:
		line = fmt.Sprintf("[%s] The competitor(%d) has started", timestamp, event.CompetitorID)
	case models.EventFiring:
		line = fmt.Sprintf(
			"[%s] The competitor(%d) is on the firing range(%s)",
			timestamp,
			event.CompetitorID,
			event.ExtraParams[0],
		)
	case models.EventHit:
		line = fmt.Sprintf(
			"[%s] The target(%s) has been hit by competitor(%d)",
			timestamp, event.ExtraParams[0], event.CompetitorID,
		)
	case models.EventLeaveFiring:
		line = fmt.Sprintf("[%s] The competitor(%d) left the firing range", timestamp, event.CompetitorID)
	case models.EventPenaltyEnter:
		line = fmt.Sprintf("[%s] The competitor(%d) entered the penalty laps", timestamp, event.CompetitorID)
	case models.EventPenaltyLeave:
		line = fmt.Sprintf("[%s] The competitor(%d) left the penalty laps", timestamp, event.CompetitorID)
	case models.EventLapEnd:
		line = fmt.Sprintf("[%s] The competitor(%d) ended the main lap", timestamp, event.CompetitorID)
	case models.EventNotContinue:
		comment := strings.Join(event.ExtraParams, " ")
		line = fmt.Sprintf(
			"[%s] The competitor(%d) can`t continue: %s",
			timestamp, event.CompetitorID, comment,
		)
	case models.EventDisqualification:
		line = fmt.Sprintf("[%s] The competitor(%d) is disqualified", timestamp, event.CompetitorID)
	case models.EventFinished:
		line = fmt.Sprintf("[%s] The competitor(%d) has finished", timestamp, event.CompetitorID)
	default:
		return
	}

	fmt.Fprintln(l.w, line)
}
