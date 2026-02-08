FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates wget
WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080

ENV PORT=8080
ENV DATABASE_URL=file:/data/dayzsmartcf.db

VOLUME /data

CMD ["./server"]
