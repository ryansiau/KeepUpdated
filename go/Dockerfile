FROM golang:1.24-alpine AS builder

RUN apk add --no-cache --update build-base

WORKDIR /app

COPY . /app

ENV CGO_ENABLED=1
RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates libc6-compat libgcc libstdc++

WORKDIR /app

COPY --from=builder /app/main /app/main

ENTRYPOINT ["./main", "/etc/keep-updated/config.yaml"]
