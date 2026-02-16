# Backend

## API Status

### Players

| Method | URL                                    | Action            | Status |
|--------|----------------------------------------|-------------------|--------|
| GET    | [`/v1/players`](#endpoint-get-players) | Show all players  |        |
| POST   | `/v1/players `                         | Create new player |        |
| PATCH  | `/v1/players/:id`                      | Update a player   |        |
| GET    | `/v1/players/:id`                      | Show a player     |        |
| DELETE | `/v1/players/:id`                      | Delete a player   |        |



#### `GET /v1/players` {#endpoint-get-players}


##### Parameters:

| key          | input                | example             | default |
|--------------|----------------------|---------------------|---------|
| `page`       | page of results      | `page=4`            | 1       |
| `page_size`  | results per page     | `page_size=40`      | 20      |
| `first_name` | players's first name | `first_name=connor` | ""      |
| `last_name`  | players's last name  | `last_name=connor`  | ""      |
| `position`   | player's position    | `position=lw`       | ""      |


##### Sorting

`/v1/players?sort=<sort_key>`

Prefix any search term with `-` to reverse the order

- `first_name`
- `last_name`
- `position`


##### Examples:

```bash
/v1/players?first_name=sydney&last_name=crosby&sort=-position
```

### Teams

| Method | URL             | Action          | Status |
|--------|-----------------|-----------------|--------|
| GET    | `/v1/teams`     | Show all teams  | WIP    |
| POST   | `/v1/teams`     | Create new team |        |
| GET    | `/v1/teams/:id` | Show a team     |        |
| DELETE | `/v1/teams/:id` | Delete a team   |        |
