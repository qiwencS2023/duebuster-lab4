package main

import (
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
)

func main() {
	// usage: coordinator <addr1> <addr2> ...
	// where addr is the address of a storage server
	// parse arguments
	args := os.Args
	println(args)
	if len(args) < 4 {
		panic("usage: coordinator -p <coordinator port> -s <addr1> <addr2> ...")
	}

	// create a coordinator server
	cPort := args[2]
	storagePorts := args[4:]

	// create a coordinator server
	coordinatorServer := NewCoordinatorServerImpl(storagePorts...)

	// create a listener on the coordinator port
	lis, err := net.Listen("tcp", ":"+cPort)
	if err != nil {
		panic(err)
	}

	// create a grpc server
	grpcServer := grpc.NewServer()
	RegisterCoordinatorServiceServer(grpcServer, coordinatorServer)

	// start the grpc server
	go func() {
		grpcServer.Serve(lis)
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
