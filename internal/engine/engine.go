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
		state = &competitorState{
			CompetitorID:     event.CompetitorID,
			LapEndTimes:      make([]time.Time, 0, e.cfg.Laps),
			PenaltyIntervals: make([]penaltyInterval, 0, e.cfg.Laps*e.cfg.FiringLines),
		}
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

func (e *Engine) Finalize() {
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
			lapTime := end.Sub(prev)
			row.LapTimes = append(row.LapTimes, lapTime)
			row.LapSpeeds = append(row.LapSpeeds, float64(e.cfg.LapLen)/lapTime.Seconds()) // metr / sec
			prev = end
		}

		var totalPen time.Duration
		for _, interval := range state.PenaltyIntervals {
			totalPen += interval.End.Sub(interval.Start)
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
