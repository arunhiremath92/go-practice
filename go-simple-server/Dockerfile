FROM golang:tip-alpine3.22 AS builder

WORKDIR /app

COPY go.mod /app/

RUN go mod download

COPY * /app/


RUN CGO_ENABLED=0 GOOS=linux go build -o simple-http-server main.go

EXPOSE 6000

CMD ["/app/simple-http-server"]