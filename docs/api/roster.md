# Roster

## View a team's roster

`GET /roster/{team_id}`

Request:

```json
{
  "roster": {
    "forwards": [
      {
        "id": 2,
        "is_active": true,
        "current_team_id": 2,
        "first_name": "Macklin",
        "last_name": "Celebrini",
        "sweater_number": 71,
        "position": "C",
        "birth_date": "2006-05-12",
        "birth_country": "CAN",
        "shoots_catches": "L"
      },
      {
        "id": 7,
        "is_active": true,
        "current_team_id": 2,
        "first_name": "William",
        "last_name": "Eklund",
        "sweater_number": 72,
        "position": "LW",
        "birth_date": "2002-10-12",
        "birth_country": "SWE",
        "shoots_catches": "L"
      },
      {
        "id": 8,
        "is_active": true,
        "current_team_id": 2,
        "first_name": "Fabian",
        "last_name": "Zetterlund",
        "sweater_number": 20,
        "position": "RW",
        "birth_date": "1999-08-25",
        "birth_country": "SWE",
        "shoots_catches": "L"
      }
    ],
    "defensemen": [
      {
        "id": 6,
        "is_active": true,
        "current_team_id": 2,
        "first_name": "Mario",
        "last_name": "Ferraro",
        "sweater_number": 38,
        "position": "D",
        "birth_date": "1998-09-17",
        "birth_country": "CAN",
        "shoots_catches": "L"
      }
    ],
    "goalies": [
      {
        "id": 9,
        "is_active": true,
        "current_team_id": 2,
        "first_name": "Mackenzie",
        "last_name": "Blackwood",
        "sweater_number": 29,
        "position": "G",
        "birth_date": "1996-12-09",
        "birth_country": "CAN",
        "shoots_catches": "L"
      }
    ]
  }
}
```
