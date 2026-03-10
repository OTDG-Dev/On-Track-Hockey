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