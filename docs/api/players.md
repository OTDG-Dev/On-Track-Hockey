# Players

## Create a player

`POST /players`

Request:

```json
{
  "first_name": "Igor",
  "last_name": "Shesterkin",
  "sweater_number": 31,
  "position": "G",
  "birth_date": "1995-12-30",
  "birth_country": "RUS",
  "shoots_catches": "L",
  "current_team_id": 1,
  "is_active": true
}
```


Response:

```json
{
  "player": {
    "id": 35,
    "is_active": true,
    "current_team_id": 1,
    "first_name": "Igor",
    "last_name": "Shesterkin",
    "sweater_number": 31,
    "position": "G",
    "birth_date": "1995-12-30",
    "birth_country": "RUS",
    "shoots_catches": "L"
  }
}
```


## View a player

`GET /player/{id}`

Response:

```json
{
  "player": {
    "id": 2,
    "is_active": true,
    "current_team_id": 2,
    "first_name": "Macklin",
    "last_name": "Celebrini",
    "sweater_number": 71,
    "position": "C",
    "birth_date": "2006-05-12",
    "birth_country": "CAN",
    "shoots_catches": "L",
    "skater_stats": {
      "current_season": {
        "games_played": 1,
        "goals": 0,
        "assists": 1,
        "points": 1,
        "pim": 0
      },
      "career_totals": {
        "games_played": 1,
        "goals": 0,
        "assists": 1,
        "points": 1,
        "pim": 0
      }
    },
    "team_full_name": "San Jose Sharks",
    "team_short_name": "SJS"
  }
}
```

Notes:

- Stats are returned on `GET /players/{id}` when a row exists in `player_stats`.
- Stats stay basic for now
- `GET /v0/updateStats` rebuilds stats
- Right now it only does skater stats

## View all players

> needs to be updated, there are filters / pagination

`GET /players`

Response:

```json
{
  "metadata": {
    "current_page": 1,
    "page_size": 20,
    "first_page": 1,
    "last_page": 1,
    "total_records": 2
  },
  "players": [
    {
      "id": 1,
      "version": 1,
      "is_active": true,
      "current_team_id": 1,
      "first_name": "Igor",
      "last_name": "Shesterkin",
      "sweater_number": 31,
      "position": "G",
      "birth_date": "1995-12-30",
      "birth_country": "RUS",
      "shoots_catches": "L",
      "team_full_name": "New York Rangers",
      "team_short_name": "NYR"
    },
    {
      "id": 2,
      "version": 1,
      "is_active": true,
      "current_team_id": 2,
      "first_name": "Macklin",
      "last_name": "Celebrini",
      "sweater_number": 71,
      "position": "C",
      "birth_date": "2006-05-12",
      "birth_country": "CAN",
      "shoots_catches": "L",
      "team_full_name": "San Jose Sharks",
      "team_short_name": "SJS"
    },
  ]
}
```

## Update a player

`PATCH /player/{id}`

Request:

```json
{
  "first_name": "Connor",
  "last_name": "McDavid",
  "sweater_number": 31,
  "position": "G",
  "birth_date": "1995-12-30",
  "birth_country": "RUS",
  "shoots_catches": "L",
  "current_team_id": 1,
  "is_active": true
}
```

Response:

```json
{
  "player": {
    "id": 1,
    "version": 3,
    "is_active": true,
    "current_team_id": 1,
    "first_name": "Connor",
    "last_name": "McDavid",
    "sweater_number": 31,
    "position": "G",
    "birth_date": "1995-12-30",
    "birth_country": "RUS",
    "shoots_catches": "L"
  }
}
```

## Delete a player

`DELETE /player/{id}`

Example response:

```json
{
  "message": "player successfully deleted"
}
```