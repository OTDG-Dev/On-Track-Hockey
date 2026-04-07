package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
	playerstats "github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/stats"
)

func TestPlayers(t *testing.T) {
	t.Run("POST creates player", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")
		divisionID := seedDivision(t, url, leagueID, "Atlantic")
		teamID := seedTeam(t, url, divisionID, "New York Rangers", "NYR")

		status := createPlayer(t, url, teamID, "Igor", "Shesterkin", 31)
		if status != http.StatusCreated {
			t.Fatalf("expected %d got %d", http.StatusCreated, status)
		}
	})

	t.Run("GET includes created player", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")
		divisionID := seedDivision(t, url, leagueID, "Atlantic")
		teamID := seedTeam(t, url, divisionID, "New York Rangers", "NYR")
		seedPlayer(t, url, teamID, "Igor", "Shesterkin", 31)

		players := getPlayers(t, url)
		if !containsPlayer(players, "Igor", "Shesterkin") {
			t.Fatalf("expected players list to contain %q %q", "Igor", "Shesterkin")
		}
	})

	t.Run("GET player includes stats", func(t *testing.T) {
		application, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")
		divisionID := seedDivision(t, url, leagueID, "Atlantic")
		teamID := seedTeam(t, url, divisionID, "New York Rangers", "NYR")
		playerID := seedPlayer(t, url, teamID, "Artemi", "Panarin", 10)

		seedSkaterStats(t, application.DB, playerID)

		player := getPlayer(t, url, playerID)

		if player.SkaterStats == nil {
			t.Fatal("expected skater stats to be included")
		}

		if player.SkaterStats.CurrentSeason.Points != 89 {
			t.Fatalf("expected current season points 89 got %d", player.SkaterStats.CurrentSeason.Points)
		}

		if player.SkaterStats.CareerTotals.Assists != 540 {
			t.Fatalf("expected career assists 540 got %d", player.SkaterStats.CareerTotals.Assists)
		}
	})

	t.Run("GET updateStats rebuilds basic skater stats", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")
		divisionID := seedDivision(t, url, leagueID, "Atlantic")
		homeTeamID := seedTeam(t, url, divisionID, "New York Rangers", "NYR")
		awayTeamID := seedTeam(t, url, divisionID, "Boston Bruins", "BOS")

		scorerID := seedPlayerWithPosition(t, url, homeTeamID, "Artemi", "Panarin", 10, data.PositionLW)
		assistID := seedPlayerWithPosition(t, url, homeTeamID, "Adam", "Fox", 23, data.PositionD)

		gameID := seedGame(t, url, homeTeamID, awayTeamID, true)
		goalEventID := seedGameEvent(t, url, gameID, homeTeamID, "goal")
		seedGameEventParticipant(t, url, goalEventID, "scorer", scorerID)
		seedGameEventParticipant(t, url, goalEventID, "assist_primary", assistID)

		penaltyEventID := seedGameEvent(t, url, gameID, homeTeamID, "penalty")
		seedGameEventParticipant(t, url, penaltyEventID, "penalty_taker", scorerID)

		updatedPlayers := updateStats(t, url)
		if updatedPlayers != 2 {
			t.Fatalf("expected 2 updated players got %d", updatedPlayers)
		}

		scorer := getPlayer(t, url, scorerID)
		if scorer.SkaterStats == nil {
			t.Fatal("expected scorer stats to be included")
		}

		if scorer.SkaterStats.CurrentSeason.GamesPlayed != 1 {
			t.Fatalf("expected scorer games played 1 got %d", scorer.SkaterStats.CurrentSeason.GamesPlayed)
		}

		if scorer.SkaterStats.CurrentSeason.Goals != 1 {
			t.Fatalf("expected scorer goals 1 got %d", scorer.SkaterStats.CurrentSeason.Goals)
		}

		if scorer.SkaterStats.CurrentSeason.Points != 1 {
			t.Fatalf("expected scorer points 1 got %d", scorer.SkaterStats.CurrentSeason.Points)
		}

		if scorer.SkaterStats.CurrentSeason.PIM != 1 {
			t.Fatalf("expected scorer pim 1 got %d", scorer.SkaterStats.CurrentSeason.PIM)
		}

		assist := getPlayer(t, url, assistID)
		if assist.SkaterStats == nil {
			t.Fatal("expected assist stats to be included")
		}

		if assist.SkaterStats.CurrentSeason.Assists != 1 {
			t.Fatalf("expected assist assists 1 got %d", assist.SkaterStats.CurrentSeason.Assists)
		}

		if assist.SkaterStats.CurrentSeason.Points != 1 {
			t.Fatalf("expected assist points 1 got %d", assist.SkaterStats.CurrentSeason.Points)
		}
	})

	t.Run("PATCH edit a player", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")
		divisionID := seedDivision(t, url, leagueID, "Atlantic")
		teamID := seedTeam(t, url, divisionID, "New York Rangers", "NYR")
		playerID := seedPlayer(t, url, teamID, "Igor", "Shesterkin", 31)

		status := patchPlayer(t, url, playerID, "Connor")
		if status != http.StatusOK {
			t.Fatalf("expected %d got %d", http.StatusOK, status)
		}

		players := getPlayers(t, url)
		if !containsPlayer(players, "Connor", "Shesterkin") {
			t.Fatalf("expected players list to contain updated player")
		}
	})

	t.Run("DELETE delete a player", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")
		divisionID := seedDivision(t, url, leagueID, "Atlantic")
		teamID := seedTeam(t, url, divisionID, "New York Rangers", "NYR")
		playerID := seedPlayer(t, url, teamID, "Igor", "Shesterkin", 31)

		status := deletePlayer(t, url, playerID)
		if status != http.StatusOK {
			t.Fatalf("expected %d got %d", http.StatusOK, status)
		}

		players := getPlayers(t, url)
		if containsPlayer(players, "Igor", "Shesterkin") {
			t.Fatalf("expected players list to not contain deleted player")
		}
	})
}

