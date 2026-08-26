# blogpost.md — Day 1

---

## Section 1 — Blog Post

### Post Title
Day 1 OTel-Shop: Naikin Infra k3d, Jaeger & Postgres

### Slug
day-1-otel-shop-naikin-infra-k3d-jaeger-postgres

### URL
http://stayrelevant.id/blog/day-1-otel-shop-naikin-infra-k3d-jaeger-postgres

### Excerpt (Meta Description)
Day 1 of building OTel-Shop from scratch: k3d cluster up, Postgres seeded, OTel Collector and Jaeger running, and the first trace made it through.

### Tag
post

### Cover Image Prompt
Flat vector illustration, wide 16:9 tech blog cover banner, a laptop on a desk showing a tracing waterfall UI with colorful horizontal bars, above the laptop float three connected rounded containers each with a small symbol (shopping cart, database cylinder, bank card) linked by glowing lines flowing right into a pipeline funnel icon, clean minimal style, dark navy background with teal and orange accent colors, subtle isometric grid faded in the background suggesting Kubernetes, no text, crisp smooth shapes, JPG output compressed under 100KB.

### Content

## Why I Started This Series

So here's the thing. I really wanted to genuinely understand OpenTelemetry and distributed tracing. Not just skim docs and go "yeah, got it," but actually build something end to end. That's how this challenge, OTel-Shop, was born: a small microservices lab in Go where all traffic gets visualized in Jaeger. Fun, right?

Today (Day 1) wasn't about features, it was about foundations. My rule: if the base is solid, everything after is just stacking bricks. So most of the day went into setting up the lab infrastructure.

## What I Got Done Today

Quick recap so I don't ramble:

- Initialized the repo from scratch with a `.gitignore` and ran a gitleaks scan so no secrets sneak in.
- Set up a Go workspace with 4 modules: order, inventory, payment, and telemetry.
- Wrote Kubernetes manifests for PostgreSQL, OpenTelemetry Collector, and Jaeger.
- Brought up a local k3d cluster. All three core components came up: database, collector, and Jaeger UI.
- Sent a tiny fake trace and watched it show up in Jaeger. The pipeline was proven end to end from day one, not just on paper.

The most satisfying part wasn't things running smoothly — it was sending that tiny little trace and seeing it land in Jaeger. That's the "oh, so this is observability" moment.

## The Drama Along the Way

Don't imagine everything was smooth — that's where the real value was.

First, k3d (the tool that creates local Kubernetes clusters) refuses to use Podman. It wants Docker, period. My machine was all-in on Podman. Brief drama, then Docker Desktop got switched on. Lesson: check what container runtime your tools actually support before you get excited about deploying.

Second, image downloads were painfully slow. One took nine minutes. My deploy script even timed out on it. Exhausting? Yes. But it reinforced why you pin image versions — so you're not silently pulling whatever "latest" is and wondering why behavior changed.

Third, a small classic error: `go build ./...` from the repo root. Turns out with a multi-module workspace that command doesn't work at the root. Fix: build per module. Small thing, but it made me understand the workspace structure better.

## Lessons Learned

Quick takeaway, four things:

- **Check your toolchain before you start.** If I'd checked that k3d needs Docker first, I wouldn't have spent time fighting Podman.
- **Pin every image version.** "Latest" is a silent commitment that makes your infra non-reproducible.
- **Get one tiny "proof" early.** Seeing a trace land in Jaeger on Day 1 is far more motivating than waiting for every feature to finish.
- **Neat infra is your future self's treat.** Tomorrow it's just deploying services, no cluster debugging first.

## Conclusion

Day 1 wrapped in a good place: cluster up, database seeded, collector running, and Jaeger already receiving its first trace. Tired, but relieved — the fiddliest foundation is done.

Tomorrow, the real services: the checkout flow, stock check, and payment. Hopefully less drama than today.

Repo is here: https://github.com/stayrelevantid/otel-shop

---

## Section 2 — LinkedIn Post

**Hook**
Day 1 of the OTel-Shop challenge: building an observability stack from zero 🚀

**Today's wins:**
• Local k3d cluster is up — PostgreSQL seeded, OTel Collector running
• Jaeger UI is live, and the first trace made it through the full pipeline
• Repo started properly: .gitignore + gitleaks scan so no secrets slip in

Small drama: k3d insists on Docker while I was all-in on Podman, plus painfully slow image pulls. Valuable lessons though: always pin image versions and check your toolchain first 🙃

Tomorrow I'm coding the three services. Wanna follow along? Blog post is here 👇

Blog: http://stayrelevant.id/blog/day-1-otel-shop-naikin-infra-k3d-jaeger-postgres
Repo: https://github.com/stayrelevantid/otel-shop

#OpenTelemetry #Golang #Kubernetes #Observability #DevOps

---

## Section 3 — Project Showcase

### Project Title
OTel-Shop Lab

### Slug
otel-shop-lab

### Showcase Image Prompt
Isometric hero illustration for portfolio project card, three connected service blocks (labeled visually with a cart icon, a database cylinder, and a bank card) sitting side by side, glowing neon arrows flowing from the three blocks into a pipeline funnel on the right that ends at a monitor screen showing a tracing waterfall chart with teal and orange bars, dark navy background, subtle grid floor, soft lighting, cinematic depth, clean modern tech aesthetic, no readable text, 16:9, JPG under 100KB.

### Description
Lab microservices Golang + OpenTelemetry end-to-end: tiga service kecil (Order, Inventory, Payment) di-deploy ke k3d Kubernetes, semua trace melewati OpenTelemetry Collector dan divisualisasikan di Jaeger. Fokus pembelajaran: distributed tracing, context propagation, baggage, hingga analisis root cause dari sebuah transaksi terdistribusi.

### Tags
golang, opentelemetry, kubernetes, observability