# Backend

## API Status

### Players

| Method | URL                                    | Action            | Status |
|--------|----------------------------------------|-------------------|--------|
| GET    | `/v1/players` | Show all players  |        |
| POST   | [`/v1/players`](#POST-v1players)                         | Create new player |        |
| PATCH  | `/v1/players/:id`                      | Update a player   |        |
| GET    | `/v1/players/:id`                      | Show a player     |        |
| DELETE | `/v1/players/:id`                      | Delete a player   |        |

#### `POST /v1/players`

body:

```json
{
    "first_name": "connor",
    "last_name": "mcdavid",
    "sweater_number": 97,
    "position": "C",
    "birth_date": "1997-01-13",
    "birth_country": "CAN",
    "shoots_catches": "L",
    "current_team_id": 1,
    "is_active": true
}
```


### Teams

| Method | URL             | Action          | Status |
|--------|-----------------|-----------------|--------|
| GET    | `/v1/teams`     | Show all teams  | WIP    |
| POST   | [`/v1/teams`](#post-v1teams)     | Create new team |        |
| GET    | `/v1/teams/:id` | Show a team     |        |
| DELETE | `/v1/teams/:id` | Delete a team   |        |

#### `POST /v1/teams`

body:

```json
{
    "full_name": "Toronto Maple Leafs",         
    "short_name": "TOR",                    // must be 3 chars
    "division_id": 1,
    "is_active": true
}
```

### Divisions

| Method | URL             | Action              | Status |
|--------|-----------------|---------------------|--------|
| GET    |  `/v1/divisions`| Show all divisions  |        |
| POST   | [`/v1/divisions`](#post-v1divisions) | Create new division |        |

#### `POST /v1/divisions`

body:

```json
{
  "name": "divison 1",
  "league_id": 1
}
```

### Leagues

| Method | URL           | Action            | Status |
|--------|---------------|-------------------|--------|
| GET    | `/v1/leagues` | Show all leagues  |        |
| POST   | [`/v1/leagues`](#post-v1leagues) | Create new league |        |


#### `POST /v1/leagues`

body:

```json
{
  "name": "NHL"
}
```