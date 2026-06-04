FROM golang:latest AS builder
WORKDIR /urlshortner/

# Copy dependency files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY internal/ internal/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux  go build -o urlshortner ./cmd/api

# Use a minimal image for the final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /urlshortner/urlshortner .

CMD ["/root/urlshortner"]
EXPOSE 8000