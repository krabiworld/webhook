FROM lukemathwalker/cargo-chef:latest-rust-1.86.0 AS chef
WORKDIR /app

FROM chef AS planner
COPY . .
RUN cargo chef prepare --recipe-path recipe.json

FROM chef AS builder
COPY --from=planner /app/recipe.json recipe.json
RUN cargo chef cook --release --recipe-path recipe.json
COPY . .
RUN cargo build --release --bin webhook

FROM gcr.io/distroless/cc-debian12
COPY --from=builder /app/target/release/webhook .
CMD ["./webhook"]
