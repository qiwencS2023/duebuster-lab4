package main

import (
	"context"
	"google.golang.org/grpc"
	"os"
)

type StorageServerImpl struct {
	StorageClient
	cancelServer func()
}

type CoordinatorServerImpl struct {
	storageServers []*StorageServerImpl
}

func NewCoordinatorServerImpl(storagePorts ...string) *CoordinatorServerImpl {
	storageServers := make([]*StorageServerImpl, len(storagePorts))
	for i, port := range storagePorts {
		client, cancelServer, err := connectStorageServer(port)
		if err != nil {
			panic(err)
		}
		storageServers[i] = &StorageServerImpl{
			client,
			cancelServer,
		}
	}
	return &CoordinatorServerImpl{
		storageServers: storageServers,
	}
}

func connectStorageServer(port string) (StorageClient, context.CancelFunc, error) {
	os.Args = []string{"storage", port}

	// run the server with a context
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		go main()
		// listen for ctx done
		<-ctx.Done()
	}(ctx)

	// create a client
	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	// create a sub context

	// close connection when context is done
	go func(ctx context.Context) {
		<-ctx.Done()
		conn.Close()
	}(ctx)

	// create a storage client
	client := NewStorageClient(conn)
	return client, cancel, err
}
