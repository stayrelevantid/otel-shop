# syntax=docker/dockerfile:1
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/payment ./cmd

FROM scratch
COPY --from=build /out/payment /payment
EXPOSE 18082
ENTRYPOINT ["/payment"]
