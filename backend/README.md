# Backend

## API Status

| Method | URL             | Action            | Status | Example                                                                      |
|--------|-----------------|-------------------|--------|------------------------------------------------------------------------------|
| GET    | /v1/healthcheck | app healthcheck   |        | `curl localhost:3000/v1/healtcheck`                                          |
| GET    | /v1/players     | Show all players  | WIP    | `curl "localhost:3000/v1/players?name=connor&page=4&page_size=40&sort=name"` |
| POST   | /v1/players     | Create new player |        |                                                                              |
| PATCH  | /v1/players/:id | Update a player   |        |                                                                              |
| GET    | /v1/players/:id | Show a player     |        |                                                                              |
| DELETE | /v1/players/:id | Delete a player   |        |                                                                              |