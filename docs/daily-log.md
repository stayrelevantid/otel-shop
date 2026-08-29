# Daily Log — OTel-Shop Lab

Catatan progres harian challenge. Satu entri per hari kerja, format konsisten di bagian bawah.

Repo: [github.com/stayrelevantid/otel-shop](https://github.com/stayrelevantid/otel-shop) ·
Status fase: [progress-tracker.md](progress-tracker.md)

---

## Template

```markdown
### Day N — YYYY-MM-DD (tema singkat)

**Tujuan:** <apa yang ingin dicapai hari ini>

**Selesai:**
- <item> 
- <item>

**Metrik/Verifikasi:**
- <hasil test/verifikasi konkret>

**Kendala & Solusi:**
- <masalah → penyelesaian>

**Next:** <fokus esok hari>
```

---

## Entri Harian

### Day 1 — 2026-08-26 (Infra-first)

**Tujuan:** Fondasi infra + workspace Go yang bisa diverifikasi sebelum koding serius.

**Selesai:**
- Git init, `.gitignore`, scan gitleaks (clean), initial commit
- Go workspace 4 module (`order`, `inventory`, `payment`, `telemetry`) + `go.work`
- `.golangci.yml` v2, `db/init.sql` (schema + seed A123/B123/C123)
- Manifest K8s: namespace, Postgres 18 (NodePort 15432), OTel Collector 0.159.0 (OTLP → batch → OTLP ke jaeger:4317), Jaeger all-in-one 1.76.0 (UI 16686)
- `scripts/cluster.sh`, `scripts/deploy.sh`
- Cluster k3d `otel-shop` hidup, semua pod Running
- Smoke trace OTLP/HTTP → Collector → terlihat di Jaeger API

**Metrik/Verifikasi:**
- Jaeger UI: HTTP 200 di localhost:16686
- Seed data: 3 produk terbaca via psql exec
- Smoke trace: ditemukan via `/api/traces?service=smoke`

**Kendala & Solusi:**
- k3d menolak Podman (wajib Docker daemon) → nyalakan Docker Desktop
- Image pull lambat (collector >9m) menyebabkan timeout deploy → jalankan ulang; rollout sukses
- `go build ./...` dari root gagal (multi-module) → build per module dengan `-o /dev/null`

**Next:** F2 penuh — model, handler, routing, Dockerfile, deploy 3 service.

---

### Day 2 — 2026-08-27 (Application Skeleton + Deploy)

**Tujuan:** Tiga service naik ke k3d dengan endpoint dasar tervalidasi dari host.

**Selesai:**
- Model: `CheckoutRequest/Response/ErrorResponse`, `InventoryResponse`, `PayRequest/Response`
- Handler + routing stdlib ServeMux (Go 1.22+): `GET /health`, `POST /checkout`, `GET /inventory/{id}`, `POST /pay` (logika stub sesuai fase)
- Dockerfile multi-stage x3 (golang:1.27-alpine → scratch, CGO off)
- Manifest `deploy/{order,inventory,payment}`: Deployment + Service NodePort 18080–18082 + probe `/health`
- `scripts/build.sh`: podman build → docker load → k3d image import
- `scripts/deploy.sh` diperluas: apply apps + rollout status per service
- Blogpost + LinkedIn post Day 2

**Metrik/Verifikasi:**
- compile+vet+gofmt OK untuk ketiga module
- Semua pod Running; rollout sukses
- E2E curl dari host: health 3×200; checkout valid `{order_id,status:paid}`; checkout invalid 400; inventory A123 `{stock:10}`; pay valid `{status:success}`; pay invalid 400

**Kendala & Solusi:**
- Tag image berubah jadi `localhost/otel-shop/...` setelah podman save/load → retag di sisi Docker sebelum `k3d image import`

**Next:** F3 DB integration (`StockStore` + pgx) dan F4 checkout flow antar-service.

---

### Day 3 — 2026-08-28 (DB Integration + Checkout Flow + E2E)

**Tujuan:** Inventory baca Postgres sungguhan, Order panggil Inventory & Payment via HTTP, alur checkout end-to-end jalan.

**Selesai (F3):**
- `internal/store/store.go`: interface `StockStore` + `ErrNotFound`
- `internal/db/db.go`: `Open(DATABASE_URL)` via pgx/v5 stdlib pool + `PingContext`
- `internal/db/product.go`: `GetStock(ctx,id)` + map `sql.ErrNoRows` → `ErrNotFound`
- Handler inventory refactor: inject store, status 200/404/500
- `main.go` inventory: buka DB, wire store, exit kalau `DATABASE_URL` kosong

**Selesai (F4):**
- `internal/client/interfaces.go`: `InventoryClient`, `PaymentClient` (interface, testable)
- `internal/client/inventory.go`: GET `/inventory/{id}`, 404 → `ErrInventory`
- `internal/client/payment.go`: POST `/pay`, non-200 → `ErrPayment`
- `internal/model/errors.go`: typed errors (`ErrInventory`, `ErrPayment`, `ErrInsufficientStock`, `ErrValidation`)
- `internal/handler/checkout.go`: `CheckoutHandler` method — validate → `GetStock` → cek stock → `Pay(qty*10)` → 200 `{order_id,status:paid}`
- Downstream gagal → **500 + `{"error":"..."}`** (sesuai keputusan)
- `main.go` order: env `INVENTORY_URL`/`PAYMENT_URL` (default localhost)

**Selesai (test):**
- `scripts/test.sh` (F12.4 dipercepat): health ×3 + inventory + 4 skenario checkout, exit non-zero bila gagal

**Metrik/Verifikasi:**
- Lokal E2E (port terpisah hindari tabrakan NodePort cluster): A123→paid, UNKNOWN→500 inventory failed, C123 qty999→400 insufficient, invalid→400
- Build image → `k3d image import` → `rollout restart` order & inventory → E2E via NodePort host: sama persis hijau
- `scripts/test.sh`: **9 passed, 0 failed**

**Kendala & Solusi:**
- Local test pertama salah karena port 18080/18081/18082 sudah dipakai NodePort cluster → test lokal pakai port lain; deploy verifikasi lewat NodePort cluster

**Next:** F5 Chaos (payment delay/error percent) + F6 unit tests (mock client/store, coverage ≥70%).

---

### Day 4 — 2026-08-29 (Chaos Engineering + Unit Tests)

**Tujuan:** Payment punya chaos (delay/error %), semua business logic ter-cover unit test ≥70%.

**Selesai (F5):**
- `internal/chaos/chaos.go`: `Config{DelayPercent,DelayMS,ErrorPercent}` dari env, `Apply(ctx)` deterministik saat 0/100%, thread-safe rng, respect `ctx.Done()`
- Payment handler refactor: `Handler` struct, chaos dipanggil sebelum sukses → 500 `{status:failed}`
- `payment/main.go`: `chaos.FromEnv()`

**Selesai (F6):**
- Order: `checkout_test.go` (mock `InventoryClient`+`PaymentClient`, 6 skenario) + `client_test.go` (httptest: inventory 200/404, payment valid/invalid/500)
- Inventory: `inventory_test.go` (mock `StockStore`: exists/404/db error) + `product_test.go` (sqlmock: found/`ErrNoRows`→`ErrNotFound`/db error) + `TestOpen_BadEndpoint`
- Payment: `chaos_test.go` (0/100% deterministik, delay timing, ctx cancel, `FromEnv`+fallback) + `payment_test.go` (sukses/chaos 500/invalid)
- Dependency baru: `DATA-DOG/go-sqlmock` (inventory)

**Metrik/Verifikasi:**
- `go test ./...` hijau di 3 service
- Coverage `./internal/...`: order client 82.1% / handler 90.9%, inventory db 90.0% / handler 82.6%, payment chaos 94.7% / handler 100% — **semua ≥70%** ✓
- Build + import + rollout restart payment; `scripts/test.sh`: 9 passed, 0 failed
- Chaos in-cluster: `PAYMENT_ERROR_PERCENT=100` → checkout `{"error":"payment failed"}`; revert 10% → `paid` ✓

**Kendala & Solusi:**
- Field duplikat `wantStatus` di struct test payment → rename `wantCode`/`wantBody`

**Next:** F7 OTel SDK (`pkg/telemetry` init + sampling) → F8 HTTP instrumentation & context propagation.