func createPlayer(t *testing.T, url string, teamID int, firstName, lastName string, sweaterNumber uint8) int {
	t.Helper()

	var body bytes.Buffer

	birthDate := data.BirthDate(time.Date(1995, time.December, 30, 0, 0, 0, 0, time.UTC))
	player := struct {
		IsActive      bool               `json:"is_active"`
		CurrentTeamID int                `json:"current_team_id"`
		FirstName     string             `json:"first_name"`
		LastName      string             `json:"last_name"`
		SweaterNumber uint8              `json:"sweater_number"`
		Position      data.Position      `json:"position"`
		BirthDate     data.BirthDate     `json:"birth_date"`
		BirthCountry  string             `json:"birth_country"`
		ShootsCatches data.ShootsCatches `json:"shoots_catches"`
	}{
		IsActive:      true,
		CurrentTeamID: teamID,
		FirstName:     firstName,
		LastName:      lastName,
		SweaterNumber: sweaterNumber,
		Position:      data.PositionG,
		BirthDate:     birthDate,
		BirthCountry:  "RUS",
		ShootsCatches: data.ShootsCatchesL,
	}

	if err := json.NewEncoder(&body).Encode(player); err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(url+"/v1/players", "application/json", &body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func seedPlayer(t *testing.T, url string, teamID int, firstName, lastName string, sweaterNumber uint8) int {
	t.Helper()

	var body bytes.Buffer

	birthDate := data.BirthDate(time.Date(1995, time.December, 30, 0, 0, 0, 0, time.UTC))
	player := struct {
		IsActive      bool               `json:"is_active"`
		CurrentTeamID int                `json:"current_team_id"`
		FirstName     string             `json:"first_name"`
		LastName      string             `json:"last_name"`
		SweaterNumber uint8              `json:"sweater_number"`
		Position      data.Position      `json:"position"`
		BirthDate     data.BirthDate     `json:"birth_date"`
		BirthCountry  string             `json:"birth_country"`
		ShootsCatches data.ShootsCatches `json:"shoots_catches"`
	}{
		IsActive:      true,
		CurrentTeamID: teamID,
		FirstName:     firstName,
		LastName:      lastName,
		SweaterNumber: sweaterNumber,
		Position:      data.PositionG,
		BirthDate:     birthDate,
		BirthCountry:  "RUS",
		ShootsCatches: data.ShootsCatchesL,
	}

	if err := json.NewEncoder(&body).Encode(player); err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(url+"/v1/players", "application/json", &body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("seedPlayer: expected %d got %d", http.StatusCreated, resp.StatusCode)
	}

	var out struct {
		Player struct {
			ID int `json:"id"`
		} `json:"player"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}

	return out.Player.ID
}

func seedPlayerWithPosition(t *testing.T, url string, teamID int, firstName, lastName string, sweaterNumber uint8, position data.Position) int {
	t.Helper()

	var body bytes.Buffer

	birthDate := data.BirthDate(time.Date(1995, time.December, 30, 0, 0, 0, 0, time.UTC))
	player := struct {
		IsActive      bool               `json:"is_active"`
		CurrentTeamID int                `json:"current_team_id"`
		FirstName     string             `json:"first_name"`
		LastName      string             `json:"last_name"`
		SweaterNumber uint8              `json:"sweater_number"`
		Position      data.Position      `json:"position"`
		BirthDate     data.BirthDate     `json:"birth_date"`
		BirthCountry  string             `json:"birth_country"`
		ShootsCatches data.ShootsCatches `json:"shoots_catches"`
	}{
		IsActive:      true,
		CurrentTeamID: teamID,
		FirstName:     firstName,
		LastName:      lastName,
		SweaterNumber: sweaterNumber,
		Position:      position,
		BirthDate:     birthDate,
		BirthCountry:  "CAN",
		ShootsCatches: data.ShootsCatchesL,
	}

	if err := json.NewEncoder(&body).Encode(player); err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(url+"/v1/players", "application/json", &body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("seedPlayerWithPosition: expected %d got %d", http.StatusCreated, resp.StatusCode)
	}

	var out struct {
		Player struct {
			ID int `json:"id"`
		} `json:"player"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}

	return out.Player.ID
}

func getPlayers(t *testing.T, url string) []*data.PlayerWithTeam {
	t.Helper()

	r, err := http.Get(url + "/v1/players")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, r.StatusCode)
	}

	var resp struct {
		Players []*data.PlayerWithTeam `json:"players"`
	}

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	return resp.Players
}

