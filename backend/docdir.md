# sitehotel-api

API REST do Mikhalateia Hotel.

## Arquivos

| Arquivo | O que faz |
| --- | --- |
| `Dockerfile` | Build multi-stage da API. |
| `docker-compose.yml` | Sobe `api` + `postgres`. |
| `.env.example` | Variáveis de ambiente de exemplo. |
| `go.mod` / `go.sum` | Módulo e dependências Go. |
| `README.md` | Como subir e lista de endpoints. |
| `docdir.md` | Documentação desta pasta. |

## Subpastas

| Pasta | O que guarda |
| --- | --- |
| `cmd/` | Entry points (`api`). |
| `internal/` | Código da aplicação. |
| `migrations/` | SQL versionado. |
| `.github/` | CI. |
