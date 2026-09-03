# blogpost.md — Day 6

---

## Section 1 — Blog Post

### Post Title
Day 6 OTel-Shop: Traces That Tell the Whole Story

### Slug
day-6-otel-shop-traces-that-tell-the-whole-story

### URL
http://stayrelevant.id/blog/day-6-otel-shop-traces-that-tell-the-whole-story

### Excerpt (Meta Description)
Day 6 of OTel-Shop: database spans, manual business spans, rich attributes, payment events, and baggage carrying the order id across every service.

### Tags
otel-shop, golang, open-telemetry, distributed-tracing, observability, kubernetes, k3d, jaeger, postgresql, microservices, devops, cloud-native

### Cover Image Prompt
Flat vector illustration, wide 16:9 tech blog cover banner, a large magnifying glass hovering over a monitor showing a tracing waterfall chart, but the bars under the glass reveal extra detail: tiny child bars branching off (a database cylinder symbol, three small step icons in a row) and floating label bubbles with a tag icon and a small id badge, a bank card block at the side emitting small confetti dots for events, clean minimal style, dark navy background with teal and orange accent colors, subtle isometric grid faded in the background suggesting Kubernetes, no text, crisp smooth shapes, JPG output compressed under 100KB.

### Content

## Where We Left Off

Day 5 delivered the first distributed trace: one trace id spanning all three services. Amazing — but honest review time. Every span had a generic name, the waterfall showed "an HTTP call happened," and that was it. If a checkout were slow, the trace still couldn't say whether validation, the database, or payment was the culprit. The trace knew the route; it didn't know the story.

## What Got Added Today

Three enrichment layers, all landing in the same waterfall.

**The business story, as spans.** Checkout now explicitly starts three named spans around its steps: `validate-order`, `check-inventory`, and `process-payment`. Each carries its own attributes — item id and quantity on validation, product id and resulting stock on the inventory check, order id, amount, and final status on payment. Reading the waterfall top to bottom now reads like the checkout procedure itself.

**The database finally speaks.** Inventory's connection pool is now opened through `otelsql`, a drop-in wrapper for `database/sql`. Zero query changes — the same `SELECT stock FROM products WHERE id = $1` — yet every query emits spans: you can literally see `sql.conn.query` and `sql.rows` nested under the inventory request, with `db.system=postgresql` attached.

**Metadata everywhere.**
- The payment span now records events as they happen: `payment_started`, then `payment_completed` (or `payment_failed`). A timeline inside a span.
- When chaos forces a failure, spans are marked with error status and the exception is recorded — the parent checkout span too, so the very first bar you see turns red.
- And the small but delightful one: **baggage**. The order id is now created at the very start of checkout and attached to the request context. OpenTelemetry's propagators carry it inside HTTP headers automatically, and both inventory and payment read it back and stamp it onto their spans. One id, visible everywhere, zero extra HTTP calls.

## What It Looks Like

Real output from one checkout request:

- Order: `POST /checkout` → `validate-order` → `check-inventory` → `process-payment`, each with its attributes.
- Inventory: `GET /inventory/{id}` stamped with `baggage.order.id=O-a921116e` and `product.stock=10`, with the actual database query nested beneath it.
- Payment: `POST /pay` with the same baggage id, `payment.amount`, status, and the two lifecycle events.

And the failure path, verified by flipping the chaos knob to 100%: the payment span turns error-red with a `payment_failed` event, the order-side payment span and the parent checkout span light up red too. Failure propagates visibly up the whole chain — exactly what you want at 3am.

## Lessons Learned

- **Auto-instrumentation gets you 80%; business spans get you meaning.** `otelsql` and `otelhttp` required almost no code, but only the hand-named spans (`check-inventory`) map to what the business actually does.
- **Baggage is a write-once, read-everywhere channel.** Set a value once at the entry point, and every downstream service can read it — perfect for correlation ids like `order.id`. (It's not for bulk data, just small metadata.)
- **Events make spans self-documenting.** A timeline of `payment_started` → `payment_failed` inside one span tells the story without needing a second trace.
- **Fail loudly, at every level.** Marking the child span, the client call, and the parent checkout span as errors means nobody has to dig to notice a broken flow.

## Conclusion

With Day 6, the tracing story is complete: routes, business steps, database queries, attributes, events, and shared context all live in one waterfall. The lab now demonstrates the full OpenTelemetry surface that real production systems rely on.

What's left is polish: a proper quality script with the coverage gate, integration tests, a chaos-test script, and the documentation that ties the whole experiment together.

Repo is here: https://github.com/stayrelevantid/otel-shop

---

## Section 2 — LinkedIn Post

**Hook**
Day 6 of the OTel-Shop challenge: our traces learned to tell the whole story 📖

**Today's wins:**
• Manual business spans: `validate-order` → `check-inventory` → `process-payment`, each with its own attributes
• Database spans via otelsql — the actual Postgres query now appears nested under the inventory request
• Baggage: the order id is set once at checkout and read back by both downstream services, zero extra calls
• Payment span events (`payment_started` → `payment_completed`) and loud error status when chaos strikes — parent span turns red too

The waterfall now reads like the checkout procedure itself. One trace = the full story.

Next: quality scripts (coverage gate + chaos test), integration tests, and documentation to wrap up the lab.

Blog: http://stayrelevant.id/blog/day-6-otel-shop-traces-that-tell-the-whole-story
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