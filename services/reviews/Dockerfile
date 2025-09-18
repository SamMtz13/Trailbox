FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/reviews/ ./services/reviews/

WORKDIR /app/services/reviews
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /service ./cmd

FROM gcr.io/distroless/static-debian12
COPY --from=builder /service /service
EXPOSE 8004
ENTRYPOINT ["/service"]
