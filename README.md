# Anti-Fraud Service — CQRS Pattern in Go

> Serviço de detecção de fraudes construído sobre o padrão **CQRS (Command Query Responsibility Segregation)**, demonstrando como separar operações de escrita e leitura para escalar cada lado independentemente.

 <img width="802" height="657" alt="image" src="https://github.com/user-attachments/assets/99de2cc3-8131-474c-9cc6-49b01a4d9332" />

---

## Por que o CQRS existe?

### O problema que ele resolve

Em sistemas tradicionais (arquitetura CRUD padrão), um único modelo de dados serve tanto para **escrever** quanto para **ler**. Isso funciona bem em sistemas simples, mas cria gargalos conforme o sistema cresce:

```
┌──────────────────────────────────────┐
│         CRUD Tradicional             │
│                                      │
│  Cliente → [Único modelo] → DB       │
│                                      │
│  Problema: leitura e escrita         │
│  competem pelos mesmos recursos      │
└──────────────────────────────────────┘
```

**Cenários concretos onde isso vira problema:**

| Situação | Consequência |
|---|---|
| Regras de negócio complexas na escrita | Queries de leitura ficam lentas esperando locks |
| Relatórios pesados na leitura | Travamentos nas operações de escrita |
| Escala horizontal | Não é possível escalar só a leitura sem escalar a escrita |
| Modelos conflitantes | Um schema que serve para tudo não serve bem para nada |

### O que o CQRS propõe

Separar completamente o fluxo de **Commands** (escrita, mutação de estado) do fluxo de **Queries** (leitura, projeções de dados):

```
┌──────────────────────────────────────────────────────────┐
│                      CQRS                               │
│                                                          │
│  ┌─────────────┐    Evento     ┌──────────────────────┐ │
│  │  Command    │ ─────────────▶│  Projector           │ │
│  │  Side       │               │  (atualiza read      │ │
│  │  PostgreSQL │               │   model no Redis)    │ │
│  └─────────────┘               └──────────────────────┘ │
│         ▲                               │               │
│         │                               ▼               │
│      Comando                     ┌─────────────┐        │
│         │                        │  Query Side │        │
│         │                        │  Redis      │        │
│  Cliente ──────────────────────▶ └─────────────┘        │
│                   Query                                  │
└──────────────────────────────────────────────────────────┘
```

**Benefícios diretos:**

- **Escalabilidade independente** — o lado de leitura (Redis) pode ter 50 réplicas sem afetar a escrita
- **Modelos otimizados** — o write model garante consistência; o read model é otimizado para cada tela/consulta
- **Performance** — leituras nunca concorrem com escritas por locks de banco
- **Auditoria natural** — cada comando gera um evento, criando um log imutável de tudo que aconteceu
- **Evolução independente** — você pode trocar o banco de leitura (de Redis para Elasticsearch) sem tocar no write model

### Por que em um serviço de antifraude?

O antifraude é um caso de uso ideal para CQRS porque:

1. **Escrita é crítica e complexa** — analisar uma transação envolve múltiplas regras, precisa ser consistente e auditável
2. **Leitura é intensiva** — dashboards, alertas e consultas de risco precisam de latência sub-milissegundo
3. **Volume assimétrico** — para cada 1 escrita (transação), podem existir 10 leituras (consultas de risco, dashboards, etc.)
4. **Histórico importa** — todo evento de fraude precisa ser rastreável

---

## Arquitetura

```
┌────────────────────────────────────────────────────────────────────┐
│                         HTTP API (:8080)                           │
│              chi router + middleware (logger, recover)             │
└───────────┬────────────────────────────────────┬───────────────────┘
            │ Commands                           │ Queries
            ▼                                   ▼
┌─────────────────────┐               ┌──────────────────────┐
│    Command Bus      │               │     Query Bus        │
│  AnalyzeTransaction │               │  GetTransactionRisk  │
│  BlockAccount       │               │  GetAccountStatus    │
│  FlagTransaction    │               │  GetFraudAlerts      │
└────────┬────────────┘               └──────────┬───────────┘
         │                                       │
         ▼                                       ▼
┌─────────────────────┐               ┌──────────────────────┐
│   Fraud Rules Engine│               │   Redis (Read Model) │
│  • Amount Rule      │               │  TransactionRiskView │
│  • Velocity Rule    │               │  AccountStatusView   │
│  • Location Rule    │               │  FraudAlertList      │
│  • Blacklist Rule   │               └──────────────────────┘
└────────┬────────────┘                          ▲
         │                                       │
         ▼                                       │
┌─────────────────────┐    Event     ┌───────────┴──────────┐
│  PostgreSQL         │ ──────────▶  │    Projectors        │
│  (Write Model)      │  (in-memory  │  TransactionProjector│
│  transactions       │   EventBus)  │  AccountProjector    │
│  accounts           │              └──────────────────────┘
└─────────────────────┘
```

