# Experiments / Dokumentasi Eksperimen — OTel-Shop Lab

**🇬🇧 EN** — Six core lab experiments (PRD §43). Each: goal, steps, expected result. Prerequisite: `./scripts/cluster.sh` + `./scripts/deploy.sh` done, Jaeger at `http://localhost:16686`.

**🇮🇩 ID** — Enam eksperimen inti lab (PRD §43). Setiap eksperimen: tujuan, langkah, dan hasil yang diharapkan. Prasyarat: `./scripts/cluster.sh` + `./scripts/deploy.sh` sudah jalan, Jaeger di `http://localhost:16686`.

---

## 1. Propagation ON vs OFF / Propagasi ON vs OFF

**🇬🇧 EN — Goal:** prove the trace id does not move between services without a propagator.
**Steps:** 1) ON: run a checkout, fetch the trace via Jaeger API (`?service=order-service&operation=POST /checkout`). 2) OFF: change the composite propagator in `pkg/telemetry` (drop `Baggage{}` or remove `otelhttp.NewTransport` in order clients) → rebuild (`./scripts/build.sh` + `rollout restart`). 3) Checkout again and compare.
**Expected:** ON: one trace ID holds order+inventory+payment (Day 5 proof: `d87f3b8f...`). OFF: each service becomes its own root trace — 3 separate traces per request; baggage `order.id` also lost.
**Lesson:** context propagation is the "glue" of distributed tracing; without it each service is blind to the same request.

**🇮🇩 ID — Tujuan:** membuktikan trace id tidak pindah antar service tanpa propagator.
**Langkah:** 1) ON: checkout, ambil trace via Jaeger API. 2) OFF: ubah propagator composite di `pkg/telemetry` (lepas `Baggage{}` atau buang `otelhttp.NewTransport` di order clients) → rebuild + rollout restart. 3) Checkout lagi, bandingkan.
**Hasil yang diharapkan:** ON: satu trace ID memuat 3 service. OFF: tiap service jadi root trace sendiri — 3 trace terpisah; baggage `order.id` juga hilang.
**Pelajaran:** context propagation adalah "lem" distributed tracing; tanpanya tiap service buta terhadap request yang sama.

---

## 2. Payment Bottleneck / Kemacetan Payment

**🇬🇧 EN — Goal:** see artificial latency in the waterfall.
**Steps:** 1) `kubectl -n otel-shop set env deploy/payment-service PAYMENT_DELAY_PERCENT=100`. 2) Wait for rollout, run 1 checkout. 3) Open the trace in Jaeger; revert env.
**Expected:** span `POST /pay` ≈ 2000ms; client span `HTTP POST` in order stretches around it; parent checkout ≈ delay + rest. Measured (Day 7): 2034ms.
**Lesson:** one slow hop is instantly visible — which hop and how long, not just "request slow".

**🇮🇩 ID — Tujuan:** melihat latensi artifisial di waterfall.
**Langkah:** 1) Set env delay 100%. 2) Tunggu rollout, checkout 1×. 3) Buka trace di Jaeger; revert env.
**Hasil yang diharapkan:** span `POST /pay` ≈ 2000ms; `HTTP POST` (order) memanjang; parent checkout ikut panjang. Terukur (Day 7): 2034ms.
**Pelajaran:** satu hop lambat langsung terlihat hop-nya dan durasinya — bukan sekadar "request lama".

---

## 3. DB Bottleneck / Kemacetan Database

**🇬🇧 EN — Goal:** observe the database query as a child span.
**Steps:** 1) Open the latest checkout trace, select `GET /inventory/{id}` in inventory-service. 2) Notice children `sql.conn.query` / `sql.rows` (otelsql). Optional simulation: grow the products table (large `INSERT ... generate_series`) or add DB-side latency (`pg_sleep`), then compare `sql.conn.query` durations.
**Expected:** the query appears as a child span with `db.system=postgresql`; DB speed changes are visible at query level, not just "inventory slow".

**🇮🇩 ID — Tujuan:** mengamati query database sebagai child span.
**Langkah:** 1) Buka trace checkout terbaru, pilih span `GET /inventory/{id}`. 2) Perhatikan child `sql.conn.query` / `sql.rows`. Simulasi opsional: besar tabel produk atau `pg_sleep`, bandingkan durasi `sql.conn.query`.
**Hasil yang diharapkan:** query muncul sebagai child span dengan `db.system=postgresql`; perubahan kecepatan DB terlihat di level query, bukan hanya "inventory lambat".

