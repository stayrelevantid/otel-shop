# Progress Tracker — OTel-Shop Lab

**🇬🇧 EN** — Phase status against `implementation_plan.md` (v2). Updated after each phase completes.

**🇮🇩 ID** — Status fase mengacu pada `implementation_plan.md` (v2). Diperbarui setelah setiap fase selesai.

Repo: [github.com/stayrelevantid/otel-shop](https://github.com/stayrelevantid/otel-shop)

---

## Phase Status / Status Fase

| Phase / Fase | Scope | Status | Date / Tanggal | Notes / Catatan |
|------|-------|--------|---------|---------|
| F1 | Foundation & Infrastructure Config | ✅ Done | 2026-08-26 | EN: Go workspace (4 modules), K8s manifests (postgres/collector/jaeger), k3d up, smoke trace passed · ID: Scaffolding Go workspace (4 module), manifest K8s, cluster k3d hidup, trace smoke lolos |
| F2 | Application Skeleton (3 services) | ✅ Done | 2026-08-27 | EN: models+handlers+routing, multi-stage Dockerfiles, deployed & E2E-verified · ID: Model+handler+routing lengkap, Dockerfile multi-stage, deploy ke k3d, endpoint tervalidasi via curl |
| F3 | Database Integration (Inventory ↔ PostgreSQL) | ✅ Done | 2026-08-28 | EN: pgx StockStore, handler 200/404/500, deployed & verified · ID: pgx StockStore, handler 200/404/500, deploy + verifikasi E2E |
| F4 | Service Integration (Order → Inventory + Payment) | ✅ Done | 2026-08-28 | EN: client interfaces + full checkout flow, 500+error on downstream failure · ID: client interfaces + checkout flow penuh, 500+error saat downstream gagal |
| F5 | Chaos Engineering (Payment) | ✅ Done | 2026-08-29 | EN: env-driven chaos, deterministic 0/100%, verified in-cluster · ID: chaos.Config dari env, deterministik 0/100%, terverifikasi in-cluster |
| F6 | Unit Testing (business logic) | ✅ Done | 2026-08-29 | EN: mocks + sqlmock; internal coverage ≥70% (82–100%) · ID: mock client/store + sqlmock; coverage internal semua ≥70% |
| F7 | OTel SDK (`pkg/telemetry`) | ✅ Done | 2026-08-30 | EN: resource/OTLP/sampler/propagator, env-driven sampling, wired to 3 services · ID: Init resource/OTLP/sampler/propagator; sampling via env; wired ke 3 service |
| F8 | HTTP Instrumentation & Context Propagation | ✅ Done | 2026-08-30 | EN: otelhttp server (×3) + client transport; 1 trace ID across 3 services in Jaeger · ID: otelhttp server + client transport; 1 trace ID lintas 3 service di Jaeger |
| F9 | Database Instrumentation | ✅ Done | 2026-08-31 | EN: otelsql + pgx; DB spans under inventory · ID: otelsql + pgx; DB span child di inventory |
| F10 | Manual Instrumentation (custom spans) | ✅ Done | 2026-08-31 | EN: validate-order / check-inventory / process-payment with attrs · ID: 3 manual spans + attrs |
| F11 | Attributes, Events, Baggage | ✅ Done | 2026-08-31 | EN: handler attrs, payment events, ERROR status, baggage order.id across services · ID: attrs handlers, events payment, status ERROR, baggage order.id lintas service |
| F12 | Scripts & Quality Pipeline | ✅ Done | 2026-09-05 | EN: check.sh (fmt/vet/lint/test + coverage gate), chaos-test.sh · ID: check.sh + chaos-test.sh; build/deploy/test sudah ada |
| F13 | Integration & E2E Testing | ✅ Done | 2026-09-05 | EN: full-stack integration test + test.sh + chaos-test.sh · ID: integration test full-stack + test.sh + chaos-test.sh |
| F14 | Documentation | ✅ Done | 2026-09-05 | EN: tracing-examples, experiments (6), final README, glossary · ID: tracing-examples, experiments (6), README final, glosarium |
| F15 | DoD Verification | ✅ Done | 2026-09-05 | EN: 15/15 checks passed — sampling 0.1 (3/20), collector-down (checkout still 200), coverage 82–97% · ID: 15 checklist PRD §45 lulus — sampling 0.1, collector-down, coverage 82–97% |

**Legend:** ✅ Done / Selesai · 🟡 Partial / Sebagian · ⏳ Pending / Menunggu

**🇬🇧 EN** — Day-by-day history: [daily-log.md](daily-log.md)

**🇮🇩 ID** — Riwayat per hari: [daily-log.md](daily-log.md)