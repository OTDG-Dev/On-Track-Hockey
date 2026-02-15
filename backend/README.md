# Backend

## API Status

### Endpoints

| Method | URL             | Action            | Status | Example                                                                      |
|--------|-----------------|-------------------|--------|------------------------------------------------------------------------------|
| GET    | /v1/healthcheck | app healthcheck   |        | `curl localhost:3000/v1/healtcheck`                                          |
| GET    | [/v1/players](#endpoint-get-players)    | Show all players  | WIP    | `curl "localhost:3000/v1/players?page=4&page_size=40&sort=first_name"` |
| POST   | /v1/players     | Create new player |        |                                                                              |
| PATCH  | /v1/players/:id | Update a player   |        |                                                                              |
| GET    | /v1/players/:id | Show a player     |        |                                                                              |
| DELETE | /v1/players/:id | Delete a player   |        |                                                                              |



#### `GET /v1/players` {#endpoint-get-players}


### Parameters:

| key          | input             | example             | default |
|--------------|-------------------|---------------------|---------|
| `page`       | page of results   | `page=4`            | 1       |
| `page_size`  | results per page  | `page_size=40`      | 20      |
| `first_name` | user's first name | `first_name=connor` | nil     |

### Sorting

`/v1/players?sort=<sort_key>`

Prefix any search term with `-` to reverse the order

- `first_name`


### Examples:

```bash
/v1/players?first_name=connnor
```

