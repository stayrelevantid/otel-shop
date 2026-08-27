# blogpost.md — Day 2

---

## Section 1 — Blog Post

### Post Title
Day 2 OTel-Shop: Shipping Three Go Services to k3d

### Slug
day-2-otel-shop-shipping-three-go-services-to-k3d

### URL
http://stayrelevant.id/blog/day-2-otel-shop-shipping-three-go-services-to-k3d

### Excerpt (Meta Description)
Day 2 of OTel-Shop: three Go services got real handlers and routes, packed into lean containers, and landed on the local k3d cluster. All green.

### Tags
otel-shop, golang, open-telemetry, distributed-tracing, observability, kubernetes, k3d, jaeger, postgresql, microservices, devops, cloud-native

### Cover Image Prompt
Flat vector illustration, wide 16:9 tech blog cover banner, three rounded cargo containers resting on a dock, each decorated with a small icon (shopping cart, database cylinder, bank card), a friendly crane lifting one container onto a platform shaped like a ship steering wheel, a small monitor nearby showing three green progress bars, clean minimal style, dark navy background with teal and orange accent colors, subtle isometric grid faded in the background suggesting Kubernetes, no text, crisp smooth shapes, JPG output compressed under 100KB.

### Content

## Where We Left Off

Yesterday was all about foundations: the cluster came up, the database got seeded, and a fake trace made it all the way to Jaeger. Nice view, empty rooms. Day 2's mission was simple to say out loud: put actual services inside that beautiful infrastructure.

## What Got Built Today

Here's the haul:

- Gave all three services their shapes and starter brains — models for checkout, inventory lookup, and payment, each wired into clean routes.
- Kept the front door consistent: every service answers `/health`, which doubles as Kubernetes' way of asking "you good?" (readiness and liveness probes are hooked up to it).
- Wrote multi-stage container recipes: a chunky builder stage does the compiling, then the final image starts from literally nothing and carries one small binary. Lean is lovely.
- Rolled the whole packaging flow into a single script: build with Podman, hand over to Docker, import into the cluster. One command instead of a ritual.
- Added deployment and service manifests so each app gets its own address inside the cluster, plus fixed door numbers on my laptop to reach them from outside.
- Finished the day poking everything with curl like an impatient customer. Everything answered politely.

A taste of what worked straight from the host:

- A checkout request returned a shiny new order id with status `paid`.
- Asking about product A123 reported ten units in stock.
- A payment call happily declared success.
- Sloppy requests (empty fields, negative numbers) got polite 400 rejections instead of mystery errors.

One honest note: the smarts are intentionally shallow right now. Inventory isn't really reading the database yet, and checkout doesn't call its friends. Those are the next episodes.

## The Container Shuffle

The day's plot twist happened at the handover between tools. Podman built the images perfectly, but when they crossed into the Docker side they picked up an extra `localhost/` prefix on their names. The cluster import step then looked at those names, shrugged, and refused everything.

Five minutes of head-scratching later, the fix was one extra retag step in the script. Classic. The lesson generalizes nicely though: when two tools meet at a boundary, naming conventions are where things go sideways — not the actual work.

## Lessons Learned

Four takeaways from today:

- **Start-from-nothing images are underrated.** A runtime with zero extras has almost nothing that can rot, need patching, or surprise you at 2am.
- **Automate the boring handoffs immediately.** The build-to-cluster chain only bit me once because it became a script right after — future me says thanks.
- **Poke from the outside.** Testing through the public ports like a real client tells you what actually matters; testing only from inside would have felt deceptively fine.
- **Health checks are cheap insurance.** Wiring probes took minutes and already keeps rollouts honest — a pod only counts as ready when it truly answers.

## Conclusion

Day 2 ends with three services running in the cluster, every endpoint answering correctly from my laptop, and the whole ship-and-deploy loop scripted. The foundation phase is officially behind us; from here on it's behavior, wiring, and telemetry.

Tomorrow the services start talking — first inventory meets the database for real, then checkout learns to orchestrate its neighbors. See you there.

Repo is here: https://github.com/stayrelevantid/otel-shop

---

## Section 2 — LinkedIn Post

**Hook**
Day 2 of the OTel-Shop challenge: three Go services are officially alive inside my local Kubernetes cluster 🚀

**Today's wins:**
• Real handlers, routes and models for Order, Inventory & Payment — plain stdlib, no frameworks
• Multi-stage builds: hefty compile stage, featherweight final image
• One-command pipeline: Podman build → Docker load → import to k3d → deploy
• Every endpoint tested green with classic curl, including graceful 400s on bad input

Small drama: images crossing from Podman to Docker picked up a `localhost/` name prefix and the cluster import wanted none of it. One retag later, peace restored 🙃

Next up: hooking inventory to PostgreSQL for real, then making checkout call its friends.

Blog: http://stayrelevant.id/blog/day-2-otel-shop-shipping-three-go-services-to-k3d
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
