FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -ldflags="-s -w" -trimpath

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/webhook .
CMD ["/webhook"]
