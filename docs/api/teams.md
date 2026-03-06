# Teams

## Create a team

`POST /teams`

Request:

```json
{
  "full_name": "New York Rangers",
  "short_name": "NYR",
  "division_id": 1,
  "is_active": true
}
```

Response:

```json
{
  "team": {
    "id": 3,
    "full_name": "New York Rangers",
    "short_name": "NYR",
    "division_id": 1,
    "is_active": true
  }
}
```

## Get a team

`GET /teams/{id}`

Response:

```json
{
  "team": {
    "id": 1,
    "version": 2,
    "full_name": "New York",
    "short_name": "NYR",
    "division_id": 1,
    "is_active": false
  }
}
```

## Get all teams

`GET /teams`

Response:

```json
{
  "teams": [
    {
      "id": 1,
      "version": 0,
      "full_name": "New York Rangers",
      "short_name": "NYR",
      "division_id": 1,
      "is_active": true
    },
    {
      "id": 2,
      "version": 0,
      "full_name": "San Jose Sharks",
      "short_name": "SJS",
      "division_id": 1,
      "is_active": true
    }
  ]
}
```

## Edit a Team

`PATCH /teams/{id}`

Partial updates are accepted

Request:

```json
{
  "full_name": "New York Rangers",
  "short_name": "NYR",
  "division_id": 1,
  "is_active": true
}
```

Response:

```json
{
  "team": {
    "id": 1,
    "full_name": "New York",
    "short_name": "NYR",
    "division_id": 1,
    "is_active": false
  }
}
```

## Delete a team

`DELETE /teams/{id}`

Response:

```json
{
  "message": "team successfully deleted"
}
```