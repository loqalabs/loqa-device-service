FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY loqa-device-service/go.mod loqa-device-service/go.sum ./
RUN go mod download

COPY loqa-device-service/ .
RUN go build -o device-service ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/device-service .

ENV NATS_URL=nats://nats:4222
ENV LOG_LEVEL=info

CMD ["./device-service"]