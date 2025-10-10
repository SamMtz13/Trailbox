FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
# Compilamos solo el HTTP server
RUN go mod tidy && go build -o server ./services/users/cmd/http

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
