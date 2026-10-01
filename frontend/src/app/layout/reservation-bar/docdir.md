# reservation-bar

Barra fixa no rodape para escolher datas e ocupacao da reserva.

## Arquivos

| Arquivo | O que faz |
| --- | --- |
| `reservation-bar.ts` | Logica do componente: liga os campos ao servico de reserva e abre o modal. |
| `reservation-bar.html` | Layout da barra: check-in, check-out, adultos, criancas e botao Reservar. |
| `reservation-bar.css` | Estilos da barra fixa no rodape. |
| `docdir.md` | Documentacao desta pasta. |

## Observacoes

- Nao aparece na pagina de login.
- Depende de `ReservationService` em `core/services/`.
