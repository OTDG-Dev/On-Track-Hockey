# Divisons

## Create a divison

`POST /divisions`

Request:

```json
{
  "league_id": 1,
  "name": "Div 1"
}
```

Response:

```json
{
  "division": {
    "id": 18,
    "league_id": 1,
    "name": "Div 1"
  }
}
```

## View a division

`GET /division/{id}`

Example response:

```json
{
  "division": {
    "id": 1,
    "league_id": 1,
    "name": "Div 1"
  }
}
```
## View all divisions

`GET /divisions`

Example response:

```json
{
  "divisions": [
    {
      "id": 1,
      "league_id": 1,
      "name": "Div 1"
    },
    {
      "id": 2,
      "league_id": 1,
      "name": "Div 2"
    },
  ]
}
```

## View all teams from a division

`GET /divisions/{division_id}/teams`

Example response:

```json
{
  "teams": [
    {
      "id": 1,
      "full_name": "New York Rangers",
      "short_name": "NYR",
      "division_id": 1,
      "is_active": true
    },
    {
      "id": 2,
      "full_name": "San Jose Sharks",
      "short_name": "SJS",
      "division_id": 1,
      "is_active": true
    }
  ]
}
```

## Update a division

`PATCH /division/{id}`

Request:

```json
{
  "name": "AwesomeDiv"
}
```

Response:

```json
{
  "division": {
    "id": 2,
    "league_id": 1,
    "name": "AwesomeDiv"
  }
}
```

## Delete a division

`DELETE /division/{id}`

Example response:

```json
{
  "message": "division successfully deleted"
}
```