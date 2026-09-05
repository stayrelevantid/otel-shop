# OTel-Shop Lab

**🇬🇧 EN** — Laboratory project to learn **OpenTelemetry distributed tracing** end to end: three Go microservices (Order, Inventory, Payment) deployed on **k3d Kubernetes**, with an OpenTelemetry Collector and Jaeger as the tracing backend.

**🇮🇩 ID** — Project laboratorium untuk belajar **OpenTelemetry distributed tracing** dari awal sampai akhir: 3 microservices Golang (Order, Inventory, Payment) di-deploy ke **k3d Kubernetes**, dengan OpenTelemetry Collector dan Jaeger sebagai backend tracing.

## Architecture / Arsitektur

```
                    ┌──────────────────── k3d cluster (otel-shop) ────────────────┐
                    │                                                             │
host :18080 ───────▶│ order-service ──HTTP──▶ inventory-service ──SQL──▶ Postgres │
host :18081         │   │  otelhttp handler ▲ otelhttp server     otelsql         │
host :18082         │   │  otelhttp client  │                                     │
host :14317/14318   │   ▼                   │                                     │
host :16686         │ payment-service (chaos delay/error %)                       │
                    │   │                                                         │
                    │   ▼ OTLP gRPC 14317                                         │
                    │ otel-collector ──OTLP──▶ jaeger ──UI :16686──▶ host         │
                    └─────────────────────────────────────────────────────────────┘
```

**🇬🇧 EN** — One `POST /checkout` request produces a **single trace ID** across all three services: `POST /checkout` → `validate-order` → `check-inventory` (+ DB spans `sql.*`) → `process-payment` → `POST /pay` (events + baggage `order.id`).

**🇮🇩 ID** — Satu request `POST /checkout` menghasilkan **satu trace ID** lintas ketiga service: `POST /checkout` → `validate-order` → `check-inventory` (+ DB span `sql.*`) → `process-payment` → `POST /pay` (events + baggage `order.id`).

## Repository Structure / Struktur Repo

```
services/          # 3 services (separate Go modules)
  order/           #   POST /checkout — orchestrator
  inventory/       #   GET  /inventory/{id} — reads PostgreSQL
  payment/         #   POST /pay — configurable chaos
pkg/telemetry/     # shared OTel SDK setup (resource, exporter, sampler, propagator)
deploy/            # Kubernetes manifests per component
db/init.sql        # schema + seed products
scripts/           # cluster.sh build.sh deploy.sh test.sh chaos-test.sh check.sh
docs/              # progress-tracker, daily-log, tracing-examples, experiments, glossary
```

**🇬🇧 EN** — Three services are independent Go modules wired together at build time via `replace` directives; the shared telemetry package keeps OTel setup in one place.

**🇮🇩 ID** — Tiga service adalah Go module terpisah yang disambung saat build lewat `replace` directive; package telemetry bersama membuat setup OTel tetap satu tempat.

## Prerequisites / Prasyarat

| Tool | Version | Notes / Catatan |
|---|---|---|
| Go | 1.27 | workspace `go.work`, 4 modules / 4 module |
| Docker | daemon running / daemon aktif | k3d runtime (k3d does not support / tidak support podman) |
| Podman | 6.x | image build (PRD D14) |
| k3d | 5.8+ | local cluster / cluster lokal |
| kubectl | 1.28+ | context `k3d-otel-shop` |
| golangci-lint | 2.x | for `check.sh` (optional) / untuk `check.sh` (opsional) |

## Ports / Port

| Port | Purpose / Tujuan |
|---|---|
| 18080 | order-service (`POST /checkout`, `GET /health`) |
| 18081 | inventory-service (`GET /inventory/{id}`, `GET /health`) |
| 18082 | payment-service (`POST /pay`, `GET /health`) |
| 15432 | PostgreSQL (NodePort → 5432, db `oteldb`) |
| 14317 / 14318 | OTel Collector OTLP gRPC / HTTP |
| 16686 | Jaeger UI |

## Quickstart

```bash
# 1. Cluster + infra (Postgres, Collector, Jaeger)
./scripts/cluster.sh
./scripts/deploy.sh

# 2. Build image (podman → docker → k3d import) + deploy 3 services / 3 service
./scripts/build.sh
./scripts/deploy.sh

# 3. Verification / Verifikasi
./scripts/test.sh          # E2E: health + checkout + error scenarios (9 checks)
./scripts/chaos-test.sh    # chaos: forced error + 2s delay on cluster / di cluster
./scripts/check.sh         # gofmt/vet/lint/test + coverage gate >=70%
```

**🇬🇧 EN** — Then try it:

```bash
curl -X POST http://localhost:18080/checkout \
  -H "Content-Type: application/json" \
  -d '{"item_id":"A123","qty":1}'
# → {"order_id":"O-xxxx","status":"paid"}
```

Open `http://localhost:16686` → service `order-service` → operation `POST /checkout` to see the full distributed trace.

**🇮🇩 ID** — Lalu cobalah:

```bash
curl -X POST http://localhost:18080/checkout \
  -H "Content-Type: application/json" \
  -d '{"item_id":"A123","qty":1}'
# → {"order_id":"O-xxxx","status":"paid"}
```

Buka `http://localhost:16686` → service `order-service` → operation `POST /checkout` untuk melihat distributed trace lengkap.

## Configuration (env) / Konfigurasi (env)

| Env | Default | What it does / Fungsi |
|---|---|---|
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `otel-collector:14317` | OTLP exporter target / tujuan OTLP exporter |
| `OTEL_TRACES_SAMPLER_ARG` | `1.0` | root trace sampling ratio (0.0–1.0) / rasio sampling trace root |
| `PAYMENT_DELAY_PERCENT` / `PAYMENT_DELAY_MS` | `20` / `2000` | chaos delay / jeda chaos |
| `PAYMENT_ERROR_PERCENT` | `10` | chaos error |
| `INVENTORY_URL` / `PAYMENT_URL` | internal DNS | order's downstream targets / target downstream order |

## Documentation / Dokumentasi

- [PRD](prd.md) — requirements & design
- [Implementation Plan](implementation_plan.md) — task breakdown / breakdown task
- [Progress Tracker](docs/progress-tracker.md) · [Daily Log](docs/daily-log.md)
- [Trace Examples / Contoh Trace](docs/tracing-examples.md) — normal / slow / error
- [Experiments / Eksperimen](docs/experiments.md) — 6 core observability experiments / 6 eksperimen inti observability
- [Glossary / Glosarium](docs/glossary.md) — beginner term dictionary (ELI5) / kamus istilah untuk pemula (ELI5)

## Teardown

```bash
k3d cluster delete otel-shop
```