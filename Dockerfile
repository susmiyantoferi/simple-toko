FROM golang:1.24.6-bookworm AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-w" -o goapp .

FROM debian:bookworm-slim
WORKDIR /app

COPY --from=builder /app/goapp .

CMD ["./goapp" ]