# Leagues

## Create a league

`POST /leagues`


```json
{
  "name": "NHL"
}
```

## View a league

`GET /leagues/{id}`

Example response:

```json
{
  "league": {
    "id": 2,
    "name": "NHL"
  }
}
```
## View all leagues

`GET /leagues`

Example response:

```json
{
  "leagues": [
    {
      "id": 1,
      "name": "AHL"
    },
    {
      "id": 2,
      "name": "NHL"
    }
  ]
}
```

## Update a league

`PATCH /leagues/{id}`

```json
{
  "name": "KHL"
}
```

## Delete a league

`DELETE /leagues/{id}`

Example response:

```json
{
  "message": "league successfully deleted"
}
```


## View all divisions in a league

`GET /v1/leagues/{league_id}/divisions`

```json
{
  "divisions": [
    {
      "id": 1,
      "league_id": 1,
      "name": "DIV1"
    },
    {
      "id": 2,
      "league_id": 1,
      "name": "DIV2"
    }
  ]
}
```