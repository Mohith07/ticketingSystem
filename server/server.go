package main

import (
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "ticketingSystem/router"
)

var (
	port = flag.Int("port", 5001, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	ser := &ticketingSystemServer{apiCount: 0}
	log.Println("server starting..")
	pb.RegisterRouteGuideServer(grpcServer, ser)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Error while loading the service")
		return
	}
}
