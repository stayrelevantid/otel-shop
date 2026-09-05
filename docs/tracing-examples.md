# Contoh Trace — OTel-Shop Lab

Referensi membandingkan hasil di Jaeger (`http://localhost:16686`) dengan
bentuk waterfall yang diharapkan (PRD §46). Filter: service `order-service`,
operation `POST /checkout` (trace `/health` dari probe akan membanjiri daftar).

---

## 1. Checkout Normal

Skenario: `POST /checkout {"item_id":"A123","qty":1}` dengan chaos default
(delay 20%, error 10%) — asumsikan chaos tidak ke-trigger.

**Struktur yang diharapkan:**

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

**Ciri:** semua span status normal (tidak ada `error=true`), satu trace ID
berisi 3 service, `baggage.order.id` sama di inventory & payment.

---

## 2. Checkout Lambat (Chaos Delay)

Skenario: `PAYMENT_DELAY_PERCENT=100` + `PAYMENT_DELAY_MS=2000`
(atau pakai `./scripts/chaos-test.sh` scenario 2).

**Ciri:**
- `POST /pay` di payment-service stretch ~2000ms.
- `HTTP POST` (client span order) dan `process-payment` juga memanjang
  (mengapit child).
- Parent `POST /checkout` total ≈ durasi delay + sisanya.

**Analisis yang bisa dilakukan:** instant terlihat bottleneck ada di hop
payment, bukan di inventory atau DB — tanpa perlu logging atau profiling.

---

## 3. Error Chain (Chaos Error)

Skenario: `PAYMENT_ERROR_PERCENT=100` (chaos-test.sh scenario 1).

**Ciri:**
- `POST /pay` → status **ERROR**, event `payment_failed`, exception record.
- `HTTP POST` di order → status ERROR.
- `process-payment` → status ERROR + `payment.status=failed`.
- Parent `POST /checkout` → ikut **ERROR** (F11.7).

Response body client: `{"error":"payment failed"}` (500). Error chain bisa
ditelusuri dari child sampai root dalam satu trace — tanpa correlation log
manual.
