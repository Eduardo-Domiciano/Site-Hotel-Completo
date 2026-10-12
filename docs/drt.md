# DRT — Mikhalateia Hotel

Documento de requisitos técnicos. Define como o sistema deve ser construído e operado para cumprir o [PRD](prd.md).

| | |
| --- | --- |
| Produto | Mikhalateia Hotel |
| Estado | Rascunho |
| Data | 12 de outubro de 2026 |
| Relacionado | `docs/prd.md` |

O que está marcado como **vigente** já vale no código. O que está marcado como **previsto** só entra quando o requisito correspondente do PRD for implementado.

## Arquitetura vigente

O sistema tem três processos. O navegador não acessa o banco.

| Peça | Tecnologia | Papel |
| --- | --- | --- |
| Site | Angular 22, standalone, signals | Páginas, barra e modais |
| API | Go 1.26, Chi, pgx | Regras, JSON e acesso ao banco |
| Banco | PostgreSQL 16 | Dados persistentes |
| Borda | Nginx 1.27 | Serve o site e encaminha `/api/` e `/health` para a API |

Em produção local, o Compose sobe `db`, `api` e `web`. O site público fica na porta `8080`. A API escuta `8080` dentro da rede Docker e é publicada no host em `8081`. O Postgres fica em `5432`.

O build do site gera `dist/garbage-hotel/browser`. O Nginx copia esse diretório. O `apiUrl` do front, em desenvolvimento e em produção, é `/api/v1`.

## Requisitos de execução

1. A API só sobe se `DATABASE_URL` e `JWT_SECRET` existirem.
2. Na subida, a API aplica as migrations em `backend/migrations/` e só depois abre o pool.
3. O pool usa no máximo 10 conexões.
4. A API encerra com SIGINT ou SIGTERM e dá até 10 segundos para as conexões abertas terminarem.
5. Arquivos copiados para as imagens Docker precisam ser legíveis pelo processo que os serve. O worker do Nginx e o usuário `appuser` (uid 10001) não rodam como dono dos arquivos do host.
6. O build de produção do Angular não baixa fontes na hora do build (`fonts: false`). As fontes continuam no `index.html` e o navegador as carrega.
7. Currículos ficam no volume `uploads`, não na imagem.

## Contrato HTTP vigente

Prefixo: `/api/v1`. Corpo JSON.

Sucesso:

```json
{ "data": {} }
```

Erro:

```json
{ "error": { "code": "CODIGO", "message": "texto em português" } }
```

JSON de entrada rejeita campos desconhecidos. Erro de parse usa o código `INVALID_JSON` e o status 400.

| Situação | Status |
| --- | --- |
| Entrada inválida | 400 |
| Sem autenticação válida | 401 |
| Proibido | 403 |
| Não encontrado | 404 |
| Conflito | 409 |
| Limite de taxa | 429 |
| Falha não prevista | 500, mensagem genérica `erro interno` |

O 500 não devolve o erro interno ao cliente. O processo registra o erro no log.

### Rotas

| Método | Caminho | Auth | Regra |
| --- | --- | --- | --- |
| GET | `/health` | não | Responde `{ "data": { "status": "ok" } }` |
| POST | `/api/v1/auth/register` | não | Até 20 chamadas por minuto por cliente |
| POST | `/api/v1/auth/login` | não | Mesmo limite |
| POST | `/api/v1/auth/refresh` | não | Mesmo limite. Corpo traz `refreshToken` |
| POST | `/api/v1/auth/logout` | não | Revoga o refresh informado |
| GET | `/api/v1/rooms` | não | Catálogo com preços |
| POST | `/api/v1/reservations/guest` | não | Reserva sem `user_id` |
| POST | `/api/v1/reservations` | Bearer | Reserva ligada ao usuário do token |
| GET | `/api/v1/reservations` | Bearer | Só as reservas daquele `user_id` |
| GET | `/api/v1/me` | Bearer | Perfil público, sem hash de senha |
| PUT | `/api/v1/me` | Bearer | Nome, telefone e nascimento. E-mail não muda |
| PUT | `/api/v1/me/password` | Bearer | Troca a senha e revoga todos os refresh do usuário |
| POST | `/api/v1/event-quotes` | não | Grava o orçamento |
| POST | `/api/v1/job-applications` | não | JSON ou `multipart`. Ver upload |

## Autenticação vigente

1. Senha armazenada com bcrypt. A API nunca devolve o hash.
2. Access token é JWT HS256. Claims incluem `uid`, `email`, `sub`, `iat` e `exp`.
3. Validade padrão do access: 15 minutos (`JWT_ACCESS_TTL_MINUTES`).
4. Refresh é um segredo opaco de 32 bytes. O banco guarda só o SHA-256.
5. Validade padrão do refresh: 168 horas (`JWT_REFRESH_TTL_HOURS`).
6. Cada refresh bem-sucedido revoga o token usado e emite um par novo.
7. Rotas privadas exigem `Authorization: Bearer <access>`.
8. O front guarda `mh_access`, `mh_refresh` e `mh_user` no `localStorage`.
9. Em 401, o interceptor tenta um refresh e repete o pedido. Se o refresh falhar, a sessão local é apagada e o usuário vai para `/login`.
10. O guard de `/perfil` libera a rota se existir `mh_user`. Ele não valida o JWT sozinho. A API continua sendo a barreira das rotas privadas.

