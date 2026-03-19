# Game Event Participants

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
- penalty_take

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
