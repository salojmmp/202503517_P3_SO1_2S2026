FROM rust:1.93-bookworm AS builder

WORKDIR /app

COPY api-rust/Cargo.toml api-rust/Cargo.lock ./
COPY api-rust/src ./src

RUN cargo build --release

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/target/release/api-rust /app/api-rust

ENV RUST_API_ADDR=0.0.0.0:8000
ENV GO_API_URL=http://api-go:8080

EXPOSE 8000

CMD ["/app/api-rust"]
