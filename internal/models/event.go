package models

import "time"

type EventID int

const (
	// Incoming events
	EventRegister     EventID = iota + 1 // The competitor registered
	EventDraw                            // The start time was set by a draw
	EventOnLine                          // The competitor is on the start line
	EventStart                           // The competitor has started
	EventFiring                          // The competitor is on the firing range
	EventHit                             // The target has been hit
	EventLeaveFiring                     // The competitor left the firing range
	EventPenaltyEnter                    // The competitor entered the penalty laps
	EventPenaltyLeave                    // The competitor left the penalty laps
	EventLapEnd                          // The competitor ended the main lap
	EventNotContinue                     // The competitor can`t continue

	// Outgoing events
	EventDisqualification = 32 // The competitor is disqualified
	EventFinished         = 33 // The competitor has finished
)

type Event struct {
	Time         time.Time
	ID           EventID
	CompetitorID int
	ExtraParams  []string
}
