FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/service ./cmd

FROM scratch
COPY --from=builder /bin/service /service
EXPOSE 8004
ENTRYPOINT ["/service"]
