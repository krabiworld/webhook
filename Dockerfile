FROM rust:slim AS build
RUN apt update && apt install -y \
    pkg-config \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY . .
RUN cargo build --release && strip target/release/webhook

FROM gcr.io/distroless/cc-debian12
COPY --from=build /app/target/release/webhook .
CMD ["./webhook"]
