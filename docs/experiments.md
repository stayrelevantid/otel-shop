# Dokumentasi Eksperimen — OTel-Shop Lab

Enam eksperimen inti lab (PRD §43). Setiap eksperimen: tujuan, langkah,
dan hasil yang diharapkan. Status cluster: `./scripts/cluster.sh` +
`./scripts/deploy.sh` sudah dijalankan, Jaeger di `http://localhost:16686`.

---

## Eksperimen 1 — Propagation ON vs OFF

**Tujuan:** membuktikan trace id tidak pindah antar service tanpa propagator.

**Langkah:**
1. ON: jalankan `POST /checkout`, ambil trace via Jaeger API
   (`?service=order-service&operation=POST /checkout`).
2. OFF: ubah propagator di `pkg/telemetry` — ganti composite menjadi hanya
   `propagation.TraceContext{}` (baggage off) **atau** buang
   `otelhttp.NewTransport` di order clients → rebuild
   (`./scripts/build.sh && kubectl rollout restart -n otel-shop
   deployment/{order,inventory,payment}-service`).
3. Checkout lagi, bandingkan.

**Hasil yang diharapkan:**
- ON: satu trace ID memuat order + inventory + payment (Day 5 bukti: `d87f3b8f...`).
- OFF: setiap service jadi root trace sendiri — 3 trace terpisah per request.
  Baggage `order.id` juga hilang jika propagator baggage dilepas.

**Pelajaran:** context propagation adalah "lem" distributed tracing;
tanpanya setiap service buta terhadap request yang sama.

---

## Eksperimen 2 — Payment Bottleneck

**Tujuan:** melihat latensi artifisial di waterfall.

**Langkah:**
1. `kubectl -n otel-shop set env deploy/payment-service PAYMENT_DELAY_PERCENT=100`
2. Tunggu rollout, jalankan 1 checkout.
3. Buka trace di Jaeger; revert env (`PAYMENT_DELAY_PERCENT=20`).

**Hasil yang diharapkan:** span `POST /pay` ≈ 2000ms; client span
`HTTP POST` di order memanjang mengapit; parent checkout ikut panjang.
Contoh terukur (Day 7): 2034ms.

**Pelajaran:** satu hop lambat langsung terlihat hop-nya dan durasinya —
bukan sekadar "request lama".

---

## Eksperimen 3 — DB Bottleneck

**Tujuan:** mengamati query database sebagai child span.

**Langkah:**
1. Buka trace checkout terbaru, pilih span `GET /inventory/{id}` di
   inventory-service.
2. Perhatikan child `sql.conn.query` / `sql.rows` (otelsql).

**Simulasi bottleneck (opsional):** tambah data besar ke tabel produk
(`INSERT ... generate_series`) atau tingkatkan latency network (pg_sleep
di sisi lain), lalu bandingkan durasi `sql.conn.query`.

**Hasil yang diharapkan:** query muncul sebagai child span dengan
`db.system=postgresql`; setiap percepatan/pelambatan DB terlihat presisi
di level query, bukan hanya "inventory lambat".

---

## Eksperimen 4 — Error Chain

**Tujuan:** menelusuri kegagalan dari child sampai root dalam satu trace.

**Langkah:**
1. `kubectl -n otel-shop set env deploy/payment-service PAYMENT_ERROR_PERCENT=100`
2. Checkout → ambil trace (lihat `docs/tracing-examples.md` #3).
3. Revert (`PAYMENT_ERROR_PERCENT=10`).

**Hasil yang diharapkan:** `POST /pay` ERROR + event `payment_failed` +
exception; `process-payment`, `HTTP POST`, dan parent `POST /checkout`
semuanya ERROR. Response: 500 `{"error":"payment failed"}`.
Bukti Day 6: traceID `54c7f43b...`.

**Pelajaran:** error status di child + parent membuat root-cause dapat
ditemukan dalam sekali klik, tanpa correlation log manual.

---

## Eksperimen 5 — Sampling 100% vs 10%

**Tujuan:** melihat efek sampler terhadap volume trace.

**Langkah:**
1. Catat jumlah trace checkout saat ini (Jaeger API, operation filter).
2. `kubectl -n otel-shop set env deploy/order-service OTEL_TRACES_SAMPLER_ARG=0.1`
   (order = pembuat root trace; sampler parent-based menjaga keputusan ini).
3. Tunggu rollout, kirim ~20 checkout, hitung berapa yang muncul di Jaeger.
4. Revert ke `1.0`.

**Hasil yang diharapkan:** dengan 0.1 hanya ~10% request (±2 dari 20)
yang menghasilkan trace; proporsi jauh lebih kecil dibanding sebelumnya.
Trace yang lolos tetap utuh (parent-based: child ikut keputusan root).

**Pelajaran:** sampling mengontrol biaya telemetry; root decision
dipropagasikan sehingga trace yang tersampling tetap konsisten.

---

## Eksperimen 6 — Collector Failure

**Tujuan:** memastikan observability tidak pernah menahan business logic
(NFR Reliability, D10).

**Langkah:**
1. `kubectl -n otel-shop scale deploy/otel-collector --replicas=0`
2. Jalankan checkout berulang → semua tetap `200 paid`.
3. Scale balik (`--replicas=1`), tunggu rollout.
4. Trace yang dikirim saat collector down hilang (batcher drop setelah
   retry); trace berikutnya normal kembali.

**Hasil yang diharapkan:** checkout tetap sukses tanpa collector;
tidak ada timeout/error di response bisnis. Hanya data telemetry sementara
yang tidak terkirim.

**Pelajaran:** telemetry harus fail-open. Batching + timeout non-blocking
membuat observer tidak pernah jadi single point of failure untuk request.
