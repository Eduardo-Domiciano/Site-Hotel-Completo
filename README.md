# Mikhalateia Hotel — projeto completo

Projeto completo com front Angular, API Go e PostgreSQL. Deploy local via Docker Compose.

- Demo só front: https://github.com/Eduardo-Domiciano/SiteHotel

## Estrutura

```
Site-Hotel-Completo/
  frontend/     # Angular
  backend/      # API Go + migrations
  docker/       # Nginx (serve SPA + proxy /api)
  docker-compose.yml
```

## Subir tudo

```bash
cp .env.example .env
docker compose up --build
```

Abre: http://localhost:8080

- Site e API no mesmo host (`/api/v1/...` via nginx)
- Health: http://localhost:8080/health
- Postgres: `localhost:5432`

## Proposta

Um site de hotel com reservas de quartos e salões de eventos. Perfil de usuario, com edição e cancelamento de reservas por aprovação. Ainda estou estudando formas eficientes de integração.
