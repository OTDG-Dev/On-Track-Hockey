package data_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/validator"
)

func TestValidateGameEvent(t *testing.T) {

	tests := []struct {
		name             string
		valid            bool
		failedCheckCount int
		ge               data.GameEvent
	}{
		{name: "ValidRequest",
			valid:            true,
			failedCheckCount: 0,
			ge: data.GameEvent{
				Period:       1,
				ClockSeconds: 120,
				EventType:    "goal",
				Situation:    "ev",
				TeamID:       2,
			},
		},
		{name: "ClockNegative-BadSituation",
			valid:            false,
			failedCheckCount: 2,
			ge: data.GameEvent{
				Period:       1,
				ClockSeconds: -4,
				EventType:    "goal",
				Situation:    "notARealSit",
				TeamID:       2,
			},
		},
		{name: "HugeClock",
			valid:            false,
			failedCheckCount: 1,
			ge: data.GameEvent{
				Period:       1,
				ClockSeconds: 25000,
				EventType:    "goal",
				Situation:    "ev",
				TeamID:       2,
			},
		},
		{name: "FailAll",
			valid:            false,
			failedCheckCount: 5,
			ge: data.GameEvent{
				Period:       16,
				ClockSeconds: -99,
				EventType:    "helloWorld",
				Situation:    "helloWorld",
				TeamID:       -2,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := validator.New()
			data.ValidateGameEvent(v, &tc.ge)

			if v.Valid() != tc.valid {
				t.Errorf("expected valid: %t, got :%t", tc.valid, v.Valid())
			}
			if len(v.Errors) != tc.failedCheckCount {
				t.Errorf("expected to fail: %d tests, actually failed :%d", tc.failedCheckCount, len(v.Errors))
			}

		})
	}

}

// a successful case
func TestGameEventModel_Insert(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	model := data.GameEventModel{DB: db}
	event := &data.GameEvent{
		Period:       1,
		ClockSeconds: 120,
		EventType:    "goal",
		Situation:    "ev",
		TeamID:       5,
	}

	rows := sqlmock.NewRows([]string{
		"id",
		"created_at",
		"version",
		"event_number",
	}).AddRow(1, time.Now(), 1, 1)

	mock.ExpectQuery(`INSERT INTO game_events`).
		WithArgs(0, 1, 120, "goal", "ev", 5).
		WillReturnRows(rows)

	err := model.Insert(event)
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}

	if event.ID != 1 {
		t.Fatalf("expected id=1 got %d", event.ID)
	}

	if event.EventNumber != 1 {
		t.Fatalf("expected event_number=1 got %d", event.EventNumber)
	}
}
