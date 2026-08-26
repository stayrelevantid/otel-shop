# Product Requirements Document (PRD)

## Project OTel-Shop Lab — OpenTelemetry Distributed Tracing

**Version:** 1.1
**Status:** Final
**Project Type:** Learning / Laboratory Project
**Primary Focus:** OpenTelemetry, Distributed Tracing, Observability, dan Software Quality
**Programming Language:** Golang
**Container Runtime:** Podman
**Kubernetes:** k3d
**Database:** PostgreSQL
**Telemetry Collector:** OpenTelemetry Collector
**Tracing Backend:** Jaeger

---

# 1. Visi Proyek

OTel-Shop Lab adalah laboratory project berbasis microservices yang dirancang untuk mempelajari **OpenTelemetry dan distributed tracing secara end-to-end**, sekaligus menerapkan praktik software engineering dasar seperti linting, unit testing, integration testing, dan end-to-end testing.

Aplikasi mensimulasikan proses checkout sederhana pada toko online:

```text
Client
  |
  v
Order Service
  |
  +----> Inventory Service ----> PostgreSQL
  |
  +----> Payment Service
```

Telemetry flow:

```text
Order Service
Inventory Service
Payment Service
        |
        | OTLP
        v
OpenTelemetry Collector
        |
        v
      Jaeger
```

Project ini sengaja memiliki business logic sederhana agar fokus pembelajaran berada pada:

```text
Application
    ↓
Instrumentation
    ↓
Context Propagation
    ↓
Distributed Trace
    ↓
OTLP
    ↓
Collector
    ↓
Jaeger
    ↓
Trace Analysis
```

---

# 2. Tujuan Pembelajaran

Setelah project selesai, developer diharapkan mampu memahami dan mengimplementasikan:

## OpenTelemetry Fundamentals

- Trace
- Span
- Trace ID
- Span ID
- Root Span
- Parent/Child Span
- Span Context
- Span Attributes
- Span Events
- Span Status
- Error Recording
- Resource
- OTLP
- Sampling

## Distributed Tracing

- HTTP server instrumentation
- HTTP client instrumentation
- Context propagation
- W3C Trace Context
- Parent-child relationship antar service
- Trace correlation
- Database tracing
- Manual instrumentation

## Advanced OpenTelemetry

- Baggage
- Span Events
- Error propagation
- Latency analysis
- Bottleneck identification
- Sampling

## OpenTelemetry Collector

- Receiver
- Processor
- Exporter
- Pipeline
- OTLP ingestion
- Collector troubleshooting

## Kubernetes

- Deployment
- Service
- ConfigMap
- Secret
- Environment Variable
- Service Discovery
- Health Check
- Pod troubleshooting

## Software Engineering

- Go unit testing
- Mocking dependency
- Code linting
- Static analysis
- Coverage
- Integration testing
- End-to-end testing
- Automated quality check

---

# 3. Scope

## 3.1 In Scope

Project mencakup:

- k3d
- Podman
- Kubernetes
- Golang
- PostgreSQL
- Order Service
- Inventory Service
- Payment Service
- OpenTelemetry SDK
- OpenTelemetry Collector
- Jaeger
- HTTP instrumentation
- HTTP client instrumentation
- Database instrumentation
- Manual span
- Span attributes
- Span events
- Span status
- Error recording
- Context propagation
- Baggage
- Sampling
- Artificial latency
- Artificial error
- OTLP
- Unit testing
- Integration testing
- End-to-end testing
- Linting
- Static analysis
- Coverage

## 3.2 Out of Scope

Tidak diperlukan:

- Frontend
- Authentication
- JWT
- User management
- Kafka
- Redis
- Service mesh
- Prometheus
- Grafana
- Centralized logging
- Real payment gateway
- Production-grade HA
- Multi-region
- Autoscaling
- Production security architecture

---

# 4. Arsitektur Sistem

## 4.1 High-Level Architecture

