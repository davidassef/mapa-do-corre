# Mapa do Corre

Mapa do Corre e uma plataforma de impacto social para dar visibilidade a microempreendedores informais e prestadores autonomos em Fortaleza. A proposta e simples: encontrar corres proximos no mapa, chamar pelo WhatsApp e medir conexoes geradas para o relatorio de extensao.

## O que o MVP entrega

- Busca de corres por proximidade usando raio geografico.
- Mapa interativo com Leaflet e OpenStreetMap.
- Cadastro publico sem conta, validado por codigo enviado por e-mail.
- Remocao do corre sem login, tambem validada por codigo enviado ao e-mail responsavel.
- Contato por WhatsApp com registro de clique.
- Dashboard de impacto com conexoes, prestadores ativos, corres removidos e recorte por categoria.
- Backend preparado para PostgreSQL/PostGIS, rate limit, CORS e SMTP real.
- Deploy de frontend na Vercel e backend em VPS com Docker Compose.

## Stack

- Backend: Go 1.25 + Fiber
- Banco: PostgreSQL 16 + PostGIS 3.4
- Frontend: React 18 + Vite + TypeScript + Tailwind CSS v3
- Mapa: Leaflet + OpenStreetMap
- Estado local: Zustand
- Deploy frontend: Vercel
- Deploy backend: Docker Compose na VPS

## Estrutura

```text
.
├── backend
│   ├── cmd/api
│   ├── internal
│   ├── migrations
│   ├── seeds
│   └── Dockerfile
├── frontend
│   ├── src
│   └── vercel.json
├── .github/skills
├── docker-compose.yml
├── docker-compose.vps.yml
└── Makefile
```

## Ambiente local

1. Crie os arquivos locais de ambiente a partir dos exemplos:

```bash
cp .env.example .env
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

2. Suba o PostGIS local:

```bash
docker compose up -d
```

3. Aplique as migrations:

```bash
docker exec -i mapa-do-corre-postgres psql -U postgres -d mapa_do_corre < backend/migrations/001_init_postgis.sql
docker exec -i mapa-do-corre-postgres psql -U postgres -d mapa_do_corre < backend/migrations/002_email_verificacao.sql
```

4. Carregue o seed de Fortaleza:

```bash
docker exec -i mapa-do-corre-postgres psql -U postgres -d mapa_do_corre < backend/seeds/001_fortaleza_prestadores.sql
```

5. Rode a API:

```bash
cd backend
go mod tidy
go run ./cmd/api
```

6. Rode o frontend em outro terminal:

```bash
cd frontend
npm install
npm run dev
```

Por padrao, o frontend local fica em `http://localhost:5173` e a API em `http://localhost:8080`.

## Variaveis de ambiente

### Backend

```env
APP_ENV=development
PORT=8080
FRONTEND_ORIGIN=http://localhost:5173
PERSISTENCE_MODE=postgres
DATABASE_URL=postgres://postgres:postgres@localhost:5432/mapa_do_corre?sslmode=disable
SMTP_HOST=
SMTP_PORT=
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM_EMAIL=
SMTP_FROM_NAME=Mapa do Corre
EMAIL_CODE_TTL_MINUTES=10
EMAIL_CODE_MAX_ATTEMPTS=5
WRITE_RATE_LIMIT_MAX=20
WRITE_RATE_LIMIT_WINDOW_SECONDS=60
```

O backend tambem aceita os aliases usados no LotoScore/Brevo:

- `SMTP_USER` como alternativa a `SMTP_USERNAME`
- `SMTP_PASS` como alternativa a `SMTP_PASSWORD`
- `EMAIL_FROM` como alternativa a `SMTP_FROM_EMAIL`
- `EMAIL_FROM_NAME` como alternativa a `SMTP_FROM_NAME`

Em `APP_ENV=development`, se o SMTP nao estiver configurado, a API devolve `debugCodigo` para teste manual. Em producao, configure SMTP real.

### Frontend

```env
VITE_API_URL=http://localhost:8080
```

Em producao, use a URL publica da API sem barra final:

```env
VITE_API_URL=https://api.lotoscore.com.br/mapa-do-corre
```

## Endpoints

### `GET /health`

Retorna o estado da API e do banco.

### `GET /corres?lat={latitude}&lon={longitude}&raioMetros={raio}&categoria={categoria}`

Lista prestadores dentro do raio informado.

### `POST /prestadores`

Solicita codigo para publicar um corre.

```json
{
	"nome": "Rita Faxina Express",
	"categoria": "diarista",
	"descricao": "Faxina residencial por diaria.",
	"whatsApp": "85999998888",
	"bairro": "Benfica",
	"email": "rita@email.com",
	"latitude": -3.743269,
	"longitude": -38.536936
}
```

