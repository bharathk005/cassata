# Start of Selection
FROM golang:1.20-alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o cassata ./cmd/server/main.go

FROM gcr.io/distroless/static

WORKDIR /app

COPY --from=builder /app/cassata .

CMD ["./cassata"] 