```text
                         +----------------+
                         |     Client     |
                         | curl / Postman |
                         +-------+--------+
                                 |
                                 | HTTP :18080
                                 v
                    +--------------------------+
                    |      Order Service       |
                    |         :18080           |
                    |                          |
                    |        Root Span         |
                    +------------+-------------+
                                 |
                  +--------------+--------------+
                  |                             |
                  | HTTP                        | HTTP
                  v                             v
        +-------------------+         +-------------------+
        | Inventory Service |         |  Payment Service  |
        |      :18081       |         |      :18082       |
        +---------+---------+         +-------------------+
                  |
                  | SQL
                  v
           +-------------+
           | PostgreSQL  |
           |   :15432    |
           +-------------+
```

Observability:

```text
Order Service --------\
Inventory Service -----+--> OTel Collector --> Jaeger
Payment Service -------/
```

---

# 5. Port Convention

Port sengaja dibuat lebih unik dan konsisten untuk kebutuhan laboratory.

| Component         |      Port | Fungsi    |
| ----------------- | --------: | --------- |
| Order Service     | **18080** | HTTP API  |
| Inventory Service | **18081** | HTTP API  |
| Payment Service   | **18082** | HTTP API  |
| PostgreSQL        | **15432** | Database  |
| OTel Collector    | **14317** | OTLP gRPC |
| OTel Collector    | **14318** | OTLP HTTP |
| Jaeger UI         | **16686** | Web UI    |

Contoh host access:

```text
http://localhost:18080
http://localhost:18081
http://localhost:18082
http://localhost:16686
```

Internal Kubernetes communication menggunakan service DNS:

```text
http://inventory-service:18081
http://payment-service:18082
```

---

# 6. Order Service

**Port:** `18080`

Order Service adalah entry point aplikasi.

## Endpoint

```http
POST /checkout
```

## Request

```json
{
  "item_id": "A123",
  "qty": 1
}
```

## Tugas

1. Menerima request.
2. Membuat Root Span.
3. Membuat order ID.
4. Membuat baggage.
5. Validasi request.
6. Memanggil Inventory Service.
7. Memeriksa stock.
8. Memanggil Payment Service.
9. Mengembalikan response.
10. Mencatat error apabila downstream gagal.

Expected trace:

```text
POST /checkout
│
├── validate-order
│
├── GET /inventory/A123
│
└── POST /pay
```

---

# 7. Inventory Service

**Port:** `18081`

## Endpoint

```http
GET /inventory/:id
```

Contoh:

```http
GET /inventory/A123
```

Response:

```json
{
  "stock": 10
}
```

## Tugas

1. Menerima request.
2. Mengambil product ID.
3. Membaca baggage/context.
4. Query PostgreSQL.
5. Mengembalikan stock.
6. Menghasilkan database span.

Target:

```text
GET /inventory/A123
        |
        └── PostgreSQL SELECT
```

---

# 8. Payment Service

**Port:** `18082`

## Endpoint

```http
POST /pay
```

Request:

```json
{
  "order_id": "O-123",
  "amount": 100
}
```

Success:

```json
{
  "status": "success"
}
```

Failure:

```json
{
  "status": "failed"
}
```

HTTP status:

```text
200 OK
500 Internal Server Error
```

## Chaos Configuration

```text
PAYMENT_DELAY_PERCENT=20
PAYMENT_DELAY_MS=2000
PAYMENT_ERROR_PERCENT=10
```

Behavior:

```text
20% request
    ↓
delay 2 seconds

10% request
    ↓
HTTP 500
```

Nilai tersebut digunakan sebagai default laboratory configuration dan dapat diubah untuk eksperimen.

---

# 9. API Contract

## Order

```http
POST /checkout
Content-Type: application/json
```

Request:

```json
{
  "item_id": "A123",
  "qty": 1
}
```

Success:

```http
200 OK
```

```json
{
  "order_id": "O-123",
  "status": "paid"
}
```

Failure:

```http
500 Internal Server Error
```

```json
{
  "error": "payment failed"
}
```

## Inventory

```http
GET /inventory/A123
```

Success:

```json
{
  "stock": 10
}
```

Not found:

```http
404 Not Found
```

## Payment

```http
POST /pay
Content-Type: application/json
```

Request:

```json
{
  "order_id": "O-123",
  "amount": 100
}
```

Success:

```http
200 OK
```

