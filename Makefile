.PHONY: all clean help grpc-client grpc-server lru mgclient pgclient pg-test secure-token shorturl tcp-chat website-status

# Default target
all: grpc-client grpc-server lru mgclient pgclient pg-test secure-token shorturl tcp-chat website-status

# Directory for build outputs
BIN_DIR := bin

# Create bin directory
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# Individual targets
grpc-client: | $(BIN_DIR)
	go build -o $(BIN_DIR)/grpc-client ./go-grpc-test/src/client

grpc-server: | $(BIN_DIR)
	go build -o $(BIN_DIR)/grpc-server ./go-grpc-test/src/server

lru: | $(BIN_DIR)
	go build -o $(BIN_DIR)/go-lru ./go-lru/cmd

mgclient: | $(BIN_DIR)
	go build -o $(BIN_DIR)/mgclient ./go-messenger/mongodb-client/client

pgclient: | $(BIN_DIR)
	go build -o $(BIN_DIR)/pgclient ./go-messenger/postgres-client/client

pg-test: | $(BIN_DIR)
	go build -o $(BIN_DIR)/go-pg-test ./go-pg-test

secure-token: | $(BIN_DIR)
	go build -o $(BIN_DIR)/secure-token ./go-sercure-token/hmac/api/route.go

shorturl: | $(BIN_DIR)
	go build -o $(BIN_DIR)/urlshortner ./go-shorturl/cmd/api

tcp-chat: | $(BIN_DIR)
	go build -o $(BIN_DIR)/go-tcp-chat ./go-tcp-chat

website-status: | $(BIN_DIR)
	go build -o $(BIN_DIR)/go-website-status ./go-website-status

# Clean target
clean:
	rm -rf $(BIN_DIR)

# Help target to display available targets
help:
	@echo "Available Makefile targets:"
	@echo "  all             - Build all projects into the 'bin' directory"
	@echo "  clean           - Remove the 'bin' directory"
	@echo "  grpc-client     - Build gRPC client"
	@echo "  grpc-server     - Build gRPC server"
	@echo "  lru             - Build LRU Cache"
	@echo "  mgclient        - Build Messenger MongoDB Client"
	@echo "  pgclient        - Build Messenger Postgres Client"
	@echo "  pg-test         - Build Postgres Test application"
	@echo "  secure-token    - Build Secure Token API"
	@echo "  shorturl        - Build URL Shortener API"
	@echo "  tcp-chat        - Build TCP Chat application"
	@echo "  website-status  - Build Website Status checker"
