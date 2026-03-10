# Games

## Create a game event

`POST /games/{game_id}/events`

Request:

```json
{
  "period": 2,
  "clock_seconds": 360,
  "event_type": "penalty",
  "team_id": 2,
  "situation": "EV"
}
```

Event types:
- `goal`
- `penalty`
- `shot`
- `save`

Situations:
- `ev` even strength
- `pp` powerplay
- `sh` short handed,
- `en` empty net

Response:

```json
{
  "game_events": {
    "id": 2,
    "event_number": 2,
    "game_id": 1,
    "period": 2,
    "clock_seconds": 12,
    "event_type": "penalty",
    "situation": "EV",
    "team_id": 2
  }
}
```

## View a game event

`GET /events/{id}`

Response:

```json
{
  "game_events": {
    "id": 1,
    "event_number": 1,
    "game_id": 1,
    "period": 2,
    "clock_seconds": 12,
    "event_type": "penalty",
    "situation": "EV",
    "team_id": 2
  }
}
```

