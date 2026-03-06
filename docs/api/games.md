# Games

## Create a game

`POST /games`

Request:

```json
{
  "home_team_id": 1,
  "away_team_id": 2,
  "start_time": "2026-03-10T19:30:00-05:00"
}
```

Response:

```json
{
  "game": {
    "HomeTeamID": 1,
    "AwayTeamID": 2,
    "StartTime": "2026-03-10T19:30:00-05:00",
    "Version": 1
  }
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
    "home_team_id": "1",
    "away_team_id": "2",
    "start_time": "2026-03-11T00:30:00Z"
  }
}
```

