FROM golang:alpine AS build
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o webhook

FROM gcr.io/distroless/static-debian12
COPY --from=build /app/webhook .
CMD ["./webhook"]
