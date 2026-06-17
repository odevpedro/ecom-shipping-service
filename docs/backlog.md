# Backlog — ecom-shipping-service

> Registro vivo do progresso do projeto. Atualizado a cada mudança de estado de uma funcionalidade.
> **Última atualização:** 2026-06-16

---

## Sobre o Projeto

API REST para cálculo de frete e rastreamento de entregas em plataforma de e-commerce. Estima custo e prazo com base em CEP de origem/destino e peso do pacote. Fornece rastreamento stub para acompanhamento de pedidos.

**Versão atual:** `0.1.0`
**Repositório:** [github.com/odevpedro/ecom-shipping-service](https://github.com/odevpedro/ecom-shipping-service)
**Stack principal:** Go 1.22 + gorilla/mux

---

## Legenda

| Símbolo | Significado |
|---------|-------------|
| `[ ]`   | Pendente |
| `[~]`   | Em andamento |
| `[x]`   | Concluído |
| `P0`    | Crítico — bloqueia outras features |
| `P1`    | Alta prioridade |
| `P2`    | Média prioridade |
| `P3`    | Melhoria / nice-to-have |
| `XS` `S` `M` `L` `XL` | Estimativa de complexidade |

---

## Em Andamento

> Features atualmente sendo desenvolvidas. Idealmente, máximo de 2–3 itens simultâneos.

_Nenhum item em andamento no momento._

---

## Pendentes

> Ordenadas por prioridade. Itens de P0 e P1 devem entrar em "Em Andamento" primeiro.

| Prioridade | Complexidade | Feature | Observação |
|------------|--------------|---------|------------|
| P1 | M | Rastreamento real com webhook | Consumir webhook de status dos Correios em vez de eventos mock |
| P3 | S | Cálculo por dimensões volumétricas | Incorporar altura, largura e comprimento no preço final |
| P3 | S | Cache de distância entre CEPs | Evitar recalcular `estimateDistance` para pares repetidos |

---

## Concluídas

> Features finalizadas com suas respectivas datas de conclusão e links de referência.

| Feature | Data | Entrega |
|---------|------|---------|
| Cálculo de frete por CEP + peso | 2026-06-16 | `POST /api/shipping/calculate` — heurística de prefixo CEP + fórmula `peso * 50 + distancia * 1` |
| Rastreamento stub | 2026-06-16 | `GET /api/shipping/{orderId}/track` — retorna 3 eventos mock fixos |
| Health checks | 2026-06-16 | Endpoints `GET /health`, `/live`, `/ready` |
| Persistência em PostgreSQL | 2026-06-17 | Camada `internal/repository` com DDL para `shipping_quotes` e `tracking_events`, conexão via `database/sql` + `lib/pq` |
| Estrutura para transportadoras reais | 2026-06-17 | Interface `Carrier` em `internal/service/carrier.go`, `StubCarrier`, `CorreiosCarrier` scaffold |
| Testes de handler (integração) | 2026-06-17 | `internal/handler/shipping_test.go` — 4 cenários com `httptest` |
| Multi-stage Docker build | 2026-06-16 | `golang:1.22-alpine` → `alpine:3.20` (~20 MB final) |
| Middleware de Request ID | 2026-06-16 | Cabeçalho `X-Request-ID` gerado automático ou herdado da request |
| Erro padronizado | 2026-06-16 | Envelope `{ data, error: { code, message, details }, meta }` |

---

## Bugs Conhecidos

> Problemas identificados que ainda não foram corrigidos.

_Nenhum bug reportado._

---

## Notas & Decisões Pendentes

> Pontos em aberto que precisam de decisão antes de serem desenvolvidos.

- Escolha da transportadora real: Correios SIGEP Web vs API unificada (como Frenet ou Melhor Envio)
- Estratégia de cache para `estimateDistance`: Redis vs mapa em memória com TTL
- Schema de banco: modelagem separada por transportadora ou tabela única com discriminator

---

## Histórico de Versões

| Versão | Data | Principais entregas |
|--------|------|---------------------|
| `0.1.0` | 2026-06-16 | Cálculo de frete, rastreamento stub, health checks, Docker |