### Stack

| Camada | Tecnologia |
|---|---|
| Linguagem | Go 1.25 |
| HTTP Router | [chi v5](https://github.com/go-chi/chi) |
| Write Database | PostgreSQL 16 via [pgx v5](https://github.com/jackc/pgx) |
| Read Database | Redis 7 via [go-redis v9](https://github.com/redis/go-redis) |
| Event Bus | In-memory (substituível por Kafka/RabbitMQ) |
| Logging | [Zap](https://github.com/uber-go/zap) |
| Testes | [testify](https://github.com/stretchr/testify) |
| Load Tests | [k6](https://k6.io) |
| Lint | [golangci-lint](https://golangci-lint.run) |

---

## Estrutura do Projeto

```
cqrs-exemple-golang/
├── cmd/api/main.go                          # Entrypoint — wiring de toda a aplicação
├── internal/
│   ├── config/                              # Configuração via variáveis de ambiente
│   ├── domain/
│   │   ├── entity/                          # Entidades do domínio (Transaction, Account)
│   │   ├── event/                           # Eventos de domínio (TransactionAnalyzed, etc.)
│   │   └── repository/                      # Interfaces dos repositórios de escrita
│   ├── command/
│   │   ├── model/                           # DTOs de comando (AnalyzeTransaction, etc.)
│   │   └── handler/                         # Handlers de comando (lógica de escrita)
│   ├── query/
│   │   ├── model/                           # Views de leitura (TransactionRiskView, etc.)
│   │   ├── repository/                      # Interfaces dos repositórios de leitura
│   │   └── handler/                         # Handlers de query (lógica de leitura)
│   ├── application/
│   │   ├── bus/                             # CommandBus e QueryBus
│   │   └── eventbus/                        # Interface do EventBus
│   ├── fraud/rules/                         # Engine de detecção de fraude e regras
│   └── infrastructure/
│       ├── http/                            # Router, handlers HTTP, middleware, response
│       ├── messaging/inmemory/              # Implementação in-memory do EventBus
│       ├── persistence/postgres/            # Repositórios de escrita (pgx)
│       ├── persistence/redis/               # Repositórios de leitura (go-redis)
│       └── projector/                       # Atualizam o read model a partir de eventos
├── tests/
│   ├── mocks/                               # Mocks testify para testes unitários
│   ├── integration/                         # Testes com Postgres + Redis reais
│   └── load/transaction_load_test.js        # Load test com k6
├── scripts/migrate.sql                      # Schema do banco (write side)
├── docker/api/Dockerfile                    # Imagem multi-stage (final: scratch)
├── docker-compose.yml                       # API + Postgres + Redis
├── docker-compose.override.yml              # Overrides para desenvolvimento local
├── .github/
│   ├── workflows/ci.yml                     # Pipeline CI (lint → test → integration → build)
│   └── PULL_REQUEST_TEMPLATE.md
├── .golangci.yml                            # Configuração do golangci-lint (20+ linters)
├── .env.example                             # Exemplo de variáveis de ambiente
└── Makefile                                 # Comandos de build, test, lint, docker
```

---

## Pré-requisitos

- [Go 1.25+](https://go.dev/dl/)
- [Docker + Docker Compose](https://docs.docker.com/get-docker/)
- [golangci-lint](https://golangci-lint.run/usage/install/) (para lint local)
- [k6](https://k6.io/docs/get-started/installation/) (para load tests)

---

## Configuração

Copie o arquivo de exemplo e ajuste conforme necessário:

```bash
cp .env.example .env
```

### Variáveis de Ambiente

| Variável | Padrão | Descrição |
|---|---|---|
| `HTTP_PORT` | `8080` | Porta do servidor HTTP |
| `LOG_LEVEL` | `info` | Nível de log (`debug`, `info`, `warn`, `error`) |
| `POSTGRES_DSN` | `postgres://postgres:postgres@localhost:5432/antifraude?sslmode=disable` | Connection string do PostgreSQL (write side) |
| `REDIS_ADDR` | `localhost:6379` | Endereço do Redis (read side) |
| `REDIS_PASSWORD` | _(vazio)_ | Senha do Redis |
| `REDIS_DB` | `0` | Número do banco Redis |
| `FRAUD_AMOUNT_THRESHOLD` | `10000` | Valor em reais acima do qual a regra de valor dispara |
| `FRAUD_MAX_TX_PER_HOUR` | `10` | Máximo de transações por hora antes da regra de velocidade disparar |

---

## Como rodar

### Com Docker (recomendado)

```bash
# Sobe toda a stack (API + Postgres + Redis)
docker compose up -d --build

# Ver logs da API
docker compose logs -f api

# Derrubar tudo
docker compose down
```

### Local (sem Docker)

```bash
# 1. Subir somente a infra
docker compose up -d postgres redis

# 2. Aplicar migrations
make migrate

# 3. Rodar a API
make run
```

---

## Makefile

```bash
make build           # Compila o binário em bin/api
make run             # Roda a API diretamente com go run
make test            # Roda os testes unitários
make cover           # Verifica cobertura >= 70% (excluindo persistence)
make cover-html      # Gera coverage.html navegável
make lint            # Executa golangci-lint
make lint-fix        # Aplica correções automáticas de lint
make integration-test # Roda testes de integração (exige POSTGRES_DSN e REDIS_ADDR)
make load-test       # Roda load test com k6
make migrate         # Aplica o schema SQL no PostgreSQL
make docker-up       # Sobe toda a stack com Docker Compose
make docker-down     # Derruba a stack
make docker-clean    # Derruba a stack e remove volumes
make tidy            # go mod tidy
```

---

## API — Endpoints e Payloads

A API base é `http://localhost:8080`.

---

### Health Check

#### `GET /health`

Verifica se o serviço está no ar.

**Response `200 OK`:**
```json
{
  "status": "ok",
  "time": "2026-03-01T12:00:00Z"
}
```

#### `GET /ready`

Verifica se o serviço está pronto para receber tráfego.

**Response `200 OK`:**
```json
{
  "status": "ready",
  "time": "2026-03-01T12:00:00Z"
}
```

---

### Transações

#### `POST /api/v1/transactions`

Submete uma transação para análise de fraude. Executa todas as regras do engine, persiste no PostgreSQL e atualiza o read model no Redis via evento.

**Request body:**
```json
{
  "account_id":  "550e8400-e29b-41d4-a716-446655440000",
  "amount":      1500.00,
  "currency":    "BRL",
  "merchant_id": "merchant-abc-123",
  "location":    "BR",
  "metadata": {
    "channel":    "mobile",
    "device_id":  "dev-xyz"
  }
}
```

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `account_id` | UUID | Sim | ID da conta que originou a transação |
| `amount` | float64 | Sim | Valor da transação (deve ser > 0) |
| `currency` | string | Sim | Código ISO 4217 de 3 letras (`BRL`, `USD`, `EUR`) |
| `merchant_id` | string | Sim | Identificador do estabelecimento |
| `location` | string | Sim | Código de país ou região |
| `metadata` | object | Não | Dados extras em chave-valor |

**Response `201 Created` — baixo risco:**
```json
{
  "transaction_id": "7f3b9a1c-4e2d-4f8a-b1c2-d3e4f5a6b7c8",
  "status":         "approved",
  "risk_score":     0.1,
  "risk_level":     "low"
}
```

**Response `201 Created` — risco médio:**
```json
{
  "transaction_id": "7f3b9a1c-4e2d-4f8a-b1c2-d3e4f5a6b7c8",
  "status":         "flagged",
  "risk_score":     0.65,
  "risk_level":     "medium"
}
```

**Response `201 Created` — alto risco (bloqueada):**
```json
{
  "transaction_id": "7f3b9a1c-4e2d-4f8a-b1c2-d3e4f5a6b7c8",
  "status":         "declined",
  "risk_score":     0.90,
  "risk_level":     "high"
}
```

**Tabela de risco:**

| `risk_score` | `risk_level` | `status` |
|---|---|---|
| 0.00 – 0.49 | `low` | `approved` |
| 0.50 – 0.79 | `medium` | `flagged` |
| 0.80 – 1.00 | `high` | `declined` |

**Response `400 Bad Request`:**
```json
{
  "error":   "Bad Request",
  "message": "amount must be greater than zero"
}
```

---

#### `GET /api/v1/transactions/{id}/risk`

Consulta o resultado da análise de risco de uma transação. Lê diretamente do Redis (read model).

**Response `200 OK`:**
```json
{
  "id":            "7f3b9a1c-4e2d-4f8a-b1c2-d3e4f5a6b7c8",
  "account_id":    "550e8400-e29b-41d4-a716-446655440000",
  "amount":        1500.00,
  "currency":      "BRL",
  "merchant_id":   "merchant-abc-123",
  "location":      "BR",
  "status":        "approved",
  "risk_score":    0.10,
  "risk_level":    "low",
  "fraud_reasons": [],
  "created_at":    "2026-03-01T12:00:00Z",
  "updated_at":    "2026-03-01T12:00:00Z"
}
```

Quando a transação tem razões de fraude detectadas:
```json
{
  "id":         "9a1c7f3b-...",
  "status":     "declined",
  "risk_score": 0.90,
  "risk_level": "high",
  "fraud_reasons": [
    "transaction amount 99999.99 exceeds threshold 10000.00",
    "transaction from suspicious location: XX"
  ]
}
```

**Response `404 Not Found`:**
```json
{
  "error":   "Not Found",
  "message": "transaction not found"
}
```

---

#### `POST /api/v1/transactions/{id}/flag`

Sinaliza manualmente uma transação como suspeita (ação de analista).

**Request body:**
```json
{
  "reason":     "padrão de compra atípico identificado pelo analista",
  "flagged_by": "analyst-001"
}
```

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `reason` | string | Sim | Motivo do flagging manual |
| `flagged_by` | string | Sim | Identificador do analista ou sistema |

**Response `204 No Content`** — sem body.

**Response `404 Not Found`:**
```json
{
  "error":   "Not Found",
  "message": "transaction not found"
}
```

---

### Contas

#### `GET /api/v1/accounts/{id}/status`

Retorna o status e métricas de risco agregadas de uma conta. Lê do Redis (read model atualizado pelos projectors).

**Response `200 OK`:**
```json
{
  "id":                 "550e8400-e29b-41d4-a716-446655440000",
  "user_id":            "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "status":             "active",
  "risk_level":         0.25,
  "total_transactions": 42,
  "flagged_count":      3,
  "declined_count":     1,
  "last_activity_at":   "2026-03-01T11:55:00Z",
  "created_at":         "2025-01-15T08:00:00Z"
}
```

Para conta bloqueada:
```json
{
  "id":         "550e8400-...",
  "status":     "blocked",
  "blocked_at": "2026-03-01T12:00:00Z"
}
```

**Valores possíveis de `status`:** `active` | `flagged` | `blocked`

---

#### `POST /api/v1/accounts/{id}/block`

Bloqueia uma conta manualmente.

**Request body:**
```json
{
  "reason":     "múltiplas tentativas de fraude confirmadas",
  "blocked_by": "fraud-team"
}
```

| Campo | Tipo | Obrigatório | Descrição |
|---|---|---|---|
| `reason` | string | Sim | Motivo do bloqueio |
| `blocked_by` | string | Sim | Identificador de quem bloqueou (operador, sistema) |

**Response `204 No Content`** — sem body.

**Response `500 Internal Server Error`** (conta já bloqueada):
```json
{
  "error":   "Internal Server Error",
  "message": "account 550e8400-... is already blocked"
}
```

---

### Alertas de Fraude

#### `GET /api/v1/fraud/alerts`

Lista as transações com risco >= 0.5 (`medium` ou `high`). Lê do Redis.

**Query params:**

| Parâmetro | Padrão | Descrição |
|---|---|---|
| _(sem params por ora)_ | page=1, limit=20 | Paginação futura |

**Response `200 OK`:**
```json
{
  "alerts": [
    {
      "id":             "a1b2c3d4-...",
      "transaction_id": "7f3b9a1c-...",
      "account_id":     "550e8400-...",
      "amount":         99999.99,
      "currency":       "USD",
      "risk_score":     0.90,
      "risk_level":     "high",
      "reasons": [
        "transaction amount 99999.99 exceeds threshold 10000.00",
        "transaction from suspicious location: XX"
      ],
      "status":     "open",
      "created_at": "2026-03-01T12:00:00Z"
    }
  ],
  "total": 1,
  "page":  1,
  "limit": 20
}
```

**Valores possíveis de `status` do alerta:** `open` | `reviewed` | `dismissed`

---

## Regras de Fraude

O engine avalia cada transação em paralelo contra 4 regras. Os scores são somados e limitados a 1.0.

| Regra | Score | Condição de disparo |
|---|---|---|
| `amount_threshold` | +0.4 | Valor da transação > `FRAUD_AMOUNT_THRESHOLD` |
| `velocity` | +0.5 | Conta fez mais de `FRAUD_MAX_TX_PER_HOUR` transações na última hora |
| `suspicious_location` | +0.5 | País da transação está na lista de locais suspeitos (`XX`, `ZZ` por padrão) |
| `blacklist` | +0.9 – 1.0 | Conta ou merchant está na blacklist |

---

## Testes

```bash
# Unitários (sem dependências externas)
make test

# Com cobertura (mínimo 70%)
make cover

# Relatório HTML de cobertura
make cover-html

# Integração (requer Postgres + Redis rodando)
POSTGRES_DSN="postgres://postgres:postgres@localhost:5432/antifraude?sslmode=disable" \
REDIS_ADDR="localhost:6379" \
make integration-test

# Load test (requer k6 instalado)
make load-test
```

### Cobertura atual

| Camada | Cobertura |
|---|---|
| Domain (entity, event) | 100% |
| Fraud Rules Engine | 97.4% |
| Application Bus | 100% |
| Command Handlers | 86.3% |
| Query Handlers | 100% |
| Infrastructure HTTP | 100% |
| Infrastructure Messaging | 100% |
| Projectors | 76.2% |
| **Total (excluindo persistence)** | **78.1%** |

> A camada `persistence/postgres` e `persistence/redis` são cobertas pelos testes de integração, não pelos testes unitários.

---

## CI/CD

O pipeline GitHub Actions (`ci.yml`) executa em todo PR e push:

```
lint → unit tests (≥70%) → integration tests → docker build
```

| Job | O que faz |
|---|---|
| **Lint** | `golangci-lint` com 20+ linters configurados em `.golangci.yml` |
| **Unit Tests** | `go test` excluindo persistence, verifica cobertura ≥ 70% |
| **Integration Tests** | Sobe Postgres e Redis via `services:`, aplica migration, roda testes |
| **Build** | Build da imagem Docker multi-stage, valida que compila sem erros |

---

## Extensibilidade

### Trocar o Event Bus por Kafka

O `EventBus` é definido por interface em `internal/application/eventbus/event_bus.go`. Para usar Kafka:

1. Crie `internal/infrastructure/messaging/kafka/event_bus.go` implementando a interface
2. Substitua a injeção no `cmd/api/main.go`

### Adicionar uma nova regra de fraude

1. Crie o arquivo em `internal/fraud/rules/minha_regra.go` implementando `rules.Rule`
2. Registre no engine em `cmd/api/main.go`:
   ```go
   rules.NewEngine(
       // regras existentes...
       rules.NewMinhaRegra(...),
   )
   ```

### Adicionar um novo endpoint

1. Crie o command/query model em `internal/command/model/` ou `internal/query/model/`
2. Implemente o handler em `internal/command/handler/` ou `internal/query/handler/`
3. Registre no bus em `cmd/api/main.go`
4. Adicione a rota em `internal/infrastructure/http/router.go`

---

## Licença

MIT
