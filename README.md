```bash
# On windows
$env:POSTGRESQL_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
migrate -database $env:POSTGRESQL_URL -path migrations up
```

```bash
# On Linux
export POSTGRESQL_UR="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
migrate -database $env:POSTGRESQL_URL -path migrations up
```

```bash
docker-compose up --build
```
