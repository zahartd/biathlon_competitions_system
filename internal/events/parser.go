package events

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zahartd/biathlon_competitions_system/internal/models"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

const (
	TimeLayoutHMSMilli = "15:04:05.000"
	timestampLen       = len(TimeLayoutHMSMilli) + 2 // [15:04:05.000]
)

func (p *Parser) ParseEvent(line string) (models.Event, error) {
	tokens := strings.Fields(line)
	if len(tokens) < 3 {
		return models.Event{},
			fmt.Errorf("invalid event line, expected [time] eventID competitorID extraParams")
	}

	timestampToken := tokens[0]
	if len(timestampToken) != timestampLen ||
		timestampToken[0] != '[' ||
		timestampToken[len(timestampToken)-1] != ']' {
		return models.Event{},
			fmt.Errorf("invalid timestamp, expected [%s], but received: %s", TimeLayoutHMSMilli, timestampToken)
	}
	timestampStr := timestampToken[1 : len(timestampToken)-1]
	timestamp, err := time.Parse(TimeLayoutHMSMilli, timestampStr)
	if err != nil {
		return models.Event{},
			fmt.Errorf("invalid timestamp, expected [%s], but received: %s", TimeLayoutHMSMilli, timestampToken)
	}

	eventID, err := strconv.Atoi(tokens[1])
	if err != nil {
		return models.Event{}, fmt.Errorf("invalid event ID %s: %w", tokens[1], err)
	}

	competitorID, err := strconv.Atoi(tokens[2])
	if err != nil {
		return models.Event{}, fmt.Errorf("invalid competitor ID %s: %w", tokens[2], err)
	}

	extraParams := tokens[3:] // maybe empty

	return models.Event{
		Time:         timestamp,
		ID:           models.EventID(eventID),
		CompetitorID: competitorID,
		ExtraParams:  extraParams,
	}, nil
}
