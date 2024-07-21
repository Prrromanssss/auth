FROM golang:1.22.5-alpine AS builder

COPY app /github.com/Prrromanssss/auth/source
WORKDIR /github.com/Prrromanssss/auth/source

RUN go mod download
RUN go build -o ./bin/auth cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /github.com/Prrromanssss/auth/source/bin/auth .

ADD config.yaml /config.yaml

CMD ["./auth"]