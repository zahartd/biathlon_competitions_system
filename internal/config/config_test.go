package config

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "bad laps type",
			input: `{
                "laps": "bad",
                "lapLen": 100,
                "penaltyLen": 50,
                "firingLines": 1,
                "start": "00:00:00",
                "startDelta": "00:00:30"
            }`,
			wantErr: true,
		},
		{
			name: "bad lapLen type",
			input: `{
                "laps": 1,
                "lapLen": "bad",
                "penaltyLen": 50,
                "firingLines": 1,
                "start": "00:00:00",
                "startDelta": "00:00:30"
            }`,
			wantErr: true,
		},
		{
			name: "bad penaltyLen type",
			input: `{
                "laps": 1,
                "lapLen": 100,
                "penaltyLen": "bad",
                "firingLines": 1,
                "start": "00:00:00",
                "startDelta": "00:00:30"
            }`,
			wantErr: true,
		},
		{
			name: "bad firingLines type",
			input: `{
                "laps": 1,
                "lapLen": 100,
                "penaltyLen": 50,
                "firingLines": "bad",
                "start": "00:00:00",
                "startDelta": "00:00:30"
            }`,
			wantErr: true,
		},
		{
			name: "bad start format",
			input: `{
                "laps": 1,
                "lapLen": 100,
                "penaltyLen": 50,
                "firingLines": 1,
                "start": "bad",
                "startDelta": "00:00:30"
            }`,
			wantErr: true,
		},
		{
			name: "bad startDelta format",
			input: `{
                "laps": 1,
                "lapLen": 100,
                "penaltyLen": 50,
                "firingLines": 1,
                "start": "00:00:00",
                "startDelta": "bad"
            }`,
			wantErr: true,
		},
		{
			name: "all good",
			input: `{
                "laps": 2,
                "lapLen": 3651,
                "penaltyLen": 50,
                "firingLines": 1,
                "start": "09:30:00",
                "startDelta": "00:00:30"
            }`,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var cfg Config
			err := json.Unmarshal([]byte(tc.input), &cfg)
			if tc.wantErr {
				assert.NotNil(t, err, "Expected error")
			} else {
				assert.Nil(t, err, "Expexted without error")
			}
		})
	}
}
