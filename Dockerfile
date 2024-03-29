FROM golang:1.21-alpine AS builder
WORKDIR /qon
COPY go.mod .
RUN go mod download && go mod verify
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN apk add --no-cache make
COPY . .
RUN make build

FROM alpine:3.16
COPY --from=builder /qon/qon /
COPY --from=builder /qon/migrations /migrations
COPY --from=builder /qon/templates  /templates
USER nobody
CMD ["/qon"]
