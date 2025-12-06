build: main.go
	gofmt -w main.go main_test.go
	go build -o file_sorter main.go
run : build
	./file_sorter
test: build
	go test ./...