FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/leaderboard/ ./services/leaderboard/

WORKDIR /app/services/leaderboard
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /service ./cmd

FROM gcr.io/distroless/static-debian12
COPY --from=builder /service /service
EXPOSE 8007
ENTRYPOINT ["/service"]
