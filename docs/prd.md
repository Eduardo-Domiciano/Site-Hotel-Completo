# PRD — Mikhalateia Hotel

Documento de requisitos do produto. Descreve para quem o site existe, o que ele deve fazer e o que já está no ar.

| | |
| --- | --- |
| Produto | Mikhalateia Hotel |
| Estado | Rascunho |
| Data | 12 de outubro de 2026 |
| Público | Hóspede, candidato e operação do hotel |

O hotel é fictício, em Varginha/MG. O site é o canal de apresentação, reserva de quartos, pedido de orçamento de eventos e candidatura a vagas.

## Problema

Quem quer se hospedar ou fazer um evento precisa ver o hotel, escolher quarto e deixar os dados sem depender de telefone. Quem já reservou logado precisa rever a reserva no perfil. A operação precisa receber esses pedidos em um lugar só.

Hoje o site cobre a vitrine e o registro do pedido. Ainda não cobre disponibilidade real, pagamento, e-mail de confirmação nem a aprovação de cancelamento.

## Personas

| Pessoa | O que quer |
| --- | --- |
| Visitante | Conhecer o hotel, ver quartos e reservar sem criar conta |
| Hóspede com conta | Entrar, reservar, editar o perfil e ver as próprias reservas |
| Organizador de evento | Pedir orçamento de cerimônia ou reunião |
| Candidato | Enviar currículo para uma área |
| Operação | Ver pedidos e, no futuro, aprovar ou recusar alterações e cancelamentos |

Não há perfil de administrador no produto atual.

## Objetivos

1. Apresentar o hotel, as acomodações e as regras de hospedagem.
2. Permitir reserva de 1 a 3 tipos de quarto, com ou sem conta.
3. Guardar conta, reservas, orçamentos e candidaturas.
4. Mostrar ao hóspede logado o detalhe da reserva: quarto, ocupação, diária, total e horários de entrada e saída.

## Fora deste ciclo

- Pagamento e nota fiscal
- Disponibilidade e bloqueio de overbooking
- E-mail automático de confirmação
- Painel da operação
- Edição e cancelamento de reserva com aprovação
- Reserva de salão como produto com preço, no mesmo fluxo dos quartos

O texto do site já fala em confirmação por e-mail, no-show às 18h e cancelamento com multa de 50%. Essas regras são exibidas. O sistema ainda não as executa.

## O que já existe

### Páginas

| Rota | Uso |
| --- | --- |
| `/` | Home: marca, contatos e pilares do hotel |
| `/sobre` | Institucional, comodidades e mapa |
| `/acomodacoes` | Catálogo de quartos e formulário de eventos |
| `/trabalhe-conosco` | Candidatura com currículo |
| `/login` | Entrar ou cadastrar |
| `/perfil` | Editar dados, trocar senha e ver reservas |

A barra de reservas fica no rodapé, exceto no login. Ela define check-in, check-out, adultos e crianças e abre o modal. O modal repete as datas e permite escolher até 3 tipos de quarto.

### Conta

- Cadastro e login com e-mail e senha (mínimo 8 caracteres).
- Sessão com access token de 15 minutos e refresh de 7 dias.
- Perfil: nome, telefone e data de nascimento. O e-mail não muda.
- Troca de senha encerra os refresh tokens da conta.

### Reserva de quarto

- Visitante reserva em `/reservations/guest`, sem vínculo com conta.
- Hóspede logado reserva em `/reservations`. Só essas aparecem em “Minhas reservas”.
- Cada reserva aceita de 1 a 3 tipos, sem repetir o mesmo tipo.
- A diária cobrada é o preço “por” da ocupação de adultos. O total da linha é diária × noites.
- Status gravado: `confirmed`. Não há outro estado.
- Entrada exibida a partir das 14h. Saída exibida até as 12h00. O banco guarda só a data, não o horário.
- No perfil, cada reserva é uma caixa. O clique abre o detalhe: quarto, foto, adultos, crianças, quantidade de diárias, totais e datas com dia da semana.

### Eventos e carreira

- Orçamento: tipo (cerimônia ou reunião), contato, até 2 salas, até 150 pessoas, data e alimentação (café, almoço, jantar).
- Candidatura: área, contato e arquivo PDF, DOC ou DOCX de até 5 MB. O arquivo fica em disco. O banco guarda o caminho.