### `POST /prestadores/confirmacoes`

Confirma publicacao com o codigo recebido.

```json
{
	"solicitacaoId": "uuid-ou-id-local",
	"codigo": "123456"
}
```

### `POST /prestadores/:id/cliques`

Registra clique no contato por WhatsApp.

### `POST /prestadores/:id/remocao`

Solicita codigo para remover um corre.

```json
{
	"email": "rita@email.com"
}
```

### `POST /prestadores/remocoes/confirmacoes`

Confirma remocao por codigo.

```json
{
	"solicitacaoId": "uuid-ou-id-local",
	"codigo": "123456"
}
```

### `GET /impacto/resumo`

Retorna conexoes historicas, prestadores ativos, prestadores removidos e metricas por categoria.

## Banco

O schema fica versionado em `backend/migrations`:

- `001_init_postgis.sql`: extensao PostGIS, prestadores, logs de clique e indices.
- `002_email_verificacao.sql`: e-mail responsavel, remocao logica e solicitacoes de codigo.

O seed local fica em `backend/seeds/001_fortaleza_prestadores.sql` e cria 5 prestadores de demonstracao em Fortaleza.

Para restaurar o banco local depois de testes manuais:

```bash
docker exec -i mapa-do-corre-postgres psql -U postgres -d mapa_do_corre -c "DELETE FROM solicitacoes_remocao_prestador; DELETE FROM solicitacoes_cadastro_prestador; DELETE FROM logs_cliques; DELETE FROM prestadores;"
docker exec -i mapa-do-corre-postgres psql -U postgres -d mapa_do_corre < backend/seeds/001_fortaleza_prestadores.sql
```

## Verificacao

```bash
make test
make build
```

Ou diretamente:

```bash
cd backend && go test ./...
cd frontend && npm run build
```

## Deploy do frontend na Vercel

O projeto da Vercel deve usar `frontend` como raiz.

```bash
cd frontend
vercel link --yes --project mapa-do-corre
vercel env add VITE_API_URL production
vercel env add VITE_API_URL preview
vercel --prod --yes
```

Valor esperado para `VITE_API_URL` em producao:

```text
https://api.lotoscore.com.br/mapa-do-corre
```

O diretorio `.vercel/` e local e nao deve ser versionado.

## Deploy do backend na VPS

O backend usa `docker-compose.vps.yml`, que sobe apenas a API e reutiliza a rede compartilhada `shared-db-network`. O Postgres compartilhado deve estar acessivel pelo alias `postgres`.

Comandos principais:

```bash
make bootstrap-vps
make deploy-backend
make status-vps
make logs-backend
```

Diretorio esperado na VPS:

```text
/opt/mapa-do-corre
```

Variaveis essenciais em `backend/.env` na VPS:

```env
APP_ENV=production
PORT=8080
FRONTEND_ORIGIN=https://mapa-do-corre.vercel.app
PERSISTENCE_MODE=postgres
DATABASE_URL=postgres://mapa_do_corre:<senha>@postgres:5432/mapa_do_corre?sslmode=disable
SMTP_HOST=<host-brevo>
SMTP_PORT=<porta-brevo>
SMTP_USERNAME=<usuario-brevo>
SMTP_PASSWORD=<senha-brevo>
SMTP_FROM_EMAIL=<email-remetente>
SMTP_FROM_NAME=Mapa do Corre
EMAIL_CODE_TTL_MINUTES=10
EMAIL_CODE_MAX_ATTEMPTS=5
WRITE_RATE_LIMIT_MAX=20
WRITE_RATE_LIMIT_WINDOW_SECONDS=60
```

Nao versionar `backend/.env`.

## Proxy publico da API

A URL publica prevista para a API e:

```text
https://api.lotoscore.com.br/mapa-do-corre
```

O proxy deve remover o prefixo `/mapa-do-corre` antes de enviar para o container. Assim, `https://api.lotoscore.com.br/mapa-do-corre/health` chega ao backend como `/health`.

## Skills do projeto

Foram adicionadas duas skills em `.github/skills`:

- `deploy-backend-vps`: fluxo de deploy e diagnostico da API na VPS.
- `deploy-frontend-vercel`: fluxo de link, env e publicacao do frontend na Vercel.

## Cuidados operacionais

- Nunca imprimir valores de `.env` no terminal ou em commits.
- Nunca subir outro Postgres para producao; reutilize o `shared-postgres`.
- Nao executar DDL ou migrations em producao sem autorizacao explicita em caixa alta.
- Deploy na VPS deve ser via push e pull do Git, nao por copia manual de arquivos.
- Ajustar `FRONTEND_ORIGIN` no backend sempre que a URL final da Vercel mudar.