Failure:

```http
500 Internal Server Error
```

---

# 10. Health Endpoint

Setiap service wajib mempunyai:

```http
GET /health
```

Response:

```text
200 OK
```

Digunakan untuk:

- Kubernetes readiness probe
- Kubernetes liveness probe
- local troubleshooting
- integration test

---

# 11. PostgreSQL

Database:

```text
oteldb
```

Table:

```sql
CREATE TABLE products (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    stock INT NOT NULL
);
```

Seed data:

```sql
INSERT INTO products (id, name, stock)
VALUES
('A123', 'Keyboard', 10),
('B123', 'Mouse', 20),
('C123', 'Monitor', 5);
```

Query utama:

```sql
SELECT stock
FROM products
WHERE id = $1;
```

Expected database trace:

```text
Inventory HTTP Span
        |
        └── PostgreSQL Span
```

---

# 12. OpenTelemetry SDK

Semua service wajib menggunakan OpenTelemetry SDK.

Minimal:

```text
TracerProvider
Resource
Sampler
Propagator
OTLP Exporter
```

Resource:

```text
service.name
service.version
deployment.environment
```

Service name:

```text
order-service
inventory-service
payment-service
```

Environment:

```text
deployment.environment=lab
```

---

# 13. OTLP

Semua service wajib mengirim trace menggunakan OTLP ke Collector.

```text
Application
     |
     | OTLP
     v
OTel Collector
     |
     v
Jaeger
```

Port:

```text
14317 = OTLP gRPC
14318 = OTLP HTTP
```

Aplikasi tidak mengirim telemetry secara langsung ke Jaeger.

---

# 14. OpenTelemetry Collector

Collector merupakan intermediary telemetry pipeline.

Minimal pipeline:

```text
OTLP Receiver
      |
      v
Batch Processor
      |
      v
Jaeger Exporter
```

Flow:

```text
Order --------\
Inventory -----+--> Collector --> Jaeger
Payment -------/
```

Collector harus dapat menerima telemetry dari seluruh service.

---

# 15. HTTP Instrumentation

Incoming request wajib diinstrumentasi.

Minimal:

```text
HTTP method
HTTP route
HTTP status
duration
```

Contoh:

```text
HTTP POST /checkout

http.request.method = POST
http.route = /checkout
http.response.status_code = 200
```

---

# 16. HTTP Client Instrumentation

Order Service melakukan call ke:

```text
Inventory
Payment
```

HTTP client instrumentation wajib mempertahankan trace context.

W3C Trace Context digunakan melalui header seperti:

```text
traceparent
```

Target:

```text
Order
 ├── Inventory
 └── Payment
```

berada di dalam satu distributed trace.

---

# 17. Context Propagation

Context wajib diteruskan:

```text
Client
   |
   v
Order
   |
   +----> Inventory
   |
   +----> Payment
```

Acceptance criteria:

Satu request `/checkout` menghasilkan satu Trace ID yang sama di downstream service.

---

# 18. Manual Instrumentation

Selain automatic instrumentation, harus ada custom span.

Minimal:

```text
validate-order
check-inventory
process-payment
```

Custom span harus merepresentasikan aktivitas bisnis, bukan sekadar menduplikasi HTTP span.

---

# 19. Span Attributes

## Order

```text
order.id
order.item_id
order.quantity
```

## Inventory

```text
product.id
product.stock
```

## Payment

```text
payment.order_id
payment.amount
payment.status
```

Dilarang memasukkan secret atau data sensitif.

---

# 20. Span Events

Minimal gunakan events:

```text
payment_started
payment_completed
payment_failed
```

Contoh struktur:

```text
Payment Span
   |
   +-- payment_started
   |
   +-- payment_completed
```

Event bukan child span.

---

# 21. Span Status dan Error

Payment HTTP 500 harus menghasilkan:

```text
Span Status = ERROR
```

Error harus direkam pada span.

Contoh:

```text
POST /pay
HTTP 500
Status: ERROR
```

Order harus menangani error downstream dan mengembalikan response yang sesuai.

---

# 22. Baggage

Gunakan baggage untuk latihan metadata propagation.

Contoh:

