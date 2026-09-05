# blogpost.md — Day 7

---

## Section 1 — Blog Post

### Post Title
Day 7 OTel-Shop: Quality Gates, Docs, and Done

### Slug
day-7-otel-shop-quality-gates-docs-and-done

### URL
http://stayrelevant.id/blog/day-7-otel-shop-quality-gates-docs-and-done

### Excerpt (Meta Description)
Day 7 of OTel-Shop, the finale: a full quality pipeline, integration tests, real documentation, and the final DoD gate — all 15 checks passed.

### Tags
otel-shop, golang, open-telemetry, distributed-tracing, observability, kubernetes, k3d, jaeger, postgresql, microservices, devops, cloud-native

### Cover Image Prompt
Flat vector illustration, wide 16:9 tech blog cover banner, a checkered finish-line banner held up between two poles, three rounded service blocks (shopping cart, database cylinder, bank card) crossing the line side by side with small motion lines, behind them a monitor showing a tracing waterfall chart, above floats a row of green check mark badges and a bar chart with all bars filled solid green, clean minimal style, dark navy background with teal and orange accent colors, subtle isometric grid faded in the background suggesting Kubernetes, no text, crisp smooth shapes, JPG output compressed under 100KB.

### Content

## Where We Left Off

Day 6 finished the tracing story: business spans, database spans, attributes, events, and baggage — the waterfall read like the checkout procedure itself. The lab was feature-complete but messy around the edges: verification meant typing commands from memory, documentation was scattered, and nothing proved the project actually met its own acceptance criteria. Day 7 was the finish line: quality gates, real tests, real docs, and a final checklist.

## What Got Built Today

**One command to judge everything.** A `check.sh` script now runs the full quality pipeline across every module: formatting check (fails if anything is unformatted), static vetting, the linter, all tests, and a coverage gate that refuses anything below 70% per package. Result: all green, with coverage sitting between 82% and 97%.

**Breaking things on schedule.** A `chaos-test.sh` script flips the payment chaos knobs against the real cluster and makes two promises: forced errors must actually fail the checkout, and a 100% delay must actually take about two seconds (measured: 2034ms). Then it reverts everything. Broken-on-demand is now reproducible with one command.

**A test that exercises the real wiring.** Unit tests mock the neighbors; the new integration test does the opposite — it uses the *real* checkout handler with the *real* HTTP clients, pointed at fake downstream servers. Three scenarios: the happy path, a missing product (fails honestly with 500), and insufficient stock (400). If the wiring ever breaks, this catches it.

**Docs that explain, not just describe.** Two new documents: one showing what traces look like in the three canonical situations (normal, slow, failing), and one walking through six core observability experiments — from turning propagation off to dialing sampling down to 10%. Plus a README that finally matches reality: architecture, ports, quickstart, teardown.

## The DoD Gate

The PRD defined fifteen "definition of done" checks, and the last day's job was to actually run them — including the ones nobody had tested live before:

- **Sampling at 10%.** Set the sampler env to 0.1, fired 20 checkouts, counted traces: 3 showed up (~15%, within expectations). Proof that sampling controls volume while traces that survive stay complete.
- **Collector down.** Scaled the telemetry collector to zero replicas — the component that receives every trace — and fired checkouts: all 200 OK. Telemetry must never hold the business hostage; when the collector came back, traces resumed flowing.

Plus the familiar ones re-verified end to end: single trace across all three services, database spans, manual spans, attributes, events, baggage, chaos behavior, full test suite, and a rebuild-from-clean loop.

All fifteen: passed. The project's own acceptance criteria, executed rather than assumed.

## Lessons Learned

- **A gate beats a habit.** Coverage above 70% was an intention for six days; it became a fact in one afternoon once a script refused anything less.
- **Deterministic chaos is the whole trick.** The chaos knobs pass at 0% and 100% and only get random in between — that's why they can be tested by a script with assertions.
- **Unit tests catch logic; integration tests catch wiring.** Mocks verified the rules; the real-handler-real-clients test verified the pipes between them. You need both.
- **Write the docs for the person you were on Day 1.** Every confusing term deserves a plain-words explanation — the repo now even ships a glossary that explains spans, baggage, NodePorts and friends like you're five.

## Conclusion

Seven days, from an empty folder to a small but complete microservices lab: three Go services in Kubernetes, a real checkout flow, honest failures, full OpenTelemetry instrumentation, chaos knobs, quality gates, and documentation that teaches. Every claim in this series ends with the same receipt: it's in the repo, and the checks say it's true.

The lab is done. What a ride.

Repo is here: https://github.com/stayrelevantid/otel-shop

---

## Section 2 — LinkedIn Post

**Hook**
Day 7 of the OTel-Shop challenge — the finale: quality gates, docs, and a clean DoD pass. The lab is done 🏁

**Today's wins:**
• `check.sh`: one command for fmt + vet + lint + tests + coverage gate — all green (82–97%)
• `chaos-test.sh`: forced errors and a measured 2-second delay, asserted and reverted automatically
• Integration test with the real handler and real HTTP clients against fake downstreams
• DoD gate: 15/15 checks passed — including live proof that sampling 10% keeps ~1 in 10 traces, and that killing the telemetry collector never breaks checkout

Seven days: cluster up → services talking → first distributed trace → chaos + tests → rich spans → quality gates. Every claim has a receipt in the repo.

Bonus: the repo now includes a plain-words glossary explaining every confusing term (spans, baggage, NodePorts...) — the doc I wish I had on Day 1.

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