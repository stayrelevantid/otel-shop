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

Riwayat per hari: [daily-log.md](daily-log.md)