```text
order.id=O-123
```

Flow:

```text
Order
 |
 | baggage: order.id=O-123
 |
 +----> Inventory
 |
 +----> Payment
```

Inventory dan Payment harus dapat membaca `order.id`.

Pembelajaran utama:

```text
Span Attribute ≠ Baggage
```

---

# 23. Sampling

Project harus menyediakan konfigurasi sampling.

Minimal eksperimen:

```text
100% sampling
10% sampling
```

Tujuan:

memahami dampak sampling terhadap jumlah trace yang dikirim dan observability.

Sampling bukan bagian dari business logic.

---

# 24. Chaos Engineering

Chaos hanya diterapkan pada Payment Service.

## Delay

```text
PAYMENT_DELAY_PERCENT=20
PAYMENT_DELAY_MS=2000
```

Expected:

```text
Order
 |
 └── Payment
       └── ~2 seconds
```

## Error

```text
PAYMENT_ERROR_PERCENT=10
```

Expected:

```text
Order
 |
 └── Payment
       └── ERROR / HTTP 500
```

---

# 25. Code Quality

Code quality adalah requirement wajib.

Semua source code Go harus melewati:

```text
go fmt
go vet
golangci-lint
go test
```

---

# 26. Linting

Gunakan:

```text
golangci-lint
```

Command:

```bash
golangci-lint run ./...
```

Konfigurasi disimpan di repository:

```text
.golangci.yml
```

Linting harus mendeteksi minimal:

- unused code
- suspicious constructs
- error handling issue
- static analysis issue
- common Go best practices
- formatting issue

Acceptance criteria:

```text
golangci-lint run ./...
```

berhasil tanpa error.

---

# 27. Unit Testing

Semua business logic utama wajib memiliki unit test menggunakan Go testing framework.

Command:

```bash
go test ./...
```

Unit test harus tidak bergantung pada:

- Kubernetes
- PostgreSQL asli
- OTel Collector
- Jaeger
- external network

Dependency eksternal harus di-mock atau di-isolasi.

## Order Service

Minimal:

```text
checkout success
invalid request
inventory failure
payment failure
```

## Inventory Service

Minimal:

```text
product exists
product not found
database error
```

## Payment Service

Minimal:

```text
successful payment
payment error
delay configuration
error probability configuration
```

---

# 28. Test Coverage

Coverage wajib dapat dihasilkan:

```bash
go test ./... -cover
```

Target minimum:

```text
≥ 70%
```

Target tersebut terutama berlaku untuk package yang memiliki business logic.

Coverage digunakan sebagai quality indicator, bukan satu-satunya ukuran kualitas test.

---

# 29. Integration Testing

Harus tersedia integration test untuk memastikan komunikasi antar komponen.

Minimal flow:

```text
Order
  |
  +--> Inventory
  |
  +--> Payment
```

Scenario:

```text
1. Checkout sukses
2. Inventory tidak ditemukan
3. Payment gagal
```

Integration test boleh menggunakan dependency yang dijalankan menggunakan container maupun environment Kubernetes lokal.

---

# 30. End-to-End Testing

Harus tersedia E2E test terhadap environment yang sudah dideploy.

Minimal:

```text
GET /health Order
GET /health Inventory
GET /health Payment
POST /checkout
```

Contoh command:

```bash
./scripts/test.sh
```

Expected:

```text
Order      -> 200
Inventory  -> 200
Payment    -> 200
Checkout   -> expected business response
```

---

# 31. Chaos Testing

Automated test atau test script harus dapat menguji:

## Payment Error

```text
PAYMENT_ERROR_PERCENT=100
```

Expected:

```text
POST /checkout
        |
        └── Payment HTTP 500
```

## Payment Delay

```text
PAYMENT_DELAY_PERCENT=100
PAYMENT_DELAY_MS=2000
```

Expected:

```text
Payment duration ≈ 2 seconds
```

---

# 32. Unified Quality Check

Repository harus memiliki satu command untuk menjalankan seluruh static/code quality checks.

Contoh:

```bash
./scripts/check.sh
```

Flow:

```text
go fmt
   ↓
go vet
   ↓
golangci-lint
   ↓
unit test
   ↓
coverage
```

