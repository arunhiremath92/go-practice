build: cmd/api/main.go
	go build -o urlshortner ./cmd/api 
run: build
	./urlshortner
test: build
	go test ./...