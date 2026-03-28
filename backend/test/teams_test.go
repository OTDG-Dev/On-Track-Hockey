package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
)

func TestTeams(t *testing.T) {
	t.Run("POST creates team", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")
		divisionID := seedDivision(t, url, leagueID, "Atlantic")

		status := createTeam(t, url, divisionID, "New York Rangers", "NYR")
		if status != http.StatusCreated {
			t.Fatalf("expected %d got %d", http.StatusCreated, status)
		}
	})

	t.Run("GET includes created team", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")
		divisionID := seedDivision(t, url, leagueID, "Atlantic")
		createTeam(t, url, divisionID, "New York Rangers", "NYR")

		teams := getTeams(t, url)
		if !containsTeam(teams, "New York Rangers") {
			t.Fatalf("expected teams list to contain %q", "New York Rangers")
		}
	})

	t.Run("PATCH edit a team", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")
		divisionID := seedDivision(t, url, leagueID, "Atlantic")
		teamID := seedTeam(t, url, divisionID, "New York Rangers", "NYR")

		newTeamName := "New York Islanders"

		status := patchTeam(t, url, teamID, newTeamName)
		if status != http.StatusOK {
			t.Fatalf("expected %d got %d", http.StatusOK, status)
		}

		teams := getTeams(t, url)
		if !containsTeam(teams, newTeamName) {
			t.Fatalf("expected teams list to contain %q", newTeamName)
		}
	})

	t.Run("DELETE delete a team", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")
		divisionID := seedDivision(t, url, leagueID, "Atlantic")
		teamID := seedTeam(t, url, divisionID, "New York Rangers", "NYR")

		status := deleteTeam(t, url, teamID)
		if status != http.StatusOK {
			t.Fatalf("expected %d got %d", http.StatusOK, status)
		}

		teams := getTeams(t, url)
		if containsTeam(teams, "New York Rangers") {
			t.Fatalf("expected teams list to not contain %q", "New York Rangers")
		}
	})
}

func createTeam(t *testing.T, url string, divisionID int, fullName, shortName string) int {
	t.Helper()

	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(map[string]any{
		"full_name":   fullName,
		"short_name":  shortName,
		"division_id": divisionID,
		"is_active":   true,
	}); err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(url+"/v1/teams", "application/json", &body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func getTeams(t *testing.T, url string) []*data.Team {
	t.Helper()

	r, err := http.Get(url + "/v1/teams")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, r.StatusCode)
	}

	var resp struct {
		Teams []*data.Team `json:"teams"`
	}

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	return resp.Teams
}

func patchTeam(t *testing.T, url string, id int, fullName string) int {
	t.Helper()

	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(map[string]any{"full_name": fullName}); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPatch, url+fmt.Sprintf("/v1/teams/%d", id), &body)
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

func deleteTeam(t *testing.T, url string, id int) int {
	t.Helper()

	req, err := http.NewRequest(http.MethodDelete, url+fmt.Sprintf("/v1/teams/%d", id), nil)
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

func containsTeam(teams []*data.Team, fullName string) bool {
	for _, team := range teams {
		if team.FullName == fullName {
			return true
		}
	}

	return false
}