Integration dan E2E dapat dijalankan menggunakan command terpisah karena membutuhkan environment runtime.

---

# 33. Kubernetes Architecture

Namespace:

```text
otel-shop
```

Resources:

```text
otel-shop
├── order-service
├── inventory-service
├── payment-service
├── postgres
├── otel-collector
└── jaeger
```

Service:

```text
order-service       :18080
inventory-service   :18081
payment-service     :18082
postgres            :15432
otel-collector      :14317 / :14318
jaeger              :16686
```

---

# 34. Kubernetes Service Discovery

Gunakan Kubernetes DNS.

Order:

```text
http://inventory-service:18081
```

Payment:

```text
http://payment-service:18082
```

Tidak boleh hardcode IP Pod.

---

# 35. Containerization

Semua service wajib memiliki container image.

Gunakan multi-stage build:

```text
Go Build
   ↓
Minimal Runtime Image
```

Tujuan:

- image kecil
- reproducible build
- runtime dependency minimum

---

# 36. Configuration

Konfigurasi harus menggunakan environment variable.

## Order

```text
PORT=18080
INVENTORY_URL=http://inventory-service:18081
PAYMENT_URL=http://payment-service:18082
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:14317
OTEL_SERVICE_NAME=order-service
```

## Inventory

```text
PORT=18081
DATABASE_URL=...
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:14317
OTEL_SERVICE_NAME=inventory-service
```

## Payment

```text
PORT=18082
PAYMENT_DELAY_PERCENT=20
PAYMENT_DELAY_MS=2000
PAYMENT_ERROR_PERCENT=10
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:14317
OTEL_SERVICE_NAME=payment-service
```

---

# 37. Struktur Repository

```text
otel-shop/
│
├── services/
│   ├── order/
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── Dockerfile
│   │   └── go.mod
│   │
│   ├── inventory/
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── Dockerfile
│   │   └── go.mod
│   │
│   └── payment/
│       ├── cmd/
│       ├── internal/
│       ├── Dockerfile
│       └── go.mod
│
├── deploy/
│   ├── namespace.yaml
│   ├── postgres/
│   ├── order/
│   ├── inventory/
│   ├── payment/
│   ├── otel-collector/
│   └── jaeger/
│
├── db/
│   └── init.sql
│
├── scripts/
│   ├── check.sh
│   ├── build.sh
│   ├── deploy.sh
│   └── test.sh
│
├── docs/
│   └── tracing-examples.md
│
├── .golangci.yml
├── go.work
└── README.md
```

Catatan: penggunaan `go.work` dimaksudkan untuk mengelola beberapa module Go dalam satu repository apabila masing-masing service dipisahkan menjadi module tersendiri.

---

# 38. Milestone Implementasi

## Fase 1 — Infrastruktur Dasar

Implementasi:

```text
k3d
PostgreSQL
OTel Collector
Jaeger
```

Acceptance criteria:

- Cluster aktif.
- PostgreSQL running.
- Collector running.
- Jaeger running.
- Jaeger UI dapat diakses pada `http://localhost:16686`.

---

## Fase 2 — Application Skeleton

Implementasi:

```text
order-service
inventory-service
payment-service
```

Masing-masing:

```text
HTTP server
/health
Dockerfile
```

---

## Fase 3 — Database

Implementasi:

```text
PostgreSQL
products
seed data
Inventory query
```

---

## Fase 4 — Code Quality & Unit Testing

Implementasi:

```text
go fmt
go vet
golangci-lint
unit testing
coverage
```

Acceptance criteria:

```text
go test ./...
golangci-lint run ./...
go vet ./...
```

berhasil.

Target coverage:

```text
>= 70%
```

---

## Fase 5 — Basic OpenTelemetry

Implementasi:

```text
TracerProvider
Resource
OTLP Exporter
Sampler
```

Acceptance criteria:

Jaeger menampilkan:

```text
order-service
inventory-service
payment-service
```

---

## Fase 6 — HTTP Instrumentation

Implementasi:

```text
HTTP server instrumentation
HTTP client instrumentation
context propagation
```

Target:

```text
Order
 ├── Inventory
 └── Payment
```

dalam satu trace.

---

## Fase 7 — Database Instrumentation

Target:

```text
Inventory
   |
   └── PostgreSQL
```

database span terlihat di Jaeger.

---

## Fase 8 — Manual Instrumentation

Tambahkan:

```text
validate-order
check-inventory
process-payment
```

---

## Fase 9 — Attributes, Events & Baggage

Implementasikan:

```text
Span Attributes
Span Events
Baggage
```

dan validasi propagation.

---

## Fase 10 — Chaos Engineering

Implementasikan:

```text
20% delay
2 second delay
10% error
```

---

## Fase 11 — Integration & E2E Testing

Test:

```text
Checkout success
Inventory not found
Payment error
Payment delay
```

---

## Fase 12 — Validation & Troubleshooting

Lakukan eksperimen:

```text
Broken propagation
Collector failure
Database latency
Payment latency
Payment failure
Sampling
```

Developer harus mampu menentukan root cause menggunakan trace.

---

# 39. Functional Requirements

| ID     | Requirement                                   |
| ------ | --------------------------------------------- |
| FR-001 | Order menyediakan `POST /checkout`            |
| FR-002 | Inventory menyediakan `GET /inventory/:id`    |
| FR-003 | Payment menyediakan `POST /pay`               |
| FR-004 | Semua service menyediakan `/health`           |
| FR-005 | Order memanggil Inventory                     |
| FR-006 | Order memanggil Payment                       |
| FR-007 | Inventory menggunakan PostgreSQL              |
| FR-008 | Semua service menggunakan OpenTelemetry SDK   |
| FR-009 | Trace dikirim melalui OTLP                    |
| FR-010 | HTTP context harus dipropagasikan             |
| FR-011 | Database query menghasilkan span              |
| FR-012 | Harus terdapat manual span                    |
| FR-013 | Harus terdapat span attributes                |
| FR-014 | Harus terdapat span events                    |
| FR-015 | Harus terdapat baggage                        |
| FR-016 | Harus terdapat configurable sampling          |
| FR-017 | Payment dapat menghasilkan delay              |
| FR-018 | Payment dapat menghasilkan HTTP 500           |
| FR-019 | OTel Collector menjadi intermediary telemetry |
| FR-020 | Trace dapat divisualisasikan melalui Jaeger   |
| FR-021 | Unit test tersedia untuk business logic       |
| FR-022 | Linting tersedia                              |
| FR-023 | Coverage report tersedia                      |
| FR-024 | Integration test tersedia                     |
| FR-025 | End-to-end test tersedia                      |

---

# 40. Non-Functional Requirements

## Reliability

Tracing tidak boleh menjadi dependency utama business transaction.

Jika Collector mengalami gangguan, service idealnya tetap dapat menjalankan business logic selama memungkinkan.

## Maintainability

Konfigurasi dipisahkan dari source code menggunakan environment variable.

## Reproducibility

Infrastructure dan deployment harus dapat dibuat ulang dari repository.

## Code Quality

Source code wajib memenuhi:

```text
go fmt
go vet
golangci-lint
unit test
```

---

# 41. Acceptance Criteria Utama

Project dianggap berhasil apabila:

### AC-001 — End-to-End Trace

Satu:

```http
POST /checkout
```

menghasilkan satu trace dengan:

```text
Order
Inventory
Payment
```

### AC-002 — Context Propagation

Trace ID sama pada seluruh downstream service.

### AC-003 — Database Span

Inventory memiliki child PostgreSQL span.

### AC-004 — Manual Span

Custom business span terlihat di Jaeger.

### AC-005 — Attributes

Span attributes dapat dilihat.

### AC-006 — Events

Span events dapat dilihat.

### AC-007 — Baggage

`order.id` terpropagasi ke downstream.

### AC-008 — Slow Request

Payment delay sekitar 2 detik terlihat pada waterfall.

### AC-009 — Error Request

Payment 500 menghasilkan error span.

### AC-010 — Collector

Semua telemetry melewati OTel Collector sebelum Jaeger.

### AC-011 — Unit Testing

