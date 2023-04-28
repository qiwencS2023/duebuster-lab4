package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: storage <port>")
	}

	port, _ := strconv.Atoi(os.Args[1])

	server := &StorageServerImpl{}

	// register grpc server
	grpcServer := grpc.NewServer()
	RegisterStorageServer(grpcServer, server)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		// handle ctrl + c
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			<-c
			grpcServer.GracefulStop()
			os.Exit(0)
		}()

		// start grpc server
		log.Fatal(grpcServer.Serve(lis))
	}()

}
