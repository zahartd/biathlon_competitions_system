package engine

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/zahartd/biathlon_competitions_system/internal/config"
	"github.com/zahartd/biathlon_competitions_system/internal/events"
	"github.com/zahartd/biathlon_competitions_system/internal/models"
	"github.com/zahartd/biathlon_competitions_system/internal/output"
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

type Engine struct {
	cfg          config.Config
	states       map[int]*competitorState
	resultLogger *output.Logger
}

func NewEngine(cfg config.Config, resultLogger *output.Logger) *Engine {
	return &Engine{
		cfg:          cfg,
		states:       make(map[int]*competitorState),
		resultLogger: resultLogger,
	}
}

func (e *Engine) ProcessEvent(event models.Event) error {
	state, ok := e.states[event.CompetitorID]
	if !ok {
		state = &competitorState{CompetitorID: event.CompetitorID}
		e.states[event.CompetitorID] = state
	}

	e.resultLogger.Write(event)

	switch event.ID {
	case models.EventRegister:
		state.RegisteredTime = event.Time
	case models.EventDraw:
		scheduled, err := time.Parse(events.TimeLayoutHMSMilli, event.ExtraParams[0])
		if err != nil {
			return fmt.Errorf("invalid start time for competitor %d: %w",
				event.CompetitorID, err)
		}
		state.ScheduledStart = scheduled
	case models.EventOnLine:
		// no op
	case models.EventStart:
		state.ActualStart = event.Time
	case models.EventFiring:
		state.lineHits = 0
	case models.EventHit:
		state.lineHits++
	case models.EventLeaveFiring:
		state.Shots += 5
		state.Hits += state.lineHits
	case models.EventPenaltyEnter:
		state.PenaltyIntervals = append(state.PenaltyIntervals, penaltyInterval{Start: event.Time})
	case models.EventPenaltyLeave:
		if len(state.PenaltyIntervals) > 0 {
			state.PenaltyIntervals[len(state.PenaltyIntervals)-1].End = event.Time
		} else {
			return fmt.Errorf("invalid PenaltyLeave event without enter for competitor %d", event.CompetitorID)
		}
	case models.EventLapEnd:
		state.LapEndTimes = append(state.LapEndTimes, event.Time)
		if len(state.LapEndTimes) == e.cfg.Laps {
			finish := models.Event{
				Time:         event.Time,
				ID:           models.EventFinished,
				CompetitorID: event.CompetitorID,
			}
			state.FinishTime = event.Time
			e.resultLogger.Write(finish)
		}
	case models.EventNotContinue:
		state.NotFinished = true
		state.NotFinishedMsg = strings.Join(event.ExtraParams, " ")
	default:
		log.Printf("Unknown eventID=%d for competitor %d", event.ID, event.CompetitorID)
	}
	return nil
}

func (e *Engine) Finilize() {
	for cid, st := range e.states {
		if st.ActualStart.IsZero() && !st.NotFinished {
			disqualification := models.Event{
				Time:         st.ScheduledStart,
				ID:           models.EventDisqualification,
				CompetitorID: cid,
			}
			e.resultLogger.Write(disqualification)
		}
	}
}

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

func (e *Engine) GetReport() []ReportRow {
	var rows []ReportRow
	for _, state := range e.states {
		row := ReportRow{
			CompetitorID:   state.CompetitorID,
			Hits:           state.Hits,
			Shots:          state.Shots,
			ScheduledStart: state.ScheduledStart,
		}

		switch {
		case state.NotFinished:
			row.Status = "NotFinished"
		case state.NotStarted:
			row.Status = "NotStarted"
		default:
			row.Status = "Finished"
		}

		prev := state.ScheduledStart
		for _, end := range state.LapEndTimes {
			dur := end.Sub(prev)
			row.LapTimes = append(row.LapTimes, dur)
			row.LapSpeeds = append(row.LapSpeeds, float64(e.cfg.LapLen)/dur.Seconds())
			prev = end
		}

		var totalPen time.Duration
		for _, iv := range state.PenaltyIntervals {
			totalPen += iv.End.Sub(iv.Start)
		}
		row.PenaltyTime = totalPen
		penCount := row.Shots - row.Hits
		if totalPen > 0 && penCount > 0 {
			row.PenaltySpeed = float64(e.cfg.PenaltyLen*penCount) / totalPen.Seconds()
		}

		rows = append(rows, row)
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].ScheduledStart.Before(rows[j].ScheduledStart)
	})
	return rows
}
