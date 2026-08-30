# blogpost.md — Day 5

---

## Section 1 — Blog Post

### Post Title
Day 5 OTel-Shop: The First Distributed Trace in Jaeger

### Slug
day-5-otel-shop-first-distributed-trace-in-jaeger

### URL
http://stayrelevant.id/blog/day-5-otel-shop-first-distributed-trace-in-jaeger

### Excerpt (Meta Description)
Day 5 of OTel-Shop: OpenTelemetry SDK integrated, every HTTP call instrumented, and the first full trace across all three services landed in Jaeger.

### Tags
otel-shop, golang, open-telemetry, distributed-tracing, observability, kubernetes, k3d, jaeger, postgresql, microservices, devops, cloud-native

### Cover Image Prompt
Flat vector illustration, wide 16:9 tech blog cover banner, three rounded service blocks (shopping cart, database cylinder, bank card) connected left to right by glowing teal dotted arrows, each arrow carrying a small glowing dot labeled visually with a shared ID badge symbol, the final arrow flows into a large monitor screen on the right showing a tracing waterfall chart with teal and orange horizontal bars of different widths, one bar noticeably stretched long, clean minimal style, dark navy background with teal and orange accent colors, subtle isometric grid faded in the background suggesting Kubernetes, no text, crisp smooth shapes, JPG output compressed under 100KB.

### Content

## Where We Left Off

After four days we had a cluster, three talking services, chaos knobs, and a test suite. Everything behaved — but we were flying blind. When checkout took two seconds, I had to guess where the time went. Day 5 is the reason this whole lab exists: wire in OpenTelemetry and finally see the request's full journey.

## The Big Picture: One Shared SDK Setup

The nice part: the entire SDK setup lives in one small shared package, `pkg/telemetry`, used by all three services. Its `Init` call does a handful of things:

- Builds a **resource** — who am I? (service name, version, "this is a lab" tag). This is how Jaeger knows which service produced which span.
- Creates an **OTLP exporter** that ships spans over gRPC to the collector, which forwards them to Jaeger.
- Configures a **sampler** driven by an environment variable — keep 100% of traces for now, dial it down later with one env change.
- Registers **propagators** (W3C Trace Context + Baggage) — the glue that lets a trace id hop from service to service inside HTTP headers.

Each service calls `Init` at startup and defers the shutdown so pending spans get flushed before exit.

## Wiring HTTP Servers and Clients

The SDK alone doesn't produce anything — requests need to create spans and forward context. Two helpers do the heavy lifting:

- On the **server side**, every service wraps its router with `otelhttp`. One line per service; now every incoming HTTP request gets a span named after its route.
- On the **client side**, the order service's HTTP clients got `otelhttp.NewTransport`. That transparently injects the `traceparent` header (and baggage) into every outgoing call.

The silent hero underneath is `context.Context` — already threaded through the whole flow on Day 3. When order calls inventory with the request's context, the transport reads the active span from it and serializes the trace id into headers. No extra code, just plumbing done right.

## The Moment of Truth

One `POST /checkout` later, I asked Jaeger's API for the latest checkout trace — and there it was:

```
traceID: d87f3b8fcfa436b1e037007a31e7c868
services: [inventory-service, order-service, payment-service]

  order-service     POST /checkout         2079.5ms
  inventory-service GET /inventory/{id}      15.0ms
  payment-service   POST /pay              2014.3ms
```

One trace id, three different services. Order fan-outs to its two neighbors, and both children nest neatly under the parent span. Even better: this particular checkout got unlucky with the chaos dice — payment's 20% chance of a 2-second delay triggered, and you can see it instantly in the waterfall. That's the "why" of this whole lab in one screenshot: not "something is slow," but exactly which hop and how much.

## Lessons Learned

- **One shared telemetry package beats copy-paste.** Three services, one `Init`, identical behavior — and a sampler change is a one-file edit.
- **`context.Context` is the highway.** Tracing breaks wherever someone drops the context; threading it early on Day 3 paid off completely today.
- **`otelhttp` is a cheat code.** One wrapper per server, one transport per client — production-grade spans and header propagation for near-zero effort.
- **Verify through the API, not just the UI.** Querying Jaeger's API (filtered by operation!) made the proof concrete — and taught me that health-check probes flood traces, so filter wisely.

## Conclusion

The core mission is done: a distributed trace now tells the truth about every checkout, end to end. From here, tracing gets richer — database spans, manual business spans (validate / check / process), attributes, events, and baggage riding along on the same context.

Repo is here: https://github.com/stayrelevantid/otel-shop

---

## Section 2 — LinkedIn Post

**Hook**
Day 5 of the OTel-Shop challenge: the moment of truth — our first end-to-end distributed trace just landed in Jaeger! 🔭

**Today's wins:**
• One shared `pkg/telemetry` package: OTLP exporter, env-driven sampler, W3C TraceContext + Baggage propagation
• Every HTTP server wrapped with otelhttp, every outgoing client call instrumented via otelhttp transport
• Single trace id spanning all three services: checkout → inventory → payment, children nested under the parent span
• Bonus proof: payment's 20% chaos delay triggered — you can see the exact 2-second hop in the waterfall

Silent hero: `context.Context` threaded everywhere since Day 3. Propagation "just worked".

Next: database spans, manual business spans, attributes/events, and baggage.

Blog: http://stayrelevant.id/blog/day-5-otel-shop-first-distributed-trace-in-jaeger
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
