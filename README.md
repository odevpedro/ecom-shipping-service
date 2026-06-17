# ecom-shipping-service

> Cálculo de frete e rastreamento de entregas para plataforma de e-commerce — estima prazos e custos com base em CEP e peso.

[![License](https://img.shields.io/github/license/odevpedro/ecom-shipping-service?style=flat-square)](./LICENSE)
[![Last Commit](https://img.shields.io/github/last-commit/odevpedro/ecom-shipping-service?style=flat-square)](https://github.com/odevpedro/ecom-shipping-service/commits/master)

---

## Sobre o Projeto

API REST responsável pelo cálculo de frete e rastreamento de entregas. Estima custo e prazo com base no CEP de origem/destino e peso do pacote. Fornece dados de rastreamento stub para acompanhamento de pedidos.

Faz parte de um ecossistema **polyglot** de microserviços (Go, Python, Node.js, Java, TypeScript).

---

## Stack & Arquitetura

| Camada        | Tecnologia                          |
|---------------|--------------------------------------|
| Runtime       | Go 1.22                              |
| Framework     | gorilla/mux                          |
| Roteador      | gorilla/mux                          |
| Infra         | Docker + Docker Compose              |
| CI/CD         | GitHub Actions                       |
| Testes        | Go testing (stdlib)                  |

> Padrão arquitetural: **Handler → Service**. Sem ORM — dados em memória (stub).

---

## Estrutura de Pastas

```
cmd/
└── server/main.go                     # Entrypoint — servidor HTTP

internal/
├── config/config.go                   # Carregamento de env vars
├── handler/
│   ├── shipping.go                    # Handlers HTTP (calculate + track)
│   └── middleware.go                  # Request ID + error response helpers
├── model/shipping.go                  # Structs de domínio
├── service/
│   ├── shipping.go                    # Lógica de negócio (frete + rastreio)
│   └── shipping_test.go               # 4 testes
└── repository/                        # (estrutura criada)
```

---

## Como Rodar Localmente

### Pré-requisitos

- Docker + Docker Compose
- Go 1.22+

### Setup

```bash
cp .env.example .env
docker compose up -d
go run ./cmd/server
```

A API estará disponível em `http://localhost:3005`.

### Variáveis de Ambiente

| Variável              | Descrição                            | Valor padrão (dev)                                      |
|-----------------------|--------------------------------------|---------------------------------------------------------|
| `PORT`                | Porta do servidor                    | `3005`                                                  |
| `DATABASE_URL`        | URL de conexão com o PostgreSQL      | `postgresql://ecom:ecom@localhost:5432/ecom_shipping`   |
| `ORDER_SERVICE_URL`   | URL do Order Service                 | `http://localhost:3003`                                 |
| `NODE_ENV`            | Ambiente de execução                 | `development`                                           |

---

## Testes

```bash
go test ./... -v
```

**4 cenários:**
| Suite                  | Arquivo                     | Cenários |
|------------------------|----------------------------|----------|
| Unitários (shipping)   | `service/shipping_test.go`  | 4        |

---

## API — Endpoints

| Método | Rota                                  | Descrição                    |
|--------|---------------------------------------|------------------------------|
| GET    | `/health`                             | Health check                 |
| GET    | `/live`                               | Liveness probe               |
| GET    | `/ready`                              | Readiness probe              |
| POST   | `/api/shipping/calculate`             | Calcula frete                |
| GET    | `/api/shipping/{orderId}/track`       | Rastreia pedido              |

---

## Documentação Técnica

| Documento                                        | Descrição                                 |
|--------------------------------------------------|-------------------------------------------|
| [Fluxos de Funcionalidades](./docs/system-feature-flows.md) | Fluxo interno de cada feature |
| [Modelo de Dados](./docs/data-model.md)          | Entidades, relacionamentos e enums        |
| [Backlog](./backlog.md)                          | Status de desenvolvimento                 |

---

## Status do Projeto

```
[x] Cálculo de frete por CEP + peso
[x] Rastreamento stub com eventos mock
[x] Health checks + Request ID + erro padronizado
[x] Multi-stage Docker build (Go 1.22 → Alpine)
[ ] Persistência em PostgreSQL
[ ] Integração com transportadoras reais
[ ] Testes de integração (handler)
```

---

## Licença

Distribuído sob a licença MIT. Veja [LICENSE](./LICENSE) para mais informações.

---

<p align="center">
  Feito com foco em qualidade por <a href="https://github.com/odevpedro">@odevpedro</a>
</p>
