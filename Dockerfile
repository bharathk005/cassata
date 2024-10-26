# Start of Selection
FROM golang:1.23-alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o cassata ./cmd/server/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -o cassata-init ./init/init.go

FROM gcr.io/distroless/static

WORKDIR /app

COPY --from=builder /app/cassata .

COPY --from=builder /app/cassata-init .

CMD ["./cassata"] 
