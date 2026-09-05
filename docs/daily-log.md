# Daily Log — OTel-Shop Lab

**🇬🇧 EN** — Day-by-day challenge notes, one entry per working day, consistent format below.

**🇮🇩 ID** — Catatan progres harian challenge, satu entri per hari kerja, format konsisten di bawah.

Repo: [github.com/stayrelevantid/otel-shop](https://github.com/stayrelevantid/otel-shop) ·
Phase status / status fase: [progress-tracker.md](progress-tracker.md)

---

## Template

```markdown
### Day N — YYYY-MM-DD (short theme / tema singkat)

**Goal / Tujuan:** <EN> / <ID>

**Completed / Selesai:**
- 🇬🇧 <item>
- 🇮🇩 <item>

**Metrics & Verification / Metrik & Verifikasi:**
- 🇬🇧 <result> / 🇮🇩 <hasil>

**Obstacles & Fixes / Kendala & Solusi:**
- 🇬🇧 <problem → fix> / 🇮🇩 <masalah → penyelesaian>

**Next / Selanjutnya:** <EN> / <ID>
```

---

## Entries / Entri Harian

### Day 1 — 2026-08-26 (Infra-first)

**Goal / Tujuan:**
- 🇬🇧 Solid infrastructure foundation + a verifiable Go workspace before serious coding.
- 🇮🇩 Fondasi infra + workspace Go yang bisa diverifikasi sebelum koding serius.

**Completed / Selesai:**
- 🇬🇧 Git init, `.gitignore`, gitleaks scan (clean), initial commit
- 🇮🇩 Git init, `.gitignore`, scan gitleaks (clean), initial commit
- 🇬🇧 Go workspace with 4 modules (`order`, `inventory`, `payment`, `telemetry`) + `go.work`
- 🇮🇩 Go workspace 4 module (`order`, `inventory`, `payment`, `telemetry`) + `go.work`
- 🇬🇧 `.golangci.yml` v2, `db/init.sql` (schema + seed A123/B123/C123)
- 🇮🇩 `.golangci.yml` v2, `db/init.sql` (schema + seed A123/B123/C123)
- 🇬🇧 K8s manifests: namespace, Postgres 18 (NodePort 15432), OTel Collector 0.159.0 (OTLP → batch → OTLP to jaeger:4317), Jaeger all-in-one 1.76.0 (UI 16686)
- 🇮🇩 Manifest K8s: namespace, Postgres 18 (NodePort 15432), OTel Collector 0.159.0 (OTLP → batch → OTLP ke jaeger:4317), Jaeger all-in-one 1.76.0 (UI 16686)
- 🇬🇧 `scripts/cluster.sh`, `scripts/deploy.sh`; cluster up, all pods Running
- 🇮🇩 `scripts/cluster.sh`, `scripts/deploy.sh`; cluster hidup, semua pod Running
- 🇬🇧 Smoke trace OTLP/HTTP → Collector → visible in Jaeger API
- 🇮🇩 Smoke trace OTLP/HTTP → Collector → terlihat di Jaeger API

**Metrics & Verification / Metrik & Verifikasi:**
- 🇬🇧 Jaeger UI HTTP 200 at localhost:16686; 3 seeded products readable via psql; smoke trace found via `/api/traces?service=smoke`
- 🇮🇩 Jaeger UI 200 di localhost:16686; 3 produk terbaca via psql; smoke trace ditemukan via `/api/traces?service=smoke`

**Obstacles & Fixes / Kendala & Solusi:**
- 🇬🇧 k3d rejects Podman (Docker daemon required) → started Docker Desktop
- 🇮🇩 k3d menolak Podman (wajib Docker daemon) → nyalakan Docker Desktop
- 🇬🇧 Slow image pulls (collector >9m) caused deploy timeout → re-run; rollout succeeded
- 🇮🇩 Image pull lambat (collector >9m) menyebabkan timeout deploy → jalankan ulang; rollout sukses
- 🇬🇧 `go build ./...` from root failed (multi-module) → build per module with `-o /dev/null`
- 🇮🇩 `go build ./...` dari root gagal (multi-module) → build per module dengan `-o /dev/null`

**Next / Selanjutnya:**
- 🇬🇧 F2 full — models, handlers, routing, Dockerfiles, deploy 3 services.
- 🇮🇩 F2 penuh — model, handler, routing, Dockerfile, deploy 3 service.

---

### Day 2 — 2026-08-27 (Application Skeleton + Deploy)

**Goal / Tujuan:**
- 🇬🇧 Three services up on k3d with base endpoints verified from the host.
- 🇮🇩 Tiga service naik ke k3d dengan endpoint dasar tervalidasi dari host.

**Completed / Selesai:**
- 🇬🇧 Models: `CheckoutRequest/Response/ErrorResponse`, `InventoryResponse`, `PayRequest/Response`
- 🇮🇩 Model: `CheckoutRequest/Response/ErrorResponse`, `InventoryResponse`, `PayRequest/Response`
- 🇬🇧 Handlers + stdlib ServeMux routing: `GET /health`, `POST /checkout`, `GET /inventory/{id}`, `POST /pay` (stub logic per phase)
- 🇮🇩 Handler + routing stdlib ServeMux (Go 1.22+): `GET /health`, `POST /checkout`, `GET /inventory/{id}`, `POST /pay` (logika stub sesuai fase)
- 🇬🇧 Multi-stage Dockerfiles ×3 (golang:1.27-alpine → scratch, CGO off)
- 🇮🇩 Dockerfile multi-stage ×3 (golang:1.27-alpine → scratch, CGO off)
- 🇬🇧 `deploy/{order,inventory,payment}` manifests: Deployment + Service NodePort 18080–18082 + `/health` probes
- 🇮🇩 Manifest `deploy/{order,inventory,payment}`: Deployment + Service NodePort 18080–18082 + probe `/health`
- 🇬🇧 `scripts/build.sh` (podman build → docker load → k3d import); `deploy.sh` extended with apps + rollout status
- 🇮🇩 `scripts/build.sh`: podman build → docker load → k3d image import; `deploy.sh` diperluas: apply apps + rollout status

**Metrics & Verification / Metrik & Verifikasi:**
- 🇬🇧 compile+vet+gofmt OK; pods Running; E2E curl: 3×200 health, checkout `{order_id,status:paid}`, invalid 400, inventory `{stock:10}`, pay `{status:success}` / invalid 400
- 🇮🇩 compile+vet+gofmt OK untuk 3 module; pods Running; E2E curl dari host hijau semua

**Obstacles & Fixes / Kendala & Solusi:**
- 🇬🇧 `podman save/load` added a `localhost/` image tag prefix → retag on the Docker side before `k3d image import`
- 🇮🇩 Tag image berubah jadi `localhost/otel-shop/...` setelah podman save/load → retag di sisi Docker sebelum `k3d image import`

**Next / Selanjutnya:**
- 🇬🇧 F3 DB integration (`StockStore` + pgx) and F4 checkout flow across services.
- 🇮🇩 F3 DB integration (`StockStore` + pgx) dan F4 checkout flow antar-service.

---

### Day 3 — 2026-08-28 (DB Integration + Checkout Flow + E2E)

**Goal / Tujuan:**
- 🇬🇧 Inventory reads real PostgreSQL; Order calls Inventory & Payment over HTTP; end-to-end checkout.
- 🇮🇩 Inventory baca Postgres sungguhan, Order panggil Inventory & Payment via HTTP, alur checkout end-to-end jalan.

**Completed (F3) / Selesai (F3):**
- 🇬🇧 `store.go`: `StockStore` interface + `ErrNotFound`; `db.go`: `Open(DATABASE_URL)` via pgx/v5 stdlib pool + PingContext; `product.go`: `GetStock` with `sql.ErrNoRows` → `ErrNotFound`; handler refactor (200/404/500); `main.go` wiring
- 🇮🇩 Interface `StockStore` + `ErrNotFound`; `Open(DATABASE_URL)` via pgx/v5 stdlib pool + PingContext; map `sql.ErrNoRows` → `ErrNotFound`; handler refactor 200/404/500; wiring `main.go`

**Completed (F4) / Selesai (F4):**
- 🇬🇧 `client/interfaces.go` (`InventoryClient`, `PaymentClient`), HTTP clients (404→`ErrInventory`, non-200→`ErrPayment`), typed errors; `CheckoutHandler`: validate → GetStock → stock check → `Pay(qty*10)` → 200 `{order_id,status:paid}`; downstream failure → 500 + `{"error":"..."}`; env `INVENTORY_URL`/`PAYMENT_URL`
- 🇮🇩 Interface client + HTTP client inventory & payment; typed errors; flow checkout penuh; downstream gagal → 500 + `{"error":"..."}`; env URL di `main.go`

**Completed (test) / Selesai (test):**
- 🇬🇧 `scripts/test.sh` (F12.4 early): health ×3 + inventory + 4 checkout scenarios, non-zero exit on failure
- 🇮🇩 `scripts/test.sh` (F12.4 dipercepat): health ×3 + inventory + 4 skenario checkout, exit non-zero bila gagal

**Metrics & Verification / Metrik & Verifikasi:**
- 🇬🇧 Local E2E (separate ports to avoid NodePort clash) all green; rebuilt+restarted in cluster, NodePort E2E same results; `test.sh`: 9 passed, 0 failed
- 🇮🇩 E2E lokal (port terpisah hindari tabrakan NodePort) hijau; deploy + E2E via NodePort sama; `test.sh`: 9 passed, 0 failed

**Obstacles & Fixes / Kendala & Solusi:**
- 🇬🇧 First local test hit old cluster pods (same NodePorts) → local tests on other ports; final verify via cluster NodePorts
- 🇮🇩 Local test pertama salah sasaran karena port 18080–18082 dipakai NodePort cluster → test lokal pakai port lain; verifikasi lewat NodePort cluster

**Next / Selanjutnya:**
- 🇬🇧 F5 chaos (payment delay/error %) + F6 unit tests (mock client/store, coverage ≥70%).
- 🇮🇩 F5 chaos (payment delay/error %) + F6 unit tests (mock client/store, coverage ≥70%).

---

### Day 4 — 2026-08-29 (Chaos Engineering + Unit Tests)

**Goal / Tujuan:**
- 🇬🇧 Payment gets chaos (delay/error %); all business logic unit-tested ≥70%.
- 🇮🇩 Payment punya chaos (delay/error %), semua business logic ter-cover unit test ≥70%.

**Completed (F5) / Selesai (F5):**
- 🇬🇧 `internal/chaos/chaos.go`: `Config{DelayPercent,DelayMS,ErrorPercent}` from env, `Apply(ctx)` deterministic at 0/100%, thread-safe rng, respects `ctx.Done()`; payment handler refactor: `Handler` struct, chaos before success → 500 `{status:failed}`; `main.go` uses `chaos.FromEnv()`
- 🇮🇩 `chaos.Config` dari env, `Apply(ctx)` deterministik 0/100%, thread-safe, respect `ctx.Done()`; handler refactor + chaos sebelum sukses → 500 `{status:failed}`; `chaos.FromEnv()` di main

**Completed (F6) / Selesai (F6):**
- 🇬🇧 Order: `checkout_test.go` (mock clients, 6 scenarios) + `client_test.go` (httptest); Inventory: handler test (mock store) + `product_test.go` (sqlmock) + `TestOpen_BadEndpoint`; Payment: `chaos_test.go` + `payment_test.go`; new dep `DATA-DOG/go-sqlmock`
- 🇮🇩 Order: test checkout (mock clients, 6 skenario) + client (httptest); Inventory: test handler (mock store) + db (sqlmock) + Open gagal; Payment: test chaos + handler; dep baru `go-sqlmock`

**Metrics & Verification / Metrik & Verifikasi:**
- 🇬🇧 `go test ./...` green ×3; coverage `./internal/...`: order client 82.1% / handler 90.9%, inventory db 90.0% / handler 82.6%, payment chaos 94.7% / handler 100% — all ≥70% ✓; build+import+restart; `test.sh` 9 passed; chaos in-cluster: `ERROR_PERCENT=100` → `{"error":"payment failed"}`, revert → `paid` ✓
- 🇮🇩 Test hijau di 3 service; coverage semua ≥70%; build+restart payment; `test.sh` 9 passed; chaos in-cluster terverifikasi + revert → paid ✓

**Obstacles & Fixes / Kendala & Solusi:**
- 🇬🇧 Duplicate `wantStatus` field in payment test struct → renamed `wantCode`/`wantBody`
- 🇮🇩 Field duplikat `wantStatus` di struct test payment → rename `wantCode`/`wantBody`

**Next / Selanjutnya:**
- 🇬🇧 F7 OTel SDK (`pkg/telemetry` init + sampling) → F8 HTTP instrumentation & context propagation.
- 🇮🇩 F7 OTel SDK (`pkg/telemetry` init + sampling) → F8 HTTP instrumentation & context propagation.

---

### Day 5 — 2026-08-30 (OpenTelemetry SDK + HTTP Instrumentation)

**Goal / Tujuan:**
- 🇬🇧 First distributed trace: one trace ID across Order → Inventory → Payment in Jaeger.
- 🇮🇩 Trace terdistribusi pertama: satu trace ID lintas Order → Inventory → Payment di Jaeger.

**Completed (F7) / Selesai (F7):**
- 🇬🇧 `pkg/telemetry/telemetry.go`: `Init(ctx, name, version)` — Resource (service.name/version, deployment.environment.name=lab), OTLP gRPC exporter (`OTEL_EXPORTER_OTLP_ENDPOINT`, default `otel-collector:14317`, insecure), `ParentBased(TraceIDRatioBased)` from `OTEL_TRACES_SAMPLER_ARG` (default 1.0), composite `TraceContext+Baggage` propagator, returns `tp.Shutdown`; OTel deps + wiring via `replace` directives (Docker-safe); `telemetry.Init` + `defer shutdown` in all three mains
- 🇮🇩 `Init` lengkap: Resource, OTLP gRPC exporter, sampler env-driven, propagator composite, return `Shutdown`; wiring via `replace` directive; Init + defer shutdown di 3 `main.go`

**Completed (F8) / Selesai (F8):**
- 🇬🇧 `otelhttp.NewHandler(mux, "<service>")` server middleware ×3; `otelhttp.NewTransport` in both order clients; `r.Context()` threading (since Day 3); env `OTEL_EXPORTER_OTLP_ENDPOINT` + `OTEL_TRACES_SAMPLER_ARG` in deployments
- 🇮🇩 otelhttp server middleware di 3 service; `otelhttp.NewTransport` di order clients; env OTel di manifests

**Build changes / Perubahan build:**
- 🇬🇧 Dockerfiles ×3: build context moved to repo root (`COPY pkg/telemetry + services/<svc>`) because `replace` needs relative paths; `scripts/build.sh`: `podman build -f services/<svc>/Dockerfile .`
- 🇮🇩 Dockerfile x3: context pindah ke repo root; `build.sh` pakai context root

**Metrics & Verification / Metrik & Verifikasi:**
- 🇬🇧 build+vet+fmt OK; unit tests green; `test.sh` 9 passed (2 initial flakes = chaos 10% error, by design); **Jaeger: one trace (`d87f3b8f...`) with 3 services** — `POST /checkout` (2079ms) → `GET /inventory/{id}` (15ms) + `POST /pay` (2014ms, chaos delay visible); W3C traceparent proven
- 🇮🇩 Semua build/test hijau; **satu trace 3 service di Jaeger**, chaos delay terlihat di waterfall; propagasi W3C terbukti

**Obstacles & Fixes / Kendala & Solusi:**
- 🇬🇧 Local `replace` breaks in the old per-module Docker context → Dockerfiles now build from repo root
- 🇮🇩 `replace` lokal tidak jalan di build context lama (per module) → Dockerfile build context di-root repo
- 🇬🇧 Jaeger `limit=1` polluted by `/health` probe traces → filter by `operation`
- 🇮🇩 Query Jaeger kena spam trace `/health` dari probe → filter `operation`

**Next / Selanjutnya:**
- 🇬🇧 F9 DB instrumentation (otelsql) → F10 manual spans → F11 attributes/events/baggage.
- 🇮🇩 F9 DB instrumentation (otelsql) → F10 manual spans → F11 attributes/events/baggage.

---

### Day 6 — 2026-08-31 (DB Spans, Manual Spans & Baggage)

**Goal / Tujuan:**
- 🇬🇧 Richer traces: DB spans, 3 manual business spans, attributes/events/status, baggage `order.id` across services.
- 🇮🇩 Trace makin kaya: DB span, 3 manual business spans, attributes/events/status, baggage `order.id` lintas service.

**Completed (F9) / Selesai (F9):**
- 🇬🇧 `db.Open` → `otelsql.Open("pgx", ...)` + `db.system=postgresql` (dep `XSAM/otelsql` v0.43.0); DB spans visible: `sql.conn.query`, `sql.rows`, `sql.conn.reset_session`
- 🇮🇩 `otelsql.Open` + attr `db.system=postgresql`; DB span muncul di bawah inventory

**Completed (F10) / Selesai (F10):**
- 🇬🇧 Manual spans in checkout: `validate-order` (order.item_id, order.quantity), `check-inventory` (product.id, product.stock), `process-payment` (payment.order_id, payment.amount, payment.status)
- 🇮🇩 3 manual span di checkout + attrs masing-masing

**Completed (F11) / Selesai (F11):**
- 🇬🇧 Handler attrs (`product.*`, `payment.*`); events `payment_started`/`payment_completed`/`payment_failed`; chaos error → span ERROR + RecordError (F11.5); baggage `order.id` set early in Order, read in Inventory & Payment (F11.6); parent checkout span marked ERROR on downstream failure (F11.7)
- 🇮🇩 Attrs handler, events payment, status ERROR + RecordError, baggage `order.id` dibaca 2 downstream, parent span ikut ERROR

**Metrics & Verification / Metrik & Verifikasi:**
- 🇬🇧 build+vet+fmt+test green; deploy + `test.sh` 9 passed; Jaeger trace `54a6c09c...`: manual spans + attrs + events + baggage + DB spans complete; chaos 100% trace `54c7f43b...`: `POST /pay`, `process-payment`, parent checkout all `error=True` with `payment_failed`; revert → paid ✓
- 🇮🇩 Semua test hijau; trace enrichment lengkap di Jaeger; chaos ERROR terverifikasi lalu revert → paid ✓

**Obstacles & Fixes / Kendala & Solusi:**
- 🇬🇧 After build+apply pods kept old images (manifest unchanged) → `kubectl rollout restart` ×3
- 🇮🇩 Setelah build+apply pods tidak otomatis pakai image baru → `rollout restart` ketiga deployment
- 🇬🇧 Jaeger tag `value` may be typed object or plain string → normalize in verification parser
- 🇮🇩 Parser tag Jaeger: value bisa object atau plain string → normalisasi

**Next / Selanjutnya:**
- 🇬🇧 F12 quality scripts → F13 integration & E2E tests → F14 docs.
- 🇮🇩 F12 lengkapi quality scripts → F13 integration & E2E tests → F14 docs.

---

### Day 7 — 2026-09-05 (Finish Line: Quality, Tests, Docs & DoD)

**Goal / Tujuan:**
- 🇬🇧 Close the project — F12+F13+F14+F15 in one day; DoD PRD §45 fully passed.
- 🇮🇩 Menutup project — F12+F13+F14+F15 sekaligus, DoD PRD §45 lulus semua.

**Completed (F12) / Selesai (F12):**
- 🇬🇧 `scripts/check.sh`: gofmt (fail if dirty) + go vet + golangci-lint + tests + coverage gate ≥70% per `./internal/...` package; `scripts/chaos-test.sh`: forced error 100% (500) + delay 100% (~2s measured) + revert; non-zero exit on failure
- 🇮🇩 `check.sh` lengkap dengan coverage gate; `chaos-test.sh` assertion error + delay + revert

**Completed (F13) / Selesai (F13):**
- 🇬🇧 `checkout_integration_test.go`: real `CheckoutHandler` + real HTTP clients → httptest downstreams (success / not found → 500 / insufficient → 400)
- 🇮🇩 Integration test full-stack: handler + client asli → httptest downstreams

**Completed (F14) / Selesai (F14):**
- 🇬🇧 `docs/tracing-examples.md` (3 scenarios per PRD §46), `docs/experiments.md` (6 experiments per PRD §43 — steps + expected + lesson), final README (architecture, structure, prerequisites, ports, quickstart, env config, teardown)
- 🇮🇩 Docs tracing examples + 6 eksperimen + README final

**Metrics & Verification (DoD PRD §45 — 15/15 passed) / Metrik & Verifikasi (DoD — 15/15 lulus):**
- 🇬🇧 `check.sh` ALL PASSED (coverage 82.1–96.7%); `chaos-test.sh` PASSED (delay 2034ms); `test.sh` 9 passed; sampling 0.1: 20 checkouts → only 3 new traces; collector down (replicas=0): 3× checkout still **200** — telemetry fail-open (D10); full enrichment trace verified Day 6; rebuild-from-repo via build+deploy
- 🇮🇩 `check.sh` hijau semua; `chaos-test.sh` lulus; `test.sh` 9 passed; sampling 0.1 → hanya 3 trace baru; collector down → checkout tetap 200; trace enrichment lengkap; rebuild jalan

**Obstacles & Fixes / Kendala & Solusi:**
- 🇬🇧 None major — scripts passed on first/second run
- 🇮🇩 Tidak ada besar — semua script lulus di run pertama/kedua

**Status:** 🇬🇧 **PROJECT DONE** — all 15 phases of implementation_plan.md complete. Remaining optional: publishing daily blogs, CI runner (outside PRD scope).
🇮🇩 **STATUS: PROJECT DONE** — 15/15 fase implementation_plan.md selesai. Sisa opsional: publikasi blog harian, CI runner (di luar PRD).