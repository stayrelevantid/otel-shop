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
