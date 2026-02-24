FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -o cashregister .

FROM alpine:3.21
COPY --from=builder /app/cashregister /usr/local/bin/cashregister
WORKDIR /data
ENTRYPOINT ["cashregister"]
