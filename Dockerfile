FROM golang:1.21-alpine AS builder
WORKDIR /qon
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -v -o qon ./cmd/main.go

FROM alpine:3.16
COPY --from=builder /qon/qon /
COPY --from=builder /qon/migrations /migrations
USER nobody
CMD ["/qon"]