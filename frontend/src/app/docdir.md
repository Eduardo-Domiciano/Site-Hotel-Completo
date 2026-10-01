# app

Nucleo da aplicacao: shell, rotas e organizacao por core, layout e pages.

## Arquivos

| Arquivo | O que faz |
| --- | --- |
| `app.ts` | Componente raiz: monta nav, outlet, rodape, barra e modal. |
| `app.html` | Template do shell da aplicacao. |
| `app.css` | Estilos do container raiz e do `main`. |
| `app.config.ts` | Providers (router, locale pt-BR). |
| `app.routes.ts` | Rotas: home, sobre, acomodacoes, login e trabalhe conosco. |
| `docdir.md` | Documentacao desta pasta. |

## Subpastas

| Pasta | O que guarda |
| --- | --- |
| `core/` | Dados, modelos e servicos compartilhados. |
| `layout/` | Componentes de casca (nav, footer, reserva). |
| `pages/` | Paginas de cada rota. |