---

## 4. Error Chain / Rantai Error

**🇬🇧 EN — Goal:** trace a failure from child to root in a single trace.
**Steps:** 1) `PAYMENT_ERROR_PERCENT=100`. 2) Checkout → fetch the trace (see `tracing-examples.md` #3). 3) Revert to `10`.
**Expected:** `POST /pay` ERROR + `payment_failed` event + exception; `process-payment`, `HTTP POST`, and parent `POST /checkout` all ERROR. Response: 500 `{"error":"payment failed"}`. Day 6 proof: traceID `54c7f43b...`.
**Lesson:** error status on child + parent makes root-cause a one-click find — no manual log correlation.

**🇮🇩 ID — Tujuan:** menelusuri kegagalan dari child sampai root dalam satu trace.
**Langkah:** 1) Set `PAYMENT_ERROR_PERCENT=100`. 2) Checkout → ambil trace (lihat `tracing-examples.md` #3). 3) Revert ke `10`.
**Hasil yang diharapkan:** `POST /pay` ERROR + event + exception; `process-payment`, `HTTP POST`, parent semuanya ERROR; response 500. Bukti Day 6: traceID `54c7f43b...`.
**Pelajaran:** error status di child + parent membuat root-cause ditemukan sekali klik — tanpa correlation log manual.

---

## 5. Sampling 100% vs 10% / Sampling 100% vs 10%

**🇬🇧 EN — Goal:** see the sampler's effect on trace volume.
**Steps:** 1) Note current checkout trace count (Jaeger API, operation filter). 2) `kubectl -n otel-shop set env deploy/order-service OTEL_TRACES_SAMPLER_ARG=0.1` (order creates the root trace; ParentBased keeps that decision). 3) Rollout, send ~20 checkouts, count new traces. 4) Revert to `1.0`.
**Expected:** at 0.1 only ~10% of requests (~2 of 20) produce a trace; surviving traces stay complete (ParentBased). Day 7 proof: 3 of 20.
**Lesson:** sampling controls telemetry cost; the root decision propagates so kept traces stay consistent.

**🇮🇩 ID — Tujuan:** melihat efek sampler terhadap volume trace.
**Langkah:** 1) Catat jumlah trace checkout saat ini. 2) Set `OTEL_TRACES_SAMPLER_ARG=0.1` di order-service (order = pembuat root; ParentBased menjaga keputusan). 3) Rollout, kirim ~20 checkout, hitung. 4) Revert ke `1.0`.
**Hasil yang diharapkan:** 0.1 → hanya ~10% (±2 dari 20) menghasilkan trace; trace yang lolos tetap utuh. Bukti Day 7: 3 dari 20.
**Pelajaran:** sampling mengontrol biaya telemetry; keputusan root dipropagasikan sehingga trace yang tersampling konsisten.

---

## 6. Collector Failure / Collector Mati

**🇬🇧 EN — Goal:** prove observability never holds business logic hostage (NFR Reliability, D10).
**Steps:** 1) `kubectl -n otel-shop scale deploy/otel-collector --replicas=0`. 2) Run repeated checkouts → all 200 paid. 3) Scale back (`--replicas=1`), wait for rollout. 4) Traces sent while down are lost (batcher drops after retries); new ones flow again.
**Expected:** checkout keeps succeeding without the collector; no timeouts/errors in business responses. Day 7 proof: 3× checkout 200 while down; traces resume after recovery.
**Lesson:** telemetry must fail-open — batching + non-blocking timeouts keep the observer from ever becoming a single point of failure for requests.

**🇮🇩 ID — Tujuan:** memastikan observability tidak pernah menahan business logic (NFR Reliability, D10).
**Langkah:** 1) Scale collector ke 0. 2) Checkout berulang → semua 200 paid. 3) Scale balik ke 1, tunggu rollout. 4) Trace saat collector down hilang (batcher drop setelah retry); trace baru normal kembali.
**Hasil yang diharapkan:** checkout tetap sukses tanpa collector; tidak ada timeout/error di response bisnis. Bukti Day 7: 3× checkout 200 saat down; trace mengalir lagi setelah naik.
**Pelajaran:** telemetry harus fail-open — batching + timeout non-blocking membuat observer tidak pernah jadi single point of failure untuk request.