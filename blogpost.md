# blogpost.md — Day 3

---

## Section 1 — Blog Post

### Post Title
Day 3 OTel-Shop: Services Start Talking via HTTP

### Slug
day-3-otel-shop-services-start-talking-via-http

### URL
http://stayrelevant.id/blog/day-3-otel-shop-services-start-talking-via-http

### Excerpt (Meta Description)
Day 3 of OTel-Shop: inventory now reads real PostgreSQL, and checkout actually calls inventory and payment over HTTP. The full flow works end to end.

### Tags
otel-shop, golang, open-telemetry, distributed-tracing, observability, kubernetes, k3d, jaeger, postgresql, microservices, devops, cloud-native

### Cover Image Prompt
Flat vector illustration, wide 16:9 tech blog cover banner, three rounded service blocks (shopping cart, database cylinder, bank card) connected by glowing animated dotted lines that clearly show a conversation: the cart block calls the database block, which answers with a small "10" bubble, then the cart calls the bank card block which responds with a check mark, all orchestrated left to right, a small monitor at the edge showing a green success bar, clean minimal style, dark navy background with teal and orange accent colors, subtle isometric grid faded in the background suggesting Kubernetes, no text, crisp smooth shapes, JPG output compressed under 100KB.

### Content

## Where We Left Off

The last two days gave us a cluster and three services that knew how to say "ok" on `/health` — but they didn't really talk. Inventory answered with a hardcoded number, and checkout never phoned anyone. Cute, but not a system. Day 3's job: make them actually communicate.

## What Got Wired Today

Two chunks of work, and they finally click together.

**Inventory meets the database.** Inventory now opens a real connection pool to PostgreSQL (pgx under the hood) and looks up stock by id. I hid the database behind a tiny `StockStore` interface so the handler doesn't care where the data lives — handy later for tests. Responses are honest now: found → 200 with the number, missing → 404, anything worse → 500.

**Checkout learns to orchestrate.** Order grew two HTTP clients — one for inventory, one for payment — both behind interfaces so they can be faked in tests. The checkout handler now does a proper little dance:
- Validate the request (no empty item, qty must be positive).
- Ask inventory for stock.
- Complain if there isn't enough.
- Charge payment for `qty * 10` (price is a flat 10 per unit for now).
- Only then return a fresh order id with status `paid`.

When a neighbor fails, checkout stays honest instead of pretending: a downstream error becomes a 500 with a clear message like `inventory failed` or `payment failed`.

How it behaved end to end (straight from the host, through the cluster):

- `POST /checkout` for A123 ×1 → `{"order_id":"O-...","status":"paid"}` ✅
- Unknown product → 500 `inventory failed` ✅
- C123 ×999 (only 5 in stock) → 400 `insufficient stock` ✅
- Garbage request → 400 ✅

And I finally added a proper `test.sh` that runs all nine of these checks and exits non-zero if anything breaks — so "does it still work?" is one command now, not a memory test.

## The Gotcha

First attempt at local testing looked green... but for the wrong reasons. My local services were trying to bind ports 18080–18082, which the cluster already occupies as NodePorts on my laptop. So my curl requests silently hit the old stub pods, not my new code. The fix: run local services on throwaway ports and point the order client at those. Silly, but it's exactly the kind of "it works on my machine because it's actually the cluster" trap worth internalizing.

## Lessons Learned

Four things stuck this time:

- **Hide infrastructure behind interfaces early.** The `StockStore` and client interfaces cost almost nothing today and are what make the flow testable tomorrow.
- **A failed call is data, not a mystery.** Returning a clear `inventory failed` / `payment failed` turned a confusing 500 into something you can act on.
- **Test like a real client, on real ports.** My port-clash slip proved that "it returned 200" means nothing if you're not sure which process answered.
- **One script beats a checklist in your head.** `test.sh` replaced five manual curls and already caught me thinking I was done when I wasn't.

## Conclusion

Day 3 closes with a checkout flow that's genuinely distributed: one request fans out to inventory and payment, talks to a real database, and reports honest results. The architecture finally earns the "microservices" label.

Next: inject some chaos into payment (random delays and errors) and start writing real unit tests behind those interfaces. See you then.

Repo is here: https://github.com/stayrelevantid/otel-shop

---

## Section 2 — LinkedIn Post

**Hook**
Day 3 of the OTel-Shop challenge: the services finally talk to each other 🔗

**Today's wins:**
• Inventory reads real PostgreSQL via a pgx pool, behind a StockStore interface (200/404/500)
• Order grew HTTP clients for inventory & payment, also behind interfaces — testable by design
• Checkout now orchestrates: validate → stock check → charge qty×10 → 200 paid, with clear 500s on downstream failure
• Added `test.sh`: 9 E2E checks, exits non-zero on any failure

Small drama: my first local test looked green because it was secretly hitting the old cluster pods on the same NodePorts. Lesson re-learned: know exactly which process answered your request.

Next: chaos engineering on payment (random delay + error) and proper unit tests behind those interfaces.

Blog: http://stayrelevant.id/blog/day-3-otel-shop-services-start-talking-via-http
Repo: https://github.com/stayrelevantid/otel-shop

#OpenTelemetry #Golang #Kubernetes #Observability #DevOps #DistributedTracing #Jaeger #PostgreSQL #Microservices #CloudNative

---

## Section 3 — Project Showcase

### Project Title
OTel-Shop Lab

### Slug
otel-shop-lab

### Showcase Image Prompt
Isometric hero illustration for portfolio project card, three connected service blocks (labeled visually with a cart icon, a database cylinder, and a bank card) sitting side by side, glowing neon arrows flowing from the three blocks into a pipeline funnel on the right that ends at a monitor screen showing a tracing waterfall chart with teal and orange bars, dark navy background, subtle grid floor, soft lighting, cinematic depth, clean modern tech aesthetic, no readable text, 16:9, JPG under 100KB.

### Description
End-to-end Golang microservices + OpenTelemetry lab: three small services (Order, Inventory, Payment) deployed on k3d Kubernetes, with every trace flowing through the OpenTelemetry Collector and visualized in Jaeger. Learning focus: distributed tracing, context propagation, baggage, all the way to root cause analysis of a distributed transaction.

### Tags
golang, opentelemetry, distributed-tracing, observability, kubernetes, k3d, jaeger, postgresql, microservices, devops, cloud-native
