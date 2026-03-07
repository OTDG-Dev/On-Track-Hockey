package data

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestEventParticipantModel_Insert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	model := EventParticipantModel{DB: db}

	tests := []struct {
		name        string
		participant *EventParticipant
		mockFn      func()
		wantErr     bool
	}{
		{
			name: "successful insert",
			participant: &EventParticipant{
				GameEventID: 101,
				PlayerID:    63,
				Role:        "scorer",
			},
			mockFn: func() {
				mock.ExpectQuery(`INSERT INTO event_participants`).
					WithArgs(101, 63, "scorer").
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(1, time.Now()))
			},
			wantErr: false,
		},
		{
			name: "foreign key violation",
			participant: &EventParticipant{
				GameEventID: 999,
				PlayerID:    63,
				Role:        "scorer",
			},
			mockFn: func() {
				mock.ExpectQuery(`INSERT INTO event_participants`).
					WithArgs(999, 63, "scorer").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := model.Insert(tt.participant)
			if (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventParticipantModel_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	model := EventParticipantModel{DB: db}
	now := time.Now()

	tests := []struct {
		name          string
		id            int
		mockFn        func()
		wantParticipant *EventParticipant
		wantErr       error
	}{
		{
			name: "existing participant",
			id:   1,
			mockFn: func() {
				mock.ExpectQuery(`SELECT id, game_event_id, player_id, role, created_at FROM event_participants WHERE id = \$1`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "game_event_id", "player_id", "role", "created_at"}).
						AddRow(1, 101, 63, "scorer", now))
			},
			wantParticipant: &EventParticipant{
				ID:          1,
				GameEventID: 101,
				PlayerID:    63,
				Role:        "scorer",
				CreatedAt:   now,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			id:   999,
			mockFn: func() {
				mock.ExpectQuery(`SELECT id, game_event_id, player_id, role, created_at FROM event_participants WHERE id = \$1`).
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			wantParticipant: nil,
			wantErr:         ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			participant, err := model.Get(tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantParticipant != nil && participant != nil {
				if participant.ID != tt.wantParticipant.ID {
					t.Errorf("Get() participant.ID = %v, want %v", participant.ID, tt.wantParticipant.ID)
				}
			}
		})
	}
}

func TestEventParticipantModel_ListByGameEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	model := EventParticipantModel{DB: db}
	now := time.Now()

	tests := []struct {
		name      string
		gameEventID int
		mockFn    func()
		wantCount int
		wantErr   bool
	}{
		{
			name:        "participants exist",
			gameEventID: 101,
			mockFn: func() {
				mock.ExpectQuery(`SELECT (.+) FROM event_participants ep INNER JOIN players`).
					WithArgs(101).
					WillReturnRows(sqlmock.NewRows([]string{"id", "game_event_id", "player_id", "role", "created_at", "first_name", "last_name", "number"}).
						AddRow(1, 101, 63, "scorer", now, "John", "Doe", sql.NullInt32{Int32: 23, Valid: true}).
						AddRow(2, 101, 81, "assist", now, "Jane", "Smith", sql.NullInt32{Int32: 14, Valid: true}).
						AddRow(3, 101, 2, "assist", now, "Bob", "Jones", sql.NullInt32{Valid: false}))
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:        "no participants",
			gameEventID: 999,
			mockFn: func() {
				mock.ExpectQuery(`SELECT (.+) FROM event_participants ep INNER JOIN players`).
					WithArgs(999).
					WillReturnRows(sqlmock.NewRows([]string{"id", "game_event_id", "player_id", "role", "created_at", "first_name", "last_name", "number"}))
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			participants, err := model.ListByGameEvent(tt.gameEventID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListByGameEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(participants) != tt.wantCount {
				t.Errorf("ListByGameEvent() returned %d participants, want %d", len(participants), tt.wantCount)
			}
		})
	}
}

func TestEventParticipantModel_ListByPlayer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	model := EventParticipantModel{DB: db}
	now := time.Now()

	tests := []struct {
		name      string
		playerID  int
		mockFn    func()
		wantCount int
		wantErr   bool
	}{
		{
			name:     "participations exist",
			playerID: 63,
			mockFn: func() {
				mock.ExpectQuery(`SELECT (.+) FROM event_participants ep INNER JOIN players`).
					WithArgs(63).
					WillReturnRows(sqlmock.NewRows([]string{"id", "game_event_id", "player_id", "role", "created_at", "first_name", "last_name", "number"}).
						AddRow(1, 101, 63, "scorer", now, "John", "Doe", sql.NullInt32{Int32: 23, Valid: true}).
						AddRow(4, 102, 63, "scorer", now, "John", "Doe", sql.NullInt32{Int32: 23, Valid: true}))
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:     "no participations",
			playerID: 999,
			mockFn: func() {
				mock.ExpectQuery(`SELECT (.+) FROM event_participants ep INNER JOIN players`).
					WithArgs(999).
					WillReturnRows(sqlmock.NewRows([]string{"id", "game_event_id", "player_id", "role", "created_at", "first_name", "last_name", "number"}))
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			participants, err := model.ListByPlayer(tt.playerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListByPlayer() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(participants) != tt.wantCount {
				t.Errorf("ListByPlayer() returned %d participants, want %d", len(participants), tt.wantCount)
			}
		})
	}
}

func TestEventParticipantModel_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	model := EventParticipantModel{DB: db}

	tests := []struct {
		name        string
		participant *EventParticipant
		mockFn      func()
		wantErr     error
	}{
		{
			name: "successful update",
			participant: &EventParticipant{
				ID:       1,
				PlayerID: 65,
				Role:     "goalie",
			},
			mockFn: func() {
				mock.ExpectExec(`UPDATE event_participants SET`).
					WithArgs(65, "goalie", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: nil,
		},
		{
			name: "not found",
			participant: &EventParticipant{
				ID:       999,
				PlayerID: 65,
				Role:     "goalie",
			},
			mockFn: func() {
				mock.ExpectExec(`UPDATE event_participants SET`).
					WithArgs(65, "goalie", 999).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := model.Update(tt.participant)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventParticipantModel_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	model := EventParticipantModel{DB: db}

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
				mock.ExpectExec(`DELETE FROM event_participants WHERE id = \$1`).
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: nil,
		},
		{
			name: "not found",
			id:   999,
			mockFn: func() {
				mock.ExpectExec(`DELETE FROM event_participants WHERE id = \$1`).
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

func TestEventParticipantModel_DeleteByGameEvent(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	model := EventParticipantModel{DB: db}

	tests := []struct {
		name        string
		gameEventID int
		mockFn      func()
		wantErr     bool
	}{
		{
			name:        "successful delete",
			gameEventID: 101,
			mockFn: func() {
				mock.ExpectExec(`DELETE FROM event_participants WHERE game_event_id = \$1`).
					WithArgs(101).
					WillReturnResult(sqlmock.NewResult(0, 3))
			},
			wantErr: false,
		},
		{
			name:        "no participants to delete",
			gameEventID: 999,
			mockFn: func() {
				mock.ExpectExec(`DELETE FROM event_participants WHERE game_event_id = \$1`).
					WithArgs(999).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn()
			err := model.DeleteByGameEvent(tt.gameEventID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteByGameEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
