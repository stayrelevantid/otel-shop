# syntax=docker/dockerfile:1
# Build context = repo root (module telemetry ikut via replace directive)
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY pkg/telemetry/ ./pkg/telemetry/
COPY services/payment/ ./services/payment/
WORKDIR /src/services/payment
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/payment ./cmd

FROM scratch
COPY --from=build /out/payment /payment
EXPOSE 18082
ENTRYPOINT ["/payment"]
