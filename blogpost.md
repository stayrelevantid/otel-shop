# blogpost.md — Day 4

---

## Section 1 — Blog Post

### Post Title
Day 4 OTel-Shop: Chaos Injection and Real Unit Tests

### Slug
day-4-otel-shop-chaos-injection-real-unit-tests

### URL
http://stayrelevant.id/blog/day-4-otel-shop-chaos-injection-real-unit-tests

### Excerpt (Meta Description)
Day 4 of OTel-Shop: payment learned to fail on purpose via configurable chaos, and the checkout flow got real unit tests above 70% coverage.

### Tags
otel-shop, golang, open-telemetry, distributed-tracing, observability, kubernetes, k3d, jaeger, postgresql, microservices, devops, cloud-native

### Cover Image Prompt
Flat vector illustration, wide 16:9 tech blog cover banner, a cartoon bank card service block standing in a light rain storm with small dice floating around it (one showing six dots, lightning bolt symbol overhead) while glowing teal arrows of requests still flow toward it from a shopping cart block, on the right side a neat row of test tubes filled with green liquid and a bar chart showing tall green coverage bars, clean minimal style, dark navy background with teal and orange accent colors, subtle isometric grid faded in the background suggesting Kubernetes, no text, crisp smooth shapes, JPG output compressed under 100KB.

### Content

## Where We Left Off

Day 3 gave us a genuinely distributed checkout: order calls inventory, inventory reads PostgreSQL, payment gets charged. Everything green. But green-everywhere is exactly when you should get suspicious — the flow had never been seen failing, and not a single line of it was covered by a unit test. Day 4 fixed both.

## What Got Built Today

**Payment learned to fail on purpose.** A tiny chaos package now sits inside the payment service. It reads three knobs from the environment: how often to wait before answering (in percent), how long that wait is (in milliseconds), and how often to outright fail. The important design choice: at 0% or 100% the behavior is completely deterministic, so it's testable — randomness only lives in the middle range. It also respects the request context, so a cancelled request doesn't keep it napping.

Then the fun part — flipping the switches in the real cluster:

- `PAYMENT_ERROR_PERCENT=100` → checkout honestly reports `payment failed` (500) ✅
- Revert to the default 10% → checkout is `paid` again ✅

No code change, no redeploy — just an environment variable. That's the whole point of chaos knobs.

**The whole flow got real tests.** Thanks to yesterday's interfaces, fakes were cheap to write:

- Order's checkout tested with mock inventory and payment clients: success, bad input, inventory failure, payment failure, and insufficient stock — every branch.
- Order's actual HTTP clients tested against a fake server (httptest), so the real request/response code is exercised too.
- Inventory's handler tested with a fake store, and its database layer tested with sqlmock — including the "row not found" path — without ever touching PostgreSQL.
- Payment's chaos behavior tested deterministically (always fails / never fails / delays actually delay / cancelled context bails out) plus the handler around it.

## The Coverage Table

The goal was 70% on every internal package. Actual numbers:

| Service | Package | Coverage |
|---|---|---|
| Order | handler | 90.9% |
| Order | client | 82.1% |
| Inventory | db | 90.0% |
| Inventory | handler | 82.6% |
| Payment | chaos | 94.7% |
| Payment | handler | 100% |

All comfortably above the bar — and every green `test.sh` run now has backup from 20+ unit tests.

## Lessons Learned

- **Chaos only works if it's deterministic at the edges.** 0% and 100% must be promises, not probabilities — otherwise you can't test the chaos itself.
- **Interfaces written "for later" pay off immediately.** Every mock in today's tests existed because of interfaces drawn on Day 3, before any test needed them.
- **A coverage gate finds real gaps.** Watching the numbers forced proper client tests instead of stopping at handler-only coverage — the HTTP layer would have stayed untested otherwise.
- **Keep the database out of unit tests.** sqlmock simulates the exact scenarios (rows, no-rows, dead connection) in milliseconds, with zero setup and zero flakes.

## Conclusion

Day 4 closes the quality-foundation chapter: the system can now be broken on demand, observed failing honestly, and every branch of the checkout flow is pinned down by tests. If tomorrow's changes regress something, the tests will say so before any human notices.

Next is the reason this whole lab exists: the OpenTelemetry SDK goes in, and the first real distributed traces appear in Jaeger. See you there.

Repo is here: https://github.com/stayrelevantid/otel-shop

---

## Section 2 — LinkedIn Post

**Hook**
Day 4 of the OTel-Shop challenge: taught payment how to fail on purpose — then tested it 🎲

**Today's wins:**
• Chaos in payment: delay % + error % from env vars, deterministic at 0/100 so it's testable
• Flipped `PAYMENT_ERROR_PERCENT=100` in the cluster — checkout honestly reported `payment failed`, reverted, `paid` again. No code change, no redeploy
• Real unit tests across all three services: mock clients (httptest), fake stores, sqlmock for the DB layer
• Coverage gate ≥70% cleared: 82.1%–100% across six internal packages

Small drama: duplicate field names in a test struct — five minutes, fixed. Tests now guard every branch of checkout.

Next: the OpenTelemetry SDK goes in and the first real distributed traces land in Jaeger 🔭

Blog: http://stayrelevant.id/blog/day-4-otel-shop-chaos-injection-real-unit-tests
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