### Catálogo inicial

Seis tipos: Express Duplo, Express Triplo, Luxo Casal, Luxo Família, Flat e Master. Preços variam de 1 a 4 adultos, conforme o limite do quarto.

## Regras de negócio em vigor

1. Check-out tem de ser depois do check-in. A estadia tem pelo menos 1 noite.
2. Adultos e crianças não podem passar do máximo do tipo de quarto.
3. Crianças entram na ocupação exibida, mas não mudam o preço. O preço usa só a quantidade de adultos.
4. E-mail de usuário é único, sem diferenciar maiúsculas.
5. Reserva de visitante não aparece no perfil de ninguém.
6. Não existe checagem de quarto livre na data.

## Requisitos ainda não feitos

Estes itens estão na proposta do produto e ainda não têm fluxo.

| ID | Requisito | Critério de aceite |
| --- | --- | --- |
| R1 | O hóspede pede alteração de uma reserva sua | O pedido fica pendente até a operação aprovar ou recusar. A reserva vigente não muda antes da aprovação |
| R2 | O hóspede pede cancelamento | O pedido fica pendente. Aprovado, a reserva deixa de valer. Recusado, continua confirmada |
| R3 | A operação vê a fila de pedidos | Existe uma lista de alterações, cancelamentos, orçamentos e candidaturas |
| R4 | Salões de eventos entram como reserva, não só como orçamento | Dá para escolher salão, data e condições, com o mesmo cuidado de aprovação se houver cancelamento |
| R5 | A política de cancelamento passa a valer no sistema | Até 24h antes do check-in, cancelamento sem multa. Depois, 50% do total, se a operação aprovar o cancelamento |

A forma de integrar aprovação, e-mail e operação ainda está em estudo. Não faz parte deste PRD escolher a ferramenta.

## Jornadas

### Reservar sem conta

1. O visitante preenche datas na barra.
2. Abre o modal, informa nome, telefone e e-mail e escolhe os quartos.
3. Confirma.
4. O sistema grava a reserva sem `user_id` e mostra o número e o total.

### Reservar com conta

1. O hóspede entra em `/login`.
2. Segue o mesmo modal. A reserva fica ligada à conta.
3. Em `/perfil`, a lista mostra a caixa da reserva.
4. O clique abre diária, totais, ocupação e horários de entrada e saída.

### Pedir evento ou vaga

1. Em Acomodações, envia o orçamento.
2. Em Trabalhe conosco, envia a área e o currículo.
3. O sistema confirma na tela. Não há acompanhamento depois.

## Dados

PostgreSQL. A API Go é o único acesso. O navegador fala só com `/api/v1`.

| Dado | Onde vive |
| --- | --- |
| Usuário e senha | `users` |
| Sessão longa | `refresh_tokens` (só o hash) |
| Tipos e preços | `room_types`, `room_prices` |
| Reserva e quartos | `reservations`, `reservation_rooms` |
| Orçamento e candidatura | `event_quotes`, `job_applications` |
| Currículo | arquivo em disco |
| Textos e fotos institucionais | arquivos do site, não do banco |

## Requisitos não funcionais

- Sobe local com Docker Compose: site em `http://localhost:8080`, API no mesmo host, Postgres em `localhost:5432`.
- Migrations rodam na subida da API.
- Erros da API voltam como código e mensagem estáveis.
- Auth pública tem limite de 20 pedidos por minuto.
- Idioma da interface: português do Brasil. Moeda: real.

## Como saber se o ciclo atual está ok

- Um visitante conclui uma reserva de quarto e recebe o total.
- Um hóspede logado vê essa reserva no perfil na primeira abertura da página, sem recarregar.
- O detalhe mostra quantidade de diárias, valor da diária, total do quarto e total da reserva.
- Orçamento de evento e candidatura com arquivo válido são aceitos.
- Arquivo maior que 5 MB ou de tipo não permitido é recusado.

## Perguntas em aberto

- Quem aprova alteração e cancelamento, e em qual prazo.
- Se a reserva de visitante pode ser reivindicada depois por uma conta com o mesmo e-mail.
- Se o orçamento de evento vira reserva de salão ou continua só como lead.
- Quando a confirmação por e-mail deixa de ser texto de política e passa a ser mensagem enviada.