Unit test berhasil:

```text
go test ./...
```

### AC-012 — Linting

Lint berhasil:

```text
golangci-lint run ./...
```

### AC-013 — Static Analysis

Berhasil:

```text
go vet ./...
```

### AC-014 — Coverage

Business logic mencapai:

```text
>= 70%
```

### AC-015 — Integration Test

Alur antar service berhasil diuji.

### AC-016 — E2E Test

Deployment Kubernetes dapat diuji menggunakan test script.

---

# 42. Skenario Pengujian

## Test 1 — Normal Checkout

```text
PAYMENT_DELAY_PERCENT=0
PAYMENT_ERROR_PERCENT=0
```

Request:

```bash
curl -X POST http://localhost:18080/checkout \
  -H "Content-Type: application/json" \
  -d '{"item_id":"A123","qty":1}'
```

Expected:

```text
HTTP 200
```

Trace:

```text
Order
 ├── Inventory
 │    └── PostgreSQL
 │
 └── Payment
```

---

## Test 2 — Slow Payment

```text
PAYMENT_DELAY_PERCENT=100
PAYMENT_DELAY_MS=2000
PAYMENT_ERROR_PERCENT=0
```

Expected:

```text
Payment duration ≈ 2 seconds
```

---

## Test 3 — Payment Error

```text
PAYMENT_DELAY_PERCENT=0
PAYMENT_ERROR_PERCENT=100
```

Expected:

```text
Payment = HTTP 500
Order = checkout failed
Trace = ERROR
```

---

## Test 4 — Inventory Not Found

```bash
curl -X POST http://localhost:18080/checkout \
  -H "Content-Type: application/json" \
  -d '{"item_id":"UNKNOWN","qty":1}'
```

Expected:

```text
Inventory = 404
Order = checkout failed
```

---

# 43. Skenario Eksperimen Pembelajaran

## Eksperimen 1 — Context Propagation

Bandingkan:

```text
Dengan propagation:

Order
 ├── Inventory
 └── Payment
```

dengan:

```text
Tanpa propagation:

Trace A
 └── Order

Trace B
 └── Inventory

Trace C
 └── Payment
```

---

## Eksperimen 2 — Payment Bottleneck

```text
PAYMENT_DELAY_PERCENT=100
PAYMENT_DELAY_MS=2000
```

Analisis waterfall.

---

## Eksperimen 3 — Database Bottleneck

Tambahkan artificial database delay.

Bandingkan DB latency dengan Payment latency.

---

## Eksperimen 4 — Error Propagation

```text
PAYMENT_ERROR_PERCENT=100
```

Analisis error chain.

---

## Eksperimen 5 — Sampling

Bandingkan:

```text
100%
```

dengan:

```text
10%
```

sampling.

---

## Eksperimen 6 — Collector Failure

Hentikan Collector dan amati:

```text
Application behavior
Telemetry behavior
Error handling
```

Tujuannya memahami separation antara business dependency dan observability dependency.

---

# 44. Output Akhir Project

Repository harus menghasilkan:

```text
1. Source code Order Service
2. Source code Inventory Service
3. Source code Payment Service
4. Container images
5. Kubernetes manifests
6. PostgreSQL schema
7. Database seed
8. OTel Collector configuration
9. Jaeger deployment
10. Unit tests
11. Integration tests
12. E2E tests
13. Lint configuration
14. Coverage report
15. Quality-check scripts
16. Deployment scripts
17. Test scripts
18. README
19. Dokumentasi tracing
20. Dokumentasi eksperimen
```

---

# 45. Definition of Done

Project dinyatakan **DONE** apabila seluruh kondisi berikut terpenuhi:

