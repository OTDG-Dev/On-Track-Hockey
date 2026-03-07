package main

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/lib/pq"
)

func TestFormatForeignKeyError(t *testing.T) {
	app := &application{}

	tests := []struct {
		name      string
		pgErr     *pq.Error
		wantPart  string
	}{
		{
			name: "game_id foreign key",
			pgErr: &pq.Error{
				Code:      "23503",
				Constraint: "game_events_game_id_fkey",
				Detail:     "Key (game_id)=(999) is not present in table \"games\"",
			},
			wantPart: "the referenced game does not exist",
		},
		{
			name: "team_id foreign key",
			pgErr: &pq.Error{
				Code:      "23503",
				Constraint: "game_events_team_id_fkey",
			},
			wantPart: "the referenced team does not exist",
		},
		{
			name: "player_id foreign key",
			pgErr: &pq.Error{
				Code:      "23503",
				Constraint: "event_participants_player_id_fkey",
			},
			wantPart: "the referenced player does not exist",
		},
		{
			name: "game_event_id foreign key",
			pgErr: &pq.Error{
				Code:      "23503",
				Constraint: "event_participants_game_event_id_fkey",
			},
			wantPart: "the referenced game event does not exist",
		},
		{
			name: "home_team_id foreign key",
			pgErr: &pq.Error{
				Code:      "23503",
				Constraint: "games_home_team_id_fkey",
			},
			wantPart: "the referenced team does not exist",
		},
		{
			name: "away_team_id foreign key",
			pgErr: &pq.Error{
				Code:      "23503",
				Constraint: "games_away_team_id_fkey",
			},
			wantPart: "the referenced team does not exist",
		},
		{
			name: "division_id foreign key",
			pgErr: &pq.Error{
				Code:      "23503",
				Constraint: "teams_division_id_fkey",
			},
			wantPart: "the referenced division does not exist",
		},
		{
			name: "league_id foreign key",
			pgErr: &pq.Error{
				Code:      "23503",
				Constraint: "divisions_league_id_fkey",
			},
			wantPart: "the referenced league does not exist",
		},
		{
			name: "unknown foreign key",
			pgErr: &pq.Error{
				Code:      "23503",
				Constraint: "some_other_constraint",
			},
			wantPart: "a related record does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := app.formatForeignKeyError(tt.pgErr)
			if got == "" {
				t.Error("formatForeignKeyError() returned empty string")
			}
			if !contains(got, tt.wantPart) {
				t.Errorf("formatForeignKeyError() = %v, want to contain %v", got, tt.wantPart)
			}
		})
	}
}

func TestHandleDatabaseError_ForeignKeyViolation(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/v1/game-events", nil)

	pgErr := &pq.Error{
		Code:       "23503",
		Constraint: "game_events_game_id_fkey",
	}

	app.handleDatabaseError(w, r, pgErr)

	resp := w.Result()
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("handleDatabaseError() status = %v, want %v", resp.StatusCode, http.StatusConflict)
	}
}

func TestHandleDatabaseError_UniqueViolation(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/v1/players", nil)

	pgErr := &pq.Error{
		Code:       "23505",
		Constraint: "players_email_key",
	}

	app.handleDatabaseError(w, r, pgErr)

	resp := w.Result()
	if resp.StatusCode != http.StatusConflict {
		t.Errorf("handleDatabaseError() status = %v, want %v", resp.StatusCode, http.StatusConflict)
	}
}

func TestHandleDatabaseError_NotNullViolation(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/v1/players", nil)

	pgErr := &pq.Error{
		Code:   "23502",
		Column: "first_name",
	}

	app.handleDatabaseError(w, r, pgErr)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("handleDatabaseError() status = %v, want %v", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestHandleDatabaseError_NonPostgresError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/v1/game-events", nil)

	// Non-postgres error should return 500
	app.handleDatabaseError(w, r, errors.New("some database error"))

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("handleDatabaseError() status = %v, want %v", resp.StatusCode, http.StatusInternalServerError)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
