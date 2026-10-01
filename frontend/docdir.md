# SiteHotel

Raiz do projeto Angular do Mikhalateia Hotel: site de apresentacao e reservas ficticias.

## Arquivos

| Arquivo | O que faz |
| --- | --- |
| `README.md` | Visao geral do projeto, features e como rodar. |
| `package.json` | Dependencias npm e scripts (`start`, `build`, `test`). |
| `package-lock.json` | Versoes travadas das dependencias. |
| `angular.json` | Configuracao do Angular CLI (build, serve, assets). |
| `tsconfig.json` | Configuracao base do TypeScript. |
| `tsconfig.app.json` | TypeScript usado na aplicacao. |
| `tsconfig.spec.json` | TypeScript usado nos testes. |
| `vercel.json` | Configuracao de deploy na Vercel. |
| `docdir.md` | Documentacao desta pasta. |

## Subpastas

| Pasta | O que guarda |
| --- | --- |
| `src/` | Codigo-fonte da aplicacao. |
| `public/` | Arquivos estaticos servidos como estao (imagens, favicon). |
| `docs/` | Modelos e guias de documentacao do projeto. |
| `dist/` | Saida do build (gerada; nao versionar logica). |
| `node_modules/` | Pacotes npm instalados (gerada). |

## Observacoes

- Nao documentar o conteudo de `node_modules/`, `dist/` ou `.git/`.
- Cada subpasta relevante tem o seu proprio `docdir.md`.
