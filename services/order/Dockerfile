# syntax=docker/dockerfile:1
# Build context = repo root (module telemetry ikut via replace directive)
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY pkg/telemetry/ ./pkg/telemetry/
COPY services/order/ ./services/order/
WORKDIR /src/services/order
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/order ./cmd

FROM scratch
COPY --from=build /out/order /order
EXPOSE 18080
ENTRYPOINT ["/order"]
