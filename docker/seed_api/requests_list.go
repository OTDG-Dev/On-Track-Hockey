package main

var requestList = []Req{
	{
		Endpoint: "leagues",
		Method:   "POST",
		ExpCode:  201,
		Payload: map[string]any{
			"name": "league",
		},
	},
	{
		Endpoint: "divisions",
		Method:   "POST",
		ExpCode:  201,
		Payload: map[string]any{
			"league_id": 1,
			"name":      "Team",
		},
	},
	{
		Endpoint: "teams",
		Method:   "POST",
		ExpCode:  201,
		Payload: map[string]any{
			"full_name":   "New York Rangers",
			"short_name":  "NYR",
			"division_id": 1,
			"is_active":   true,
		},
	},
	{
		Endpoint: "teams",
		Method:   "POST",
		ExpCode:  201,
		Payload: map[string]any{
			"full_name":   "San Jose Sharks",
			"short_name":  "SJS",
			"division_id": 1,
			"is_active":   true,
		},
	},
	{
		Endpoint: "players",
		Method:   "POST",
		ExpCode:  201,
		Payload: map[string]any{
			"first_name":      "Igor",
			"last_name":       "Shesterkin",
			"sweater_number":  31,
			"position":        "G",
			"birth_date":      "1995-12-30",
			"birth_country":   "RUS",
			"shoots_catches":  "L",
			"current_team_id": 1,
			"is_active":       true,
		},
	},
	{
		Endpoint: "players",
		Method:   "POST",
		ExpCode:  201,
		Payload: map[string]any{
			"first_name":      "Macklin",
			"last_name":       "Celebrini",
			"sweater_number":  71,
			"position":        "C",
			"birth_date":      "2006-05-12",
			"birth_country":   "CAN",
			"shoots_catches":  "L",
			"current_team_id": 2,
			"is_active":       true,
		},
	},
}
