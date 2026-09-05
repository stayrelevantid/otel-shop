# Trace Examples / Contoh Trace — OTel-Shop Lab

**🇬🇧 EN** — Reference for comparing what you see in Jaeger (`http://localhost:16686`) with the expected waterfall shapes (PRD §46). Filter: service `order-service`, operation `POST /checkout` (probe `/health` traces will flood the list otherwise).

**🇮🇩 ID** — Referensi membandingkan hasil di Jaeger (`http://localhost:16686`) dengan bentuk waterfall yang diharapkan (PRD §46). Filter: service `order-service`, operation `POST /checkout` (trace `/health` dari probe akan membanjiri daftar).

---

## 1. Normal Checkout / Checkout Normal

**🇬🇧 EN — Scenario:** `POST /checkout {"item_id":"A123","qty":1}` with default chaos (delay 20%, error 10%) — assuming chaos does not trigger.

**Expected structure / Struktur yang diharapkan:**

```
order-service        POST /checkout           (parent, ~10-30ms)
├─ validate-order    (0ms)                    order.item_id, order.quantity
├─ check-inventory   (~1-8ms)                 product.id, product.stock
│  └─ HTTP GET       (client span)
│     └─ inventory-service GET /inventory/{id}
│        ├─ sql.conn.query   (DB span, otelsql)
│        └─ sql.rows
└─ process-payment   (~1-2ms)                 payment.order_id, amount, status
   └─ HTTP POST      (client span)
      └─ payment-service POST /pay            events: payment_started → payment_completed
```

**🇬🇧 EN — Markers:** all spans normal (no `error=true`), one trace ID holds 3 services, same `baggage.order.id` in inventory & payment.

**🇮🇩 ID — Skenario:** `POST /checkout` dengan chaos default — asumsikan chaos tidak ke-trigger.

**🇮🇩 ID — Ciri:** semua span normal (tidak ada `error=true`), satu trace ID berisi 3 service, `baggage.order.id` sama di inventory & payment.

---

## 2. Slow Checkout (Chaos Delay) / Checkout Lambat (Chaos Delay)

**🇬🇧 EN — Scenario:** `PAYMENT_DELAY_PERCENT=100` + `PAYMENT_DELAY_MS=2000` (or use scenario 2 of `./scripts/chaos-test.sh`).

**🇬🇧 EN — Markers:** `POST /pay` stretches to ~2000ms; the order-side `HTTP POST` and `process-payment` spans stretch to wrap it; parent `POST /checkout` total ≈ delay + rest. Analysis: the bottleneck hop is visible instantly — no logging or profiling needed.

**🇮🇩 ID — Skenario:** `PAYMENT_DELAY_PERCENT=100` + `PAYMENT_DELAY_MS=2000` (atau scenario 2 `chaos-test.sh`).

**🇮🇩 ID — Ciri:** `POST /pay` stretch ~2000ms; span `HTTP POST` (order) dan `process-payment` ikut memanjang; parent total ≈ delay + sisanya. Analisis: bottleneck terlihat instan di hop payment, bukan di inventory/DB.

---

## 3. Error Chain (Chaos Error) / Rantai Error (Chaos Error)

**🇬🇧 EN — Scenario:** `PAYMENT_ERROR_PERCENT=100` (chaos-test.sh scenario 1).

**🇬🇧 EN — Markers:** `POST /pay` → status **ERROR** with `payment_failed` event + recorded exception; `HTTP POST` in order → ERROR; `process-payment` → ERROR + `payment.status=failed`; parent `POST /checkout` → also **ERROR** (F11.7). Client response: `{"error":"payment failed"}` (500). The failure chain is traceable from child to root in one trace — no manual log correlation.

**🇮🇩 ID — Skenario:** `PAYMENT_ERROR_PERCENT=100` (scenario 1 chaos-test.sh).

**🇮🇩 ID — Ciri:** `POST /pay` → status **ERROR**, event `payment_failed`, exception record; `HTTP POST` di order → ERROR; `process-payment` → ERROR + `payment.status=failed`; parent `POST /checkout` → ikut **ERROR** (F11.7). Response: 500 `{"error":"payment failed"}`. Rantai error bisa ditelusuri child → root dalam satu trace — tanpa correlation log manual.