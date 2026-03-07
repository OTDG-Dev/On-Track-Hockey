package data

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGameEventModel_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	model := GameEventModel{DB: db}

	tests := []struct {
		name    string
		event   *GameEvent
		mockFn  func()
		wantErr bool
	}{
		{
			name: "successful insert",
			event: &GameEvent{
				GameID:    1,
				Period:    1,
				Clock:     "00:14:34",
				EventType: "goal",
				TeamID:    2,
			},
			mockFn: func() {
				mock.ExpectQuery(`INSERT INTO game_events`).
					WithArgs(1, 1, "00:14:34", "goal", 2).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "version"}).
						AddRow(1, time.Now(), 1))
			},
			wantErr: false,
		},
		{
			name: "database error",
			event: &GameEvent{
				GameID:    1,
				Period:    1,
				Clock:     "00:14:34",
				EventType: "goal",
				TeamID:    2,
			},
			mockFn: func() {
				mock.ExpectQuery(`INSERT INTO game_events`).
					WithArgs(1, 1, "00:14:34", "goal", 2).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := model.Insert(tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGameEventModel_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	model := GameEventModel{DB: db}
	now := time.Now()

	tests := []struct {
		name      string
		id        int
		mockFn    func()
		wantEvent *GameEvent
		wantErr   error
	}{
		{
			name: "existing event",
			id:   1,
			mockFn: func() {
				mock.ExpectQuery(`SELECT (.+) FROM game_events WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "game_id", "period", "clock", "event_type", "team_id", "created_at", "version"}).
						AddRow(1, 101, 1, "00:14:34", "goal", 2, now, 1))
			},
			wantEvent: &GameEvent{
				ID:        1,
				GameID:    101,
				Period:    1,
				Clock:     "00:14:34",
				EventType: "goal",
				TeamID:    2,
				CreatedAt: now,
				Version:   1,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			id:   999,
			mockFn: func() {
				mock.ExpectQuery(`SELECT (.+) FROM game_events WHERE id = \$1`).
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			wantEvent: nil,
			wantErr:   ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			event, err := model.Get(tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantEvent != nil {
				if event.ID != tt.wantEvent.ID {
					t.Errorf("Get() event.ID = %v, want %v", event.ID, tt.wantEvent.ID)
				}
			}
		})
	}
}

func TestGameEventModel_ListByGame(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	model := GameEventModel{DB: db}
	now := time.Now()

	tests := []struct {
		name       string
		gameID     int
		mockFn     func()
		wantCount  int
		wantErr    bool
	}{
		{
			name:   "events exist",
			gameID: 101,
			mockFn: func() {
				mock.ExpectQuery(`SELECT (.+) FROM game_events ge INNER JOIN teams`).
					WithArgs(101).
					WillReturnRows(sqlmock.NewRows([]string{"id", "game_id", "period", "clock", "event_type", "team_id", "created_at", "version", "short_name"}).
						AddRow(1, 101, 1, "00:14:34", "goal", 2, now, 1, "Team A").
						AddRow(2, 101, 2, "00:05:22", "penalty", 3, now, 1, "Team B"))
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:   "no events",
			gameID: 999,
			mockFn: func() {
				mock.ExpectQuery(`SELECT (.+) FROM game_events ge INNER JOIN teams`).
					WithArgs(999).
					WillReturnRows(sqlmock.NewRows([]string{"id", "game_id", "period", "clock", "event_type", "team_id", "created_at", "version", "short_name"}))
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			events, err := model.ListByGame(tt.gameID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListByGame() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(events) != tt.wantCount {
				t.Errorf("ListByGame() returned %d events, want %d", len(events), tt.wantCount)
			}
		})
	}
}

func TestGameEventModel_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	model := GameEventModel{DB: db}

	tests := []struct {
		name    string
		event   *GameEvent
		mockFn  func()
		wantErr error
	}{
		{
			name: "successful update",
			event: &GameEvent{
				ID:        1,
				Period:    2,
				Clock:     "00:10:00",
				EventType: "penalty",
				TeamID:    3,
				Version:   1,
			},
			mockFn: func() {
				mock.ExpectQuery(`UPDATE game_events SET`).
					WithArgs(2, "00:10:00", "penalty", 3, 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow(2))
			},
			wantErr: nil,
		},
		{
			name: "edit conflict",
			event: &GameEvent{
				ID:      1,
				Version: 999,
			},
			mockFn: func() {
				mock.ExpectQuery(`UPDATE game_events SET`).
					WithArgs(0, "", "", 0, 1, 999).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: ErrEditConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := model.Update(tt.event)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGameEventModel_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	model := GameEventModel{DB: db}

	tests := []struct {
		name    string
		id      int
		mockFn  func()
		wantErr error
	}{
		{
			name: "successful delete",
			id:   1,
			mockFn: func() {
				mock.ExpectExec(`DELETE FROM game_events WHERE id = \$1`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: nil,
		},
		{
			name: "not found",
			id:   999,
			mockFn: func() {
				mock.ExpectExec(`DELETE FROM game_events WHERE id = \$1`).
					WithArgs(999).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := model.Delete(tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
