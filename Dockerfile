FROM golang:1.24.2-alpine AS builder
WORKDIR /app

COPY ./go.mod /app/
COPY ./go.sum /app/
RUN go mod download

COPY ./app /app/app/
RUN go build app/cmd/auth/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main /app/main
ENTRYPOINT [ "/app/main" ]


