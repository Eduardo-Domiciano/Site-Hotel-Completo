# Mikhalateia Hotel (frontend)

Front Angular do monorepo completo. Cópia independente do repositório demo `SiteHotel`.

Em Docker, o build de produção usa `apiUrl: '/api/v1'` (mesmo host, proxy Nginx).

## Dev local

```bash
npm install
npm start
```

Requer a API acessível (Compose `db` + `api`, ou stack completa).
