package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Laps        int           // Amount of laps for main distance
	LapLen      int           // Length of each main lap
	PenaltyLen  int           // Length of each penalty lap
	FiringLines int           // Number of firing lines per lap
	Start       time.Time     // Planned start time for the first competitor
	StartDelta  time.Duration // Planned interval between starts
}

const timeForm = "15:04:05"

type rawConfig struct {
	Laps        int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`
}

func (c *Config) UnmarshalJSON(data []byte) error {
	var raw rawConfig
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	c.Laps = raw.Laps
	c.LapLen = raw.LapLen
	c.PenaltyLen = raw.PenaltyLen
	c.FiringLines = raw.FiringLines

	startTime, err := time.Parse(timeForm, raw.Start)
	if err != nil {
		return fmt.Errorf("failed to parse start time %q: %w", raw.Start, err)
	}
	c.Start = startTime

	var h, m, s int
	if _, err := fmt.Sscanf(raw.StartDelta, "%02d:%02d:%02d", &h, &m, &s); err != nil {
		return fmt.Errorf("failed to parse startDelta %q: %w", raw.StartDelta, err)
	}
	c.StartDelta = time.Duration(h)*time.Hour +
		time.Duration(m)*time.Minute +
		time.Duration(s)*time.Second

	return nil
}

func Load(path *string) (Config, error) {
	data, err := os.ReadFile(*path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
