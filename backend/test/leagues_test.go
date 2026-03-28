package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
)

func TestLeagues(t *testing.T) {
	t.Run("POST creates league", func(t *testing.T) {
		_, url := setup(t)

		status := createLeague(t, url, "NHL")
		if status != http.StatusCreated {
			t.Fatalf("expected %d got %d", http.StatusCreated, status)
		}
	})

	t.Run("GET includes created league", func(t *testing.T) {
		_, url := setup(t)

		seedLeague(t, url, "NHL")

		leagues := getLeagues(t, url)
		if !containsLeague(leagues, "NHL") {
			t.Fatalf("expected leagues list to contain %q", "NHL")
		}
	})

	t.Run("PATCH edit a league", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")

		newLeagueName := "AHL"

		status := patchLeague(t, url, leagueID, newLeagueName)
		if status != http.StatusOK {
			t.Fatalf("expected %d got %d", http.StatusOK, status)
		}

		leagues := getLeagues(t, url)
		if !containsLeague(leagues, newLeagueName) {
			t.Fatalf("expected leagues list to contain %q", newLeagueName)
		}
	})

	t.Run("DELETE delete a league", func(t *testing.T) {
		_, url := setup(t)

		leagueID := seedLeague(t, url, "NHL")

		status := deleteLeague(t, url, leagueID)
		if status != http.StatusOK {
			t.Fatalf("expected %d got %d", http.StatusOK, status)
		}

		leagues := getLeagues(t, url)
		if containsLeague(leagues, "NHL") {
			t.Fatalf("expected leagues list to not contain %q", "NHL")
		}
	})
}

func createLeague(t *testing.T, url, name string) int {
	t.Helper()

	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(map[string]any{"name": name}); err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(url+"/v1/leagues", "application/json", &body)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func getLeagues(t *testing.T, url string) []*data.LeagueWithMetadata {
	t.Helper()

	r, err := http.Get(url + "/v1/leagues")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, r.StatusCode)
	}

	var resp struct {
		Leagues []*data.LeagueWithMetadata `json:"leagues"`
	}

	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	return resp.Leagues
}

func patchLeague(t *testing.T, url string, id int, name string) int {
	t.Helper()

	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(map[string]any{"name": name}); err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPatch, url+fmt.Sprintf("/v1/leagues/%d", id), &body)
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

func deleteLeague(t *testing.T, url string, id int) int {
	t.Helper()

	req, err := http.NewRequest(http.MethodDelete, url+fmt.Sprintf("/v1/leagues/%d", id), nil)
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

func containsLeague(leagues []*data.LeagueWithMetadata, name string) bool {
	for _, league := range leagues {
		if league.Name == name {
			return true
		}
	}

	return false
}
