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
	cfg    config.Config
	states map[int]*models.CompetitorState
	logger *output.Logger
}

func NewEngine(cfg config.Config, logger *output.Logger) *Engine {
	return &Engine{
		cfg:    cfg,
		states: make(map[int]*models.CompetitorState),
		logger: logger,
	}
}

func (e *Engine) ProcessEvent(event models.Event) error {
	state, ok := e.states[event.CompetitorID]
	if !ok {
		e.states[event.CompetitorID] = &models.CompetitorState{
			CompetitorID: event.CompetitorID,
		}
	}

	e.logger.Write(event)

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
		{
		}
	case models.EventStart:
		state.ActualStart = event.Time
	case models.EventFiring:
		{
		}
	case models.EventHit:
		state.Shots++
		state.Hits++
	case models.EventLeaveFiring:
		{
		}
	case models.EventPenaltyEnter:
		state.PenaltyStart = event.Time
	case models.EventPenaltyLeave:
		state.PenaltyEnd = event.Time
	case models.EventLapEnd:
		state.LapEndTimes = append(state.LapEndTimes, event.Time)
		if len(state.LapEndTimes) == e.cfg.Laps {
			finish := models.Event{
				Time:         event.Time,
				ID:           models.EventFinished,
				CompetitorID: event.CompetitorID,
			}
			state.FinishTime = event.Time
			e.logger.Write(finish)
		}
	case models.EventNotContinue:
		state.NotFinished = true
		state.NotFinishedMsg = strings.Join(event.ExtraParams, " ")
	default:
		log.Printf("Uknown eventID = %d: %v", event.ID, event)
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
			e.logger.Write(disqualification)
		}
	}
}

func (e *Engine) GetReport() []models.ReportRow {
	var rows []models.ReportRow
	for _, st := range e.states {
		row := models.ReportRow{
			CompetitorID: st.CompetitorID,
			Hits:         st.Hits,
			Shots:        st.Shots,
			StartTime:    st.ActualStart,
		}

		switch {
		case st.NotFinished:
			row.Status = "NotFinished"
		case st.ActualStart.IsZero():
			row.Status = "NotStarted"
		default:
			row.Status = "Finished"
		}

		prev := st.ActualStart
		if prev.IsZero() {
			prev = st.ScheduledStart
		}
		for _, end := range st.LapEndTimes {
			dur := end.Sub(prev)
			row.LapTimes = append(row.LapTimes, dur)
			row.LapSpeeds = append(row.LapSpeeds, float64(e.cfg.LapLen)/dur.Seconds())
			prev = end
		}

		if !st.PenaltyStart.IsZero() && !st.PenaltyEnd.IsZero() {
			penDur := st.PenaltyEnd.Sub(st.PenaltyStart)
			row.PenaltyTime = penDur
			row.PenaltySpeed = float64(e.cfg.PenaltyLen*len(st.LapEndTimes)) / penDur.Seconds()
		}

		rows = append(rows, row)
	}

	sort.Slice(rows, func(i, j int) bool {
		ri, rj := rows[i], rows[j]
		ti := ri.StartTime
		tj := rj.StartTime
		return ti.Before(tj)
	})
	return rows
}
