FROM golang:1.24-alpine AS builder
LABEL stage=builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o calculator ./cmd/restcalculator


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/calculator ./

ENTRYPOINT [ "./calculator" ]