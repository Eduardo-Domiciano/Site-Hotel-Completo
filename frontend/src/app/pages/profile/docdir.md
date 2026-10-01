# profile

Pagina do usuario autenticado: editar perfil, trocar senha e listar reservas.

## Arquivos

| Arquivo | O que faz |
| --- | --- |
| `profile.ts` | Carrega `/me`, salva perfil/senha e lista reservas da API. |
| `profile.html` | Formularios de perfil e senha, lista de reservas e o modal de detalhe. |
| `profile.css` | Estilos do perfil. |
| `docdir.md` | Documentacao desta pasta. |

## Observacoes

- Rota: `/perfil` (protegida por `authGuard`).
