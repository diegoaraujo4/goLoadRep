
FROM golang:1.22.5 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o goloadrep .
ENTRYPOINT  ["./goloadrep"]