# Progress Tracker — OTel-Shop Lab

Repo: [github.com/stayrelevantid/otel-shop](https://github.com/stayrelevantid/otel-shop)

---

## Status Fase

Mengacu pada `implementation_plan.md` (v2). Update setelah setiap fase selesai.

| Fase | Scope | Status | Tanggal | Catatan |
|------|-------|--------|---------|---------|
| F1 | Foundation & Infrastructure Config | ✅ Done | 2026-08-26 | Scaffolding Go workspace (4 module), manifests K8s (postgres/collector/jaeger), cluster k3d hidup, trace smoke lolos |
| F2 | Application Skeleton (3 services) | 🟡 Partial | 2026-08-26 | Baru stub `/health` per service (verify compile+vet). Handler, model, routing, Dockerfile belum |
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