FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /service ./services/map/cmd/grpcmain

FROM gcr.io/distroless/static-debian12
COPY --from=builder /service /service

EXPOSE 50051
ENTRYPOINT ["/service"]
