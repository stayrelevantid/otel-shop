# Progress Tracker — OTel-Shop Lab

Repo: [github.com/stayrelevantid/otel-shop](https://github.com/stayrelevantid/otel-shop)

---

## Status Fase

Mengacu pada `implementation_plan.md` (v2). Update setelah setiap fase selesai.

| Fase | Scope | Status | Tanggal | Catatan |
|------|-------|--------|---------|---------|
| F1 | Foundation & Infrastructure Config | ✅ Done | 2026-08-26 | Scaffolding Go workspace (4 module), manifests K8s (postgres/collector/jaeger), cluster k3d hidup, trace smoke lolos |
| F2 | Application Skeleton (3 services) | ✅ Done | 2026-08-27 | Model+handler+routing lengkap, Dockerfile multi-stage, deploy ke k3d, semua endpoint tervalidasi via curl dari host |
| F3 | Database Integration (Inventory ↔ PostgreSQL) | ✅ Done | 2026-08-28 | pgx StockStore, handler 200/404/500, deploy + E2E verifikasi |
| F4 | Service Integration (Order → Inventory + Payment) | ✅ Done | 2026-08-28 | client interfaces + checkout flow (validate→stock→pay qty*10), 500+error saat downstream gagal |
| F5 | Chaos Engineering (Payment) | ✅ Done | 2026-08-29 | chaos.Config dari env, deterministik 0/100%, handler 500 saat forced error; terverifikasi in-cluster |
| F6 | Unit Testing (business logic) | ✅ Done | 2026-08-29 | Mock clients/store + sqlmock db; coverage internal semua ≥70% (82–100%) |
| F7 | OTel SDK (`pkg/telemetry`) | ✅ Done | 2026-08-30 | Init resource/OTLP/sampler/propagator; sampling via OTEL_TRACES_SAMPLER_ARG; wired ke 3 service |
| F8 | HTTP Instrumentation & Context Propagation | ✅ Done | 2026-08-30 | otelhttp server (3 service) + client transport (order); terverifikasi 1 trace ID lintas 3 service di Jaeger |
| F9 | Database Instrumentation | ✅ Done | 2026-08-31 | otelsql + pgx; DB spans (sql.conn.query/rows) child di inventory |
| F10 | Manual Instrumentation (custom spans) | ✅ Done | 2026-08-31 | validate-order / check-inventory / process-payment + attrs |
| F11 | Attributes, Events, Baggage | ✅ Done | 2026-08-31 | attrs handlers, payment events, ERROR status, baggage order.id lintas service |
| F12 | Scripts & Quality Pipeline | 🟡 Partial | 2026-08-28 | test.sh E2E ada; check/build/chaos belum |
| F13 | Integration & E2E Testing | ⏳ Pending | — | — |
| F14 | Documentation | ⏳ Pending | — | tracing-examples, experiments, README |
| F15 | DoD Verification | ⏳ Pending | — | Gate akhir |

Legend: ✅ Done · 🟡 Partial · ⏳ Pending

Riwayat per hari: [daily-log.md](daily-log.md)