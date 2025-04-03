FROM golang:1.23 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o loadtest ./cmd/loadtest/main.go

FROM scratch

WORKDIR /root/

COPY --from=builder /app/loadtest .

ENTRYPOINT ["./loadtest"]
