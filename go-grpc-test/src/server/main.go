package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/arunhiremath/go-grpc-test/server/pb"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedGreeterServer
	pb.UnimplementedDataHandlerServer
}

func (srv *server) SayHello(ctxt context.Context, helloPkt *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("received hello from the user ", helloPkt.Name)
	fmt.Println("sending a reply back to the messenger")
	return &pb.HelloReply{
		Name: "server",
	}, nil
}

func (srv *server) SendACK(ctxt context.Context, helloPkt *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("received hello from the user ", helloPkt.Name)
	fmt.Println("sending a reply back to the messenger")
	return &pb.HelloReply{
		Name: "server",
	}, nil
}

func main() {
	var srv server
	flag.Parse()
	listner, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterDataHandlerServer(grpcServer, &srv)
	pb.RegisterGreeterServer(grpcServer, &srv)
	log.Printf("server listening at %v", listner.Addr())
	if err := grpcServer.Serve(listner); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
