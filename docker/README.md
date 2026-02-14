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
docker compose up -d --build backend
```