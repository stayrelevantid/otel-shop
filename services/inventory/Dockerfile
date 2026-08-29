# syntax=docker/dockerfile:1
# Build context = repo root (module telemetry ikut via replace directive)
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY pkg/telemetry/ ./pkg/telemetry/
COPY services/inventory/ ./services/inventory/
WORKDIR /src/services/inventory
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/inventory ./cmd

FROM scratch
COPY --from=build /out/inventory /inventory
EXPOSE 18081
ENTRYPOINT ["/inventory"]