func getPlayer(t *testing.T, url string, id int) *data.PlayerWithTeam {
	t.Helper()

	r, err := http.Get(url + fmt.Sprintf("/v1/players/%d", id))
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, r.StatusCode)
	}

	var resp struct {
		Player *data.PlayerWithTeam `json:"player"`
	}

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	return resp.Player
}

func seedSkaterStats(t *testing.T, db *sql.DB, playerID int) {
	t.Helper()

	skaterStats := playerstats.SkaterStatSet{
		CurrentSeason: playerstats.SkaterStats{
			GamesPlayed: 67,
			Goals:       31,
			Assists:     58,
			Points:      89,
			PIM:         22,
		},
		CareerTotals: playerstats.SkaterStats{
			GamesPlayed: 700,
			Goals:       260,
			Assists:     540,
			Points:      800,
			PIM:         180,
		},
	}

	payload, err := json.Marshal(skaterStats)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, `
		INSERT INTO player_stats (player_id, skater_stats)
		VALUES ($1, $2)
	`, playerID, payload)
	if err != nil {
		t.Fatal(err)
	}
}

func seedGameEvent(t *testing.T, url string, gameID, teamID int, eventType string) int {
	t.Helper()

	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(map[string]any{
		"period":        1,
		"clock_seconds": 120,
		"event_type":    eventType,
		"team_id":       teamID,
		"situation":     "EV",
	}); err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(url+fmt.Sprintf("/v1/games/%d/events", gameID), "application/json", &body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("seedGameEvent: expected %d got %d", http.StatusCreated, resp.StatusCode)
	}

	var out struct {
		GameEvents struct {
			ID int `json:"id"`
		} `json:"game_events"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}

	return out.GameEvents.ID
}

func seedGameEventParticipant(t *testing.T, url string, eventID int, role string, playerID int) {
	t.Helper()

	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(map[string]any{
		"role":      role,
		"player_id": playerID,
	}); err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(url+fmt.Sprintf("/v1/events/%d/participants", eventID), "application/json", &body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("seedGameEventParticipant: expected %d got %d", http.StatusCreated, resp.StatusCode)
	}
}

func updateStats(t *testing.T, url string) int {
	t.Helper()

	resp, err := http.Get(url + "/v0/updateStats")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, resp.StatusCode)
	}

	var out struct {
		UpdatedPlayers int `json:"updated_players"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}

	return out.UpdatedPlayers
}

func patchPlayer(t *testing.T, url string, id int, firstName string) int {
	t.Helper()

	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(map[string]any{"first_name": firstName}); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPatch, url+fmt.Sprintf("/v1/players/%d", id), &body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func deletePlayer(t *testing.T, url string, id int) int {
	t.Helper()

	req, err := http.NewRequest(http.MethodDelete, url+fmt.Sprintf("/v1/players/%d", id), nil)
	if err != nil {
		t.Fatal(err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func containsPlayer(players []*data.PlayerWithTeam, firstName, lastName string) bool {
	for _, player := range players {
		if player.FirstName == firstName && player.LastName == lastName {
			return true
		}
	}

	return false
}
