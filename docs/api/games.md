# Games

## Create a game

`POST /games`

Request:

```json
{
  "home_team_id": 1,
  "away_team_id": 2,
  "start_time": "2026-03-10T19:30:00-05:00",
  "is_finished": false
}
```

Response:

```json
{
  "game": {
    "id": 2,
    "home_team_id": 1,
    "away_team_id": 2,
    "start_time": "2026-03-10T19:30:00-05:00",
    "is_finished": false
  }
}
```

## Get all games

`GET /games`

Response:

```json
{
  "games": [
    {
      "id": 1,
      "home_team": "NYR",
      "away_team": "SJS",
      "home_team_id": 1,
      "away_team_id": 2,
      "is_finished": false,
      "start_time": "2026-03-11T00:30:00Z"
    },
    {
      "id": 2,
      "home_team": "NYR",
      "away_team": "SJS",
      "home_team_id": 1,
      "away_team_id": 2,
      "is_finished": true,
      "start_time": "2026-03-11T00:30:00Z"
    }
  ]
}
```

## View a game

`GET /games/{id}`

Response:

```json
{
  "game": {
    "home_team": "NYR",
    "away_team": "SJS",
    "home_team_id": 1,
    "away_team_id": 2,
    "is_finished": false,
    "start_time": "2026-03-11T00:30:00Z",
    "game_events": [
      {
        "id": 1,
        "event_number": 1,
        "period": 2,
        "clock_seconds": 12,
        "event_type": "penalty",
        "situation": "EV",
        "team_id": 2
      },
      {
        "id": 2,
        "event_number": 2,
        "period": 2,
        "clock_seconds": 12,
        "event_type": "penalty",
        "situation": "EV",
        "team_id": 2
      }
    ]
  }
}
```

