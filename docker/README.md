# OnTrackHockey Docker

Bring up containers, rebuild images

```bash
docker compose up --build -d  
```

Stop containers, remove volumes

```bash
docker compose down -v
```

Rebuild select images

```bash
# rebuild frontend
docker compose up -d --build frontend
# rebuild backend
docker compose up -d --build api
```

If you already have a running database volume from before a migration was added, apply migrations with:

```bash
docker compose up -d migrate
```

If you want a clean local database instead, recreate the stack and volumes:

```bash
docker compose down -v
docker compose up --build -d
```

Enter database container

```bash
docker compose exec db sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB"'
```