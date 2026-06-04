FROM golang:1.24
WORKDIR /build
COPY ./ ./
RUN go mod download
EXPOSE 8080
RUN GOOS=linux GOARCH=amd64 go build -o /build/route hmac/api/route.go
CMD ["/build/route"]