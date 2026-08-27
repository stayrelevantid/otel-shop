# Progress Tracker — OTel-Shop Lab

Repo: [github.com/stayrelevantid/otel-shop](https://github.com/stayrelevantid/otel-shop)

---

## Status Fase

Mengacu pada `implementation_plan.md` (v2). Update setelah setiap fase selesai.

| Fase | Scope | Status | Tanggal | Catatan |
|------|-------|--------|---------|---------|
| F1 | Foundation & Infrastructure Config | ✅ Done | 2026-08-26 | Scaffolding Go workspace (4 module), manifests K8s (postgres/collector/jaeger), cluster k3d hidup, trace smoke lolos |
| F2 | Application Skeleton (3 services) | ✅ Done | 2026-08-27 | Model+handler+routing lengkap, Dockerfile multi-stage, deploy ke k3d, semua endpoint tervalidasi via curl dari host |
| F3 | Database Integration (Inventory ↔ PostgreSQL) | ⏳ Pending | — | Interface `StockStore`, pgx open, product store |
| F4 | Service Integration (Order → Inventory + Payment) | ⏳ Pending | — | Client interfaces + full checkout flow |
| F5 | Chaos Engineering (Payment) | ⏳ Pending | — | Delay/error percent |
| F6 | Unit Testing (business logic) | ⏳ Pending | — | Mock clients/store; coverage ≥70% |
| F7 | OTel SDK (`pkg/telemetry`) | ⏳ Pending | — | Init + sampling digabung |
| F8 | HTTP Instrumentation & Context Propagation | ⏳ Pending | — | otelhttp server + client |
| F9 | Database Instrumentation | ⏳ Pending | — | otelsql |
| F10 | Manual Instrumentation (custom spans) | ⏳ Pending | — | validate/check/process |
| F11 | Attributes, Events, Baggage | ⏳ Pending | — | — |
| F12 | Scripts & Quality Pipeline | ⏳ Pending | — | check/build/deploy/test/chaos |
| F13 | Integration & E2E Testing | ⏳ Pending | — | — |
| F14 | Documentation | ⏳ Pending | — | tracing-examples, experiments, README |
| F15 | DoD Verification | ⏳ Pending | — | Gate akhir |

Legend: ✅ Done · 🟡 Partial · ⏳ Pending

---

## Daily Log

### Day 1 — 2026-08-26 (Infra-first)

**Tujuan:** Fondasi infra + workspace yang bisa diverifikasi dulu sebelum koding serius.

**Selesai:**
- Git init + `.gitignore` + scan gitleaks (clean) + initial commit
- Go workspace 4 module (`order`, `inventory`, `payment`, `telemetry`) + `go.work`
- `.golangci.yml` (v2), `db/init.sql` (schema + seed A123/B123/C123), README skeleton
- Manifest K8s: namespace, Postgres (secret/configmap/deploy/svc NodePort 15432), OTel Collector (OTLP → batch → OTLP ke jaeger:4317), Jaeger all-in-one (UI 16686)
- Script: `cluster.sh`, `deploy.sh`
- Cluster k3d `otel-shop` up (Docker runtime), semua pods Running/Ready
- Smoke test: OTLP/HTTP trace → Collector → muncul di Jaeger API

**Kendala:**
- k3d nolak Podman (wajib Docker daemon) → nyalain Docker Desktop
- Image pull lambat (postgres:18 ~2m38s, jaeger 54s, collector >9m) → timeout `deploy.sh`
- `go build ./...` di root gagal karena multi-module → build per module dengan `-o /dev/null`

**Next (Day 2):** F2 penuh — model, handler (`/checkout`, `/inventory/{id}`, `/pay`), routing, Dockerfile, deploy 3 service ke cluster.

### Day 2 — 2026-08-27 (Application Skeleton + Deploy)

**Tujuan:** 3 service naik ke k3d, semua endpoint dasar tervalidasi dari host.

**Selesai:**
- Model: `CheckoutRequest/Response/ErrorResponse`, `InventoryResponse`, `PayRequest/Response`
- Handler + routing (Go 1.22+ ServeMux): `GET /health`, `POST /checkout`, `GET /inventory/{id}`, `POST /pay` (stub logic sesuai fase)
- Dockerfile multi-stage x3 (golang:1.27-alpine → scratch, CGO off)
- Manifest `deploy/{order,inventory,payment}` deployment+service NodePort 18080-18082 + probe `/health`
- `scripts/build.sh`: podman build → docker load → k3d image import
- `deploy.sh` diperluas: apply apps + rollout status per service
- Semua pods Running; E2E curl hijau (200 checkout `{order_id,status:paid}`, inventory `{stock:10}`, pay `{status:success}`, 400 untuk invalid)

**Kendala:**
- podman save/load menambah prefix `localhost/` pada tag → perlu `docker tag` sebelum `k3d image import`

**Next (Day 3):** F3 DB integration (`StockStore`, pgx) + F4 service integration (checkout flow antar-service).