FROM golang:alpine AS builder
WORKDIR /app
COPY . .
ENV CGO_ENABLED=0
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -ldflags="-s -w" -trimpath ./cmd/webhook

FROM alpine AS certs
RUN apk add --no-cache curl
RUN curl -sSL https://i.pki.goog/r4.pem -o /tmp/rootCA.pem

FROM scratch
COPY --from=builder /app/webhook .
COPY --from=certs /tmp/rootCA.pem /etc/ssl/certs/ca-certificates.crt
COPY <<EOF /etc/passwd
nobody:x:65534:65534:nobody:/:
EOF
USER nobody
CMD ["/webhook"]