- [ ] k3d cluster berhasil dijalankan
- [ ] PostgreSQL running
- [ ] OTel Collector running
- [ ] Jaeger running
- [ ] Order Service running
- [ ] Inventory Service running
- [ ] Payment Service running
- [ ] `/health` semua service menghasilkan 200
- [ ] Checkout berhasil
- [ ] Distributed trace muncul
- [ ] Trace context terpropagasi
- [ ] Inventory DB span muncul
- [ ] Manual span muncul
- [ ] Span attributes muncul
- [ ] Span events muncul
- [ ] Baggage berhasil dipropagasikan
- [ ] Sampling dapat dikonfigurasi
- [ ] Payment delay dapat diamati
- [ ] Payment error dapat diamati
- [ ] `go fmt` berhasil
- [ ] `go vet ./...` berhasil
- [ ] `golangci-lint run ./...` berhasil
- [ ] `go test ./...` berhasil
- [ ] Coverage business logic ≥ 70%
- [ ] Integration test berhasil
- [ ] E2E test berhasil
- [ ] Chaos test berhasil
- [ ] Seluruh project dapat direbuild dan redeploy dari repository

---

# 46. Target Trace Akhir

## Normal

```text
POST /checkout                         ~100ms
│
├── validate-order                      ~1ms
│
├── GET /inventory/A123                ~20ms
│   │
│   └── PostgreSQL SELECT               ~5ms
│
└── POST /pay                           ~50ms
```

## Slow Payment

```text
POST /checkout                         ~2.1s
│
├── validate-order                      ~1ms
│
├── GET /inventory/A123                ~20ms
│   │
│   └── PostgreSQL SELECT               ~5ms
│
└── POST /pay                           ~2.0s
```

## Payment Error

```text
POST /checkout
│
├── GET /inventory/A123
│   └── PostgreSQL SELECT
│
└── POST /pay
      │
      └── ERROR: HTTP 500
```

---

# 47. Struktur Final Sistem

```text
                         +------------------+
                         |      Client      |
                         +--------+---------+
                                  |
                                  | :18080
                                  v
                       +----------------------+
                       |    Order Service     |
                       |       :18080         |
                       +----------+-----------+
                                  |
                    +-------------+-------------+
                    |                           |
                    v                           v
          +-------------------+       +-------------------+
          | Inventory Service |       |  Payment Service  |
          |      :18081       |       |      :18082       |
          +---------+---------+       +-------------------+
                    |
                    v
             +-------------+
             | PostgreSQL  |
             |   :15432   |
             +-------------+

                    All Services
                           |
                           | OTLP
                           v
                +----------------------+
                |  OTel Collector      |
                | :14317 / :14318      |
                +----------+-----------+
                           |
                           v
                    +-------------+
                    |    Jaeger   |
                    |    :16686   |
                    +-------------+
```

---

# 48. Goal Akhir Pembelajaran

Developer harus mampu menggunakan OTel-Shop Lab untuk menjawab pertanyaan operasional:

```text
Kenapa checkout lambat?
```

```text
Service mana yang menjadi bottleneck?
```

```text
Apakah bottleneck berasal dari HTTP atau PostgreSQL?
```

```text
Service mana yang menghasilkan error?
```

```text
Apakah trace context berhasil diteruskan?
```

```text
Apakah baggage berhasil dipropagasikan?
```

```text
Apakah Collector menerima telemetry?
```

```text
Apakah sampling mempengaruhi trace yang terlihat?
```

Dan yang paling penting:

```text
Bagaimana menemukan root cause
sebuah distributed transaction
hanya dengan membaca trace?
```

Target akhir pembelajaran:

```text
                    CODE
                      │
                      ▼
              ┌───────────────┐
              │ Unit Testing  │
              │    + Lint     │
              └───────┬───────┘
                      │
                      ▼
                CONTAINER
                      │
                      ▼
                 KUBERNETES
                      │
                      ▼
              ┌───────────────┐
              │ OpenTelemetry │
              │     SDK       │
              └───────┬───────┘
                      │
                      ▼
              CONTEXT PROPAGATION
                      │
                      ▼
                    OTLP
                      │
                      ▼
              OTel Collector
                      │
                      ▼
                   JAeger
                      │
                      ▼
              TRACE ANALYSIS
                      │
                      ▼
             ROOT CAUSE ANALYSIS
```

Dengan demikian OTel-Shop Lab menjadi satu project pembelajaran yang mencakup **software quality → containerization → Kubernetes → OpenTelemetry instrumentation → distributed tracing → telemetry pipeline → chaos → testing → troubleshooting → root cause analysis** secara end-to-end.
