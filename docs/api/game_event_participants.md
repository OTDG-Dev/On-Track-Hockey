# Game Event Participants

## List participants for a game event

`GET /events/{event_id}/participants`

Response:

```json
{
    "game_event_participants": [
        {
            "id": 1,
            "role": "scorer",
            "event_id": 1,
            "player_id": 1
        },
        {
            "id": 2,
            "role": "assist_primary",
            "event_id": 1,
            "player_id": 2
        }
    ]
}
```

## Create a game event participant

`POST /events/{event_id}/participants`

Request:

```json
{
    "role": "scorer",
    "player_id": 1,
}
```

Roles:


- scorer
- assist_primary
- assist_secondary
- penalty_taker

Response:

```json
{
    "game_event_participant": {
        "role": "assist_secondary",
        "event_id": 1,
        "player_id": 3
    }
}
```
