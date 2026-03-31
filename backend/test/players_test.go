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
