FROM golang:1.26-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /out && go build -o /out/metrics-exporter ./metrics-exporter

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /out/metrics-exporter /usr/local/bin/metrics-exporter

EXPOSE 9100

ENTRYPOINT ["/usr/local/bin/metrics-exporter"]
