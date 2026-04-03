package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
)

func TestGames(t *testing.T) {
	t.Run("POST creates game", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")
		divisionID := seedDivision(t, url, leagueID, "Atlantic")
		homeTeamID := seedTeam(t, url, divisionID, "New York Rangers", "NYR")
		awayTeamID := seedTeam(t, url, divisionID, "Boston Bruins", "BOS")

		status := createGame(t, url, homeTeamID, awayTeamID, false)
		if status != http.StatusCreated {
			t.Fatalf("expected %d got %d", http.StatusCreated, status)
		}
	})

	t.Run("GET includes created game", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")
		divisionID := seedDivision(t, url, leagueID, "Atlantic")
		homeTeamID := seedTeam(t, url, divisionID, "New York Rangers", "NYR")
		awayTeamID := seedTeam(t, url, divisionID, "Boston Bruins", "BOS")

		seedGame(t, url, homeTeamID, awayTeamID, true)

		games := getGames(t, url)
		if len(games) != 1 {
			t.Fatalf("expected 1 game got %d", len(games))
		}
		if !games[0].IsFinished {
			t.Fatal("expected game list to include finished game")
		}
	})

	t.Run("PATCH edit a game", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")
		divisionID := seedDivision(t, url, leagueID, "Atlantic")
		homeTeamID := seedTeam(t, url, divisionID, "New York Rangers", "NYR")
		awayTeamID := seedTeam(t, url, divisionID, "Boston Bruins", "BOS")
		newAwayTeamID := seedTeam(t, url, divisionID, "Montreal Canadiens", "MTL")
		gameID := seedGame(t, url, homeTeamID, awayTeamID, false)

		status := patchGame(t, url, gameID, homeTeamID, newAwayTeamID, true)
		if status != http.StatusOK {
			t.Fatalf("expected %d got %d", http.StatusOK, status)
		}

		game := getGame(t, url, gameID)
		if game.AwayTeamID != newAwayTeamID {
			t.Fatalf("expected away_team_id %d got %d", newAwayTeamID, game.AwayTeamID)
		}
		if !game.IsFinished {
			t.Fatal("expected patched game to be finished")
		}
	})
}

func createGame(t *testing.T, url string, homeTeamID, awayTeamID int, isFinished bool) int {
	t.Helper()

	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(map[string]any{
		"home_team_id": homeTeamID,
		"away_team_id": awayTeamID,
		"start_time":   time.Date(2026, time.January, 2, 19, 0, 0, 0, time.UTC),
		"is_finished":  isFinished,
	}); err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(url+"/v1/games", "application/json", &body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func seedGame(t *testing.T, url string, homeTeamID, awayTeamID int, isFinished bool) int {
	t.Helper()

	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(map[string]any{
		"home_team_id": homeTeamID,
		"away_team_id": awayTeamID,
		"start_time":   time.Date(2026, time.January, 2, 19, 0, 0, 0, time.UTC),
		"is_finished":  isFinished,
	}); err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(url+"/v1/games", "application/json", &body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("seedGame: expected %d got %d", http.StatusCreated, resp.StatusCode)
	}

	var out struct {
		Game struct {
			ID int `json:"id"`
		} `json:"game"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}

	return out.Game.ID
}

func getGames(t *testing.T, url string) []*data.GameListView {
	t.Helper()

	r, err := http.Get(url + "/v1/games")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, r.StatusCode)
	}

	var resp struct {
		Games []*data.GameListView `json:"games"`
	}

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	return resp.Games
}

func getGame(t *testing.T, url string, id int) *data.GameView {
	t.Helper()

	r, err := http.Get(url + fmt.Sprintf("/v1/games/%d", id))
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, r.StatusCode)
	}

	var resp struct {
		Game *data.GameView `json:"game"`
	}

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	return resp.Game
}

func patchGame(t *testing.T, url string, id, homeTeamID, awayTeamID int, isFinished bool) int {
	t.Helper()

	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(map[string]any{
		"home_team_id": homeTeamID,
		"away_team_id": awayTeamID,
		"is_finished":  isFinished,
	}); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPatch, url+fmt.Sprintf("/v1/games/%d", id), &body)
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
