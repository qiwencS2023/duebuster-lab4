package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
)

type StorageCommand struct {
	// [-d <database> -dp <data port> -h <host> -u <user> -pw <password>]
	Port     int
	Database string
	DataPort int
	Host     string
	User     string
	Password string
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: storage -p <port> [-d <database> -dp <data port> -h <host> -u <user> -pw <password>]")
	}

	// parse all command line arguments
	var command StorageCommand
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-p":
			i++
			command.Port, _ = strconv.Atoi(os.Args[i])
		case "-d":
			i++
			command.Database = os.Args[i]
		case "-dp":
			i++
			command.DataPort, _ = strconv.Atoi(os.Args[i])
		case "-h":
			i++
			command.Host = os.Args[i]
		case "-u":
			i++
			command.User = os.Args[i]
		case "-pw":
			i++
			command.Password = os.Args[i]
		}
	}

	server := &StorageServerImpl{}

	// register grpc server
	grpcServer := grpc.NewServer()
	RegisterStorageServer(grpcServer, server)

	// if all arguments are set, register database
	if command.Database != "" && command.DataPort != 0 && command.Host != "" && command.User != "" && command.Password != "" {
		// register database
		_, err := server.Register(nil, &Database{
			Type:     "mysql",
			Host:     command.Host,
			Port:     int32(command.DataPort),
			Database: command.Database,
			User:     command.User,
			Password: command.Password,
		})

		if err != nil {
			log.Fatal(err)
		}
	}

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", command.Port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		fmt.Printf("storage[storage] Listening on port %d\n\tregistration info: %s:%d %s %s %s\n", command.Port, command.Host, command.DataPort, command.Database, command.User, command.Password)

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

	// wait for ctrl + c
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
}

func startStorageServerWithCtx(ctx context.Context) {
	if len(os.Args) < 2 {
		log.Fatal("Usage: storage -p <port> [-d <database> -dp <data port> -h <host> -u <user> -pw <password>]")
	}

	// parse all command line arguments
	var command StorageCommand
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-p":
			i++
			command.Port, _ = strconv.Atoi(os.Args[i])
		case "-d":
			i++
			command.Database = os.Args[i]
		case "-dp":
			i++
			command.DataPort, _ = strconv.Atoi(os.Args[i])
		case "-h":
			i++
			command.Host = os.Args[i]
		case "-u":
			i++
			command.User = os.Args[i]
		case "-pw":
			i++
			command.Password = os.Args[i]
		}
	}

	server := &StorageServerImpl{}

	// register grpc server
	grpcServer := grpc.NewServer()
	RegisterStorageServer(grpcServer, server)

	// if all arguments are set, register database
	if command.Database != "" && command.DataPort != 0 && command.Host != "" && command.User != "" && command.Password != "" {
		// register database
		_, err := server.Register(nil, &Database{
			Type:     "mysql",
			Host:     command.Host,
			Port:     int32(command.DataPort),
			Database: command.Database,
			User:     command.User,
			Password: command.Password,
		})

		if err != nil {
			log.Fatal(err)
		}
	}

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", command.Port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		fmt.Printf("storage[storage] Listening on port %d\n\tregistration info: %s:%d %s %s %s\n", command.Port, command.Host, command.DataPort, command.Database, command.User, command.Password)
		// start grpc server
		log.Print(grpcServer.Serve(lis))
	}()

	select {
	case <-ctx.Done():
		grpcServer.GracefulStop()
		return
	}
}
