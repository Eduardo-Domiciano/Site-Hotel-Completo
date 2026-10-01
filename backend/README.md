# Mikhalateia Hotel API

API Go do monorepo Site-Hotel-Completo. Em produção sobe pelo Compose da raiz.

Para desenvolvimento isolado:

```bash
# na raiz do monorepo
docker compose up db
cd backend
cp ../.env.example .env   # ajuste DATABASE_URL para localhost
# DATABASE_URL=postgres://mikhalateia:mikhalateia@localhost:5432/mikhalateia?sslmode=disable
go run ./cmd/api
```

Ou use só o Compose da raiz (`docker compose up --build`).