## Dados vigentes

1. Identificadores de usuário, reserva, orçamento e candidatura são UUID gerados no banco.
2. O id do tipo de quarto é texto estável (`express-duplo`, `express-triplo`, `luxo-casal`, `luxo-familia`, `flat`, `master`).
3. E-mail de usuário é único em minúsculas (`LOWER(email)`).
4. Datas de hospedagem são `DATE`. Check-out tem de ser maior que check-in. Noites são no mínimo 1.
5. Dinheiro usa `NUMERIC`, não ponto flutuante. A diária gravada é `price_to` da ocupação de adultos. O total da linha é diária × noites.
6. A reserva e as linhas de quarto entram na mesma transação. Se uma linha falha, nada fica gravado.
7. No máximo 3 tipos por reserva, sem repetir `roomTypeId`.
8. Adultos e crianças são validados contra `max_adults` e `max_children` do tipo.
9. Status inicial da reserva é `confirmed`. Não há outro valor produzido pela API hoje.
10. Apagar um usuário apaga os refresh tokens. A reserva permanece, com `user_id` nulo.
11. O seed dos seis quartos e dos preços faz parte da migration `000002`.

## Front vigente

1. Aplicação standalone, sem `NgModule`. Locale `pt-BR`.
2. Estado que a tela precisa redesenhar sozinha usa signal. Lista de reservas do perfil é signal: a lista aparece quando a resposta chega, sem novo carregamento da página.
3. Check-in e check-out da barra e do modal são o mesmo estado. Abrir o modal mostra o que já estiver preenchido.
4. Entrada e saída no detalhe da reserva são texto de interface: a partir das 14h e até as 12h00. Esses horários não são colunas no banco.
5. O catálogo estático em `core/data/rooms.ts` é a fonte do modal de reserva e do nome do quarto no perfil. A página de acomodações prefere `GET /rooms` e cai no catálogo estático se a API falhar.
6. Nomes de quarto no perfil saem desse catálogo, pelo `roomTypeId` devolvido pela API.

## Upload vigente

1. Currículo aceita `.pdf`, `.doc` e `.docx`.
2. Tamanho máximo: `MAX_UPLOAD_MB`, padrão 5.
3. O Nginx aceita corpo de até 10 MB.
4. O caminho gravado em `job_applications.resume_path` é local ao servidor. A API não publica o arquivo de volta.

## Segurança vigente

1. CORS só libera as origens de `CORS_ORIGINS`.
2. Cadastro, login e refresh têm limite de 20 pedidos por minuto, em memória no processo da API.
3. Segredos não entram no repositório. O exemplo fica em `.env.example`. O arquivo `.env` fica de fora do Git.
4. A pasta `.cursor/` fica no `.gitignore` e não entra em commit.
5. Não há TLS dentro do Compose. Quem publicar na internet coloca HTTPS na frente do Nginx.

## Verificação vigente

1. `GET /health` pela porta do site responde status ok.
2. Reserva de visitante e reserva autenticada persistem e devolvem total coerente com diária × noites.
3. `GET /reservations` não devolve reserva de outro usuário nem reserva de visitante.
4. Perfil autenticado mostra a lista na primeira abertura.
5. Upload acima do limite ou com extensão fora da lista volta erro de aplicação, não 500 opaco.
6. CI do backend, em `backend/.github/workflows/ci.yml`, roda `go vet`, `go test ./...` e `go build`. A versão de Go do workflow deve ser a mesma do `go.mod` (1.26). Hoje o workflow pede 1.24.

## Requisitos previstos

Valem quando os itens R1 a R5 do PRD forem construídos. Ainda não há tabela, rota nem tela para eles.

1. Alteração e cancelamento não sobrescrevem a reserva confirmada. O pedido fica em estado pendente, com referência à reserva e ao usuário.
2. Aprovar ou recusar é uma transição explícita. A reserva só muda de status depois da aprovação.
3. Toda transição guarda quem decidiu e quando.
4. A fila da operação é uma rota autenticada, separada das rotas do hóspede.
5. Reserva de salão, se entrar no mesmo banco, não reutiliza `room_types` sem um tipo próprio. Quarto e salão têm regras de capacidade diferentes.
6. Envio de e-mail, se for adotado, fica fora do pedido HTTP do hóspede: a API grava o fato e o envio não bloqueia a resposta além do tempo já aceito para a reserva.

A escolha do meio de integração (fila, e-mail, painel) continua em aberto no PRD. Este DRT não fixa fornecedor.
