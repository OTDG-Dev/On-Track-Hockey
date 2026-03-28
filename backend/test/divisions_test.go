package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
)

func TestDivisions(t *testing.T) {
	t.Run("POST creates division", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")

		status := createDivision(t, url, leagueID, "Atlantic")
		if status != http.StatusCreated {
			t.Fatalf("expected %d got %d", http.StatusCreated, status)
		}
	})

	t.Run("GET includes created division", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")
		createDivision(t, url, leagueID, "Atlantic")

		divisions := getDivisions(t, url)
		if !containsDivision(divisions, "Atlantic") {
			t.Fatalf("expected divisions list to contain %q", "Atlantic")
		}
	})

	t.Run("PATCH edit a division", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")
		divisionID := seedDivision(t, url, leagueID, "Atlantic")

		newDivisionName := "Metropolitan"

		status := patchDivision(t, url, divisionID, newDivisionName)
		if status != http.StatusOK {
			t.Fatalf("expected %d got %d", http.StatusOK, status)
		}

		divisions := getDivisions(t, url)
		if !containsDivision(divisions, newDivisionName) {
			t.Fatalf("expected divisions list to contain %q", newDivisionName)
		}
	})

	t.Run("DELETE delete a division", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")
		divisionID := seedDivision(t, url, leagueID, "Atlantic")

		status := deleteDivision(t, url, divisionID)
		if status != http.StatusOK {
			t.Fatalf("expected %d got %d", http.StatusOK, status)
		}

		divisions := getDivisions(t, url)
		if containsDivision(divisions, "Atlantic") {
			t.Fatalf("expected divisions list to not contain %q", "Atlantic")
		}
	})
}

func createDivision(t *testing.T, url string, leagueID int, name string) int {
	t.Helper()

	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(map[string]any{
		"name":      name,
		"league_id": leagueID,
	}); err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(url+"/v1/divisions", "application/json", &body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func getDivisions(t *testing.T, url string) []*data.Division {
	t.Helper()

	r, err := http.Get(url + "/v1/divisions")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, r.StatusCode)
	}

	var resp struct {
		Divisions []*data.Division `json:"divisions"`
	}

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	return resp.Divisions
}

func patchDivision(t *testing.T, url string, id int, name string) int {
	t.Helper()

	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(map[string]any{"name": name}); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPatch, url+fmt.Sprintf("/v1/divisions/%d", id), &body)
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

func deleteDivision(t *testing.T, url string, id int) int {
	t.Helper()

	req, err := http.NewRequest(http.MethodDelete, url+fmt.Sprintf("/v1/divisions/%d", id), nil)
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

func containsDivision(divisions []*data.Division, name string) bool {
	for _, division := range divisions {
		if division.Name == name {
			return true
		}
	}

	return false
}
