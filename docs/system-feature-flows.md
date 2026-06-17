# System Feature Flows

> Registro histórico e incremental dos fluxos internos de cada funcionalidade.
> Este documento cresce a cada nova feature implementada e **nunca tem seções removidas**.

---

## Índice

- [Visão Geral da Arquitetura](#visão-geral-da-arquitetura)
- [Convenções deste Documento](#convenções-deste-documento)
- [Feature: Cálculo de Frete](#feature-cálculo-de-frete)
- [Feature: Rastreamento de Pedido](#feature-rastreamento-de-pedido)

---

## Visão Geral da Arquitetura

> API REST monolítica (sem banco) usando roteador gorilla/mux. Dados em memória com stubs. Sem ORM ou camada de persistência.

**Padrão arquitetural:** Handler → Service (2 camadas)

**Fluxo global de uma requisição:**

```
HTTP Request
    └── Handler (roteamento + parse de request)
            └── Service (lógica de negócio)
```

**Camadas e responsabilidades:**

| Camada    | Responsabilidade                                        |
|-----------|---------------------------------------------------------|
| `handler` | Receber requisições, decodificar JSON, validar entrada, formatar resposta, escrever erros padronizados |
| `service` | Regras de negócio: cálculo de frete, estimativa de distância, geração de rastreamento stub |

---

## Convenções deste Documento

- **Erros de domínio** são retornados como respostas HTTP com envelope `{ data, error: { code, message, details }, meta }`
- **DTOs** (`model.CalculateInput`, `model.CalculateOutput`, `model.TrackingOutput`) trafegam entre handler e service
- **Não há transações ou persistência** — dados são processados e descartados em memória
- **Rastreamento é stub** — não consulta API externa nem banco

---

# Feature: Cálculo de Frete

> **Versão:** 1.0.0
> **Implementada em:** 2026-06-16
> **Status:** Concluída

---

## Resumo

Calcula o valor do frete e o prazo estimado de entrega com base no CEP de origem, CEP de destino e peso do pacote. A distância é estimada por heurística de prefixo CEP, e o preço é derivado de fórmula linear.

**Motivação:** Sem esta feature, o e-commerce não consegue exibir custo de frete no checkout — um requisito obrigatório para conversão de vendas.
**Resultado:** O sistema agora calcula frete em ms, sem dependência externa, com precisão suficiente para o MVP.

---

## Fluxo Principal

### 1. Ponto de Entrada

- **Tipo:** HTTP REST
- **Arquivo:** `cmd/server/main.go:42`
- **Rota/Evento:** `POST /api/shipping/calculate`
- **Autenticação:** Pública

O handler `ShippingHandler.Calculate` recebe a requisição, decodifica o body JSON e delega ao serviço.

---

### 2. Validação de Entrada

- **Arquivo:** `internal/handler/shipping.go:20`
- **Biblioteca:** `encoding/json` (stdlib)

| Campo | Tipo | Obrigatório | Regra de validação |
|-------|------|-------------|---------------------|
| `from_cep` | `string` | Não (runtime) | Deve ser string de 8 dígitos (validado indiretamente pelo service) |
| `to_cep` | `string` | Não (runtime) | Deve ser string de 8 dígitos (validado indiretamente pelo service) |
| `weight_kg` | `float64` | Não (runtime) | Deve ser > 0 |

**Falha de validação:** Se o JSON for malformado, retorna `400` com `{ error: { code: "INVALID_REQUEST", message: "invalid request body" } }`. Campos inválidos (CEP curto, weight zero) produzem resultados com preço zero — não há validação explícita de domínio ainda (melhoria pendente).

---

### 3. Orquestração da Aplicação

- **Arquivo:** `internal/service/shipping.go:16`

O método `Calculate` executa em 4 passos:

1. Extrai os prefixos (primeiros 5 dígitos) de `from_cep` e `to_cep`
2. Calcula a distância estimada: `|prefixo_from - prefixo_to| * 50` (mínimo 50 km)
3. Calcula o preço: `peso * 50 + distancia * 1` (resultado em centavos)
4. Calcula os dias estimados: `distancia / 200`, clampado entre 1 e 15

---

### 4. Regras de Negócio

| Regra | Descrição | Localização no Código |
|-------|-----------|----------------------|
| Distância mínima | Se a diferença entre prefixos CEP for < 1, assume 50 km | `service/shipping.go:50` |
| Fórmula de preço | `price = weight * 50 + distance * 1` (centavos) | `service/shipping.go:18-19` |
| Fórmula de prazo | `days = distance / 200`, mínimo 1, máximo 15 | `service/shipping.go:56-65` |
| Transportadora fixa | Sempre `Correios` / `PAC` | `service/shipping.go:25-26` |
| Moeda fixa | Sempre `BRL` | `service/shipping.go:29` |

---

### 5. Persistência / Integrações

**Repositórios utilizados:** Nenhum (cálculo puramente funcional).

**Integrações externas:** Nenhuma (distância estimada por heurística interna).

---

### 6. Resposta Final

**Sucesso — `200`:**

```json
{
  "carrier": "Correios",
  "service_name": "PAC",
  "price_cents": 255,
  "estimated_days": 3,
  "currency": "BRL"
}
```

**Campos retornados:**

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `carrier` | `string` | Transportadora |
| `service_name` | `string` | Nome do serviço |
| `price_cents` | `int` | Valor em centavos (peso * 50 + distância * 1) |
| `estimated_days` | `int` | Prazo em dias (distância / 200, clamp 1–15) |
| `currency` | `string` | Moeda (sempre `BRL`) |

---

## Fluxos Alternativos e Erros

| Cenário | HTTP Status | Código de Erro | Mensagem |
|---------|-------------|----------------|----------|
| JSON inválido no body | `400` | `INVALID_REQUEST` | `invalid request body` |

> Todos os erros retornam o mesmo envelope:
> ```json
> { "data": null, "error": { "code": "ERROR_CODE", "message": "...", "details": {} }, "meta": { "requestId": "abc123" } }
> ```

---

## Diagrama de Sequência

```mermaid
sequenceDiagram
    actor Client
    participant Handler
    participant Service

    Client->>Handler: POST /api/shipping/calculate
    Note over Client,Handler: { from_cep, to_cep, weight_kg }

    Handler->>Handler: Decodifica JSON
    alt JSON inválido
        Handler-->>Client: 400 INVALID_REQUEST
    end

    Handler->>Service: Calculate(input)

    Service->>Service: extractPrefix(from_cep) → 01001
    Service->>Service: extractPrefix(to_cep) → 20020
    Service->>Service: diff = |01001 - 20020| = 10019
    Service->>Service: distance = 10019 * 50 = 500950 km
    Service->>Service: price = weight * 50 + distance * 1
    Service->>Service: days = distance / 200 (clamped 1-15)

    Service-->>Handler: CalculateOutput
    Handler-->>Client: 200 { carrier, service_name, price_cents, estimated_days, currency }
```

---

## Decisões Técnicas

### ADR-001 — Heurística de distância via prefixo CEP

| Campo | Detalhe |
|-------|---------|
| **Status** | Aceita |
| **Data** | 2026-06-16 |
| **Contexto** | Calcular frete sem API de geolocalização nem banco de CEPs. |
| **Decisão** | Usar os primeiros 5 dígitos do CEP como proxy de região geográfica. A diferença absoluta entre prefixos multiplicada por 50 produz uma distância aproximada em km. |
| **Consequências** | Cálculo rápido (~μs), sem dependências externas. Precisão limitada — substituir por integração real de transportadora no futuro. |

---

# Feature: Rastreamento de Pedido

> **Versão:** 1.0.0
> **Implementada em:** 2026-06-16
> **Status:** Concluída

---

## Resumo

Retorna eventos de rastreamento mock para um `orderId`. O rastreamento é completamente stub — dados são fixos e não refletem estado real de entrega.

**Motivação:** Fornecer um contrato de API de rastreamento para que outros serviços (Order Service, frontend) possam consumi-lo e desenvolver suas integrações em paralelo.
**Resultado:** O endpoint `GET /api/shipping/{orderId}/track` está disponível e retorna dados no formato esperado, permitindo o desenvolvimento paralelo das interfaces dependentes.

---

## Fluxo Principal

### 1. Ponto de Entrada

- **Tipo:** HTTP REST
- **Arquivo:** `cmd/server/main.go:43`
- **Rota/Evento:** `GET /api/shipping/{orderId}/track`
- **Autenticação:** Pública

---

### 2. Validação de Entrada

- **Arquivo:** `internal/handler/shipping.go:32`
- **Biblioteca:** `gorilla/mux.Vars`

| Campo | Tipo | Obrigatório | Regra de validação |
|-------|------|-------------|---------------------|
| `orderId` | `string` | Sim | Extraído da URL via mux.Vars |

**Falha de validação:** Se `orderId` estiver vazio, retorna `400` com `{ error: { code: "INVALID_REQUEST", message: "orderId is required" } }`.

---

### 3. Orquestração da Aplicação

- **Arquivo:** `internal/service/shipping.go:33`

O método `Track` executa em 1 passo:

1. Retorna um `TrackingOutput` fixo com 3 eventos mock e status `"in_transit"`

---

### 4. Regras de Negócio

| Regra | Descrição | Localização no Código |
|-------|-----------|----------------------|
| Stub de rastreamento | Eventos são fixos, independentes do orderId ou carrier real | `service/shipping.go:34-43` |
| Status fixo | Sempre `"in_transit"` | `service/shipping.go:37` |

---

### 5. Persistência / Integrações

**Repositórios utilizados:** Nenhum.

**Integrações externas:** Nenhuma (stub).

---

### 6. Resposta Final

**Sucesso — `200`:**

```json
{
  "order_id": "order-123",
  "carrier": "Correios",
  "status": "in_transit",
  "events": [
    {
      "date": "2026-06-14 08:30",
      "location": "São Paulo, SP",
      "description": "Objeto postado"
    },
    {
      "date": "2026-06-15 14:15",
      "location": "Curitiba, PR",
      "description": "Em trânsito para unidade de distribuição"
    },
    {
      "date": "2026-06-16 09:00",
      "location": "Curitiba, PR",
      "description": "Saiu para entrega ao destinatário"
    }
  ]
}
```

**Campos retornados:**

| Campo | Tipo | Descrição |
|-------|------|-----------|
| `order_id` | `string` | Identificador do pedido (ecoado da URL) |
| `carrier` | `string` | Transportadora (sempre `"Correios"`) |
| `status` | `string` | Status atual (sempre `"in_transit"`) |
| `events` | `array` | Lista de eventos de rastreamento |

---

## Fluxos Alternativos e Erros

| Cenário | HTTP Status | Código de Erro | Mensagem |
|---------|-------------|----------------|----------|
| orderId vazio | `400` | `INVALID_REQUEST` | `orderId is required` |

> Todos os erros retornam o mesmo envelope:
> ```json
> { "data": null, "error": { "code": "ERROR_CODE", "message": "...", "details": {} }, "meta": { "requestId": "abc123" } }
> ```

---

## Diagrama de Sequência

```mermaid
sequenceDiagram
    actor Client
    participant Handler
    participant Service

    Client->>Handler: GET /api/shipping/{orderId}/track
    Note over Client,Handler: orderId = "order-123"

    Handler->>Handler: Extrai orderId da URL
    alt orderId vazio
        Handler-->>Client: 400 INVALID_REQUEST
    end

    Handler->>Service: Track("order-123")

    Service-->>Handler: TrackingOutput (eventos fixos)

    Handler-->>Client: 200 { order_id, carrier, status, events }
```

---

## Decisões Técnicas

### ADR-002 — Rastreamento stub para viabilizar desenvolvimento paralelo

| Campo | Detalhe |
|-------|---------|
| **Status** | Aceita |
| **Data** | 2026-06-16 |
| **Contexto** | O Order Service e o frontend dependem do endpoint de tracking, mas a integração real com Correios ainda não estava disponível. |
| **Decisão** | Implementar o endpoint com dados mock fixos, mas com o contrato de resposta idêntico ao que será usado futuramente. |
| **Consequências** | Permite desenvolvimento paralelo. Quando a integração real for implementada, o contrato da API não muda — apenas a fonte dos eventos troca de stub para real. |
