build: main.go
	go build -o simpleserver main.go
test: build
	go test ./...
run: build test
	./simpleserver
