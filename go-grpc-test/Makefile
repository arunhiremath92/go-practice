# Define your final executable names as variables
clean:
	rm -f ../grpc-client ../grpc-server

protodir: 
	mkdir -p src/client/pb
	mkdir -p src/server/pb

clientprotos: 
	rm -f src/client/pb/*
	protoc --go_out=src/client/pb --go_opt=paths=source_relative --proto_path=src/proto src/proto/helloworld.proto
	protoc --go-grpc_out=src/client/pb --go-grpc_opt=paths=source_relative --proto_path=src/proto src/proto/helloworld.proto

client: protodir clientprotos
	go mod -C src tidy
	go -C src/client build -o ../grpc-client  ./ 

server: protodir
	protoc --go_out=src/server/pb --go_opt=paths=source_relative --proto_path=src/proto src/proto/helloworld.proto
	protoc --go-grpc_out=src/server/pb --go-grpc_opt=paths=source_relative --proto_path=src/proto src/proto/helloworld.proto
	go mod -C src tidy
	go -C src/server build -o ../grpc-server  ./


