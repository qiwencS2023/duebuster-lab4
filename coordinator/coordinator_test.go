package main

import (
	"context"
	"google.golang.org/grpc"
	"os"
	"strconv"
	"testing"
)

var table = &Table{
	Name: "test_table",
	Columns: map[string]string{
		"id":          "int",
		"test_column": "varchar(255)",
	},
	PrimaryKey: "id",
}

func startStorageServer(port string) (StorageClient, context.CancelFunc, error) {
	os.Args = []string{"storage", port}

	// run the server with a context
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		go main()
		// listen for ctx done
		<-ctx.Done()
	}(ctx)

	// create a client
	conn, err := grpc.Dial("localhost"+port, grpc.WithInsecure())
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

func mockStorageCluster(numStorage int) ([]string, context.CancelFunc, error) {
	storageServers := make([]*StorageServerImpl, numStorage)
	for i := 0; i < numStorage; i++ {
		port := strconv.Itoa(9000 + i)
		client, cancelServer, err := startStorageServer(port)
		if err != nil {
			panic(err)
		}
		storageServers[i] = &StorageServerImpl{
			client,
			cancelServer,
		}
	}
	ports := make([]string, numStorage)
	for i := 0; i < numStorage; i++ {
		ports[i] = strconv.Itoa(9000 + i)
	}
	return ports, func() {
		for _, server := range storageServers {
			server.cancelServer()
		}
	}, nil
}

func mockCoordinator(port string, storagePorts []string) (CoordinatorServiceClient, context.CancelFunc, error) {
	os.Args = []string{"", "-p", "localhost:" + port, "-s"}
	os.Args = append(os.Args, storagePorts...)

	// run the server with a context
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		go main()
		// listen for ctx done
		<-ctx.Done()
	}(ctx)

	// create a client
	conn, err := grpc.Dial("localhost:8999", grpc.WithInsecure())
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
	client := NewCoordinatorServiceClient(conn)
	return client, cancel, err
}

func TestCoordinatorServerImpl_CreateTable(t *testing.T) {
	// mock storage cluster
	ports, cancelStorage, err := mockStorageCluster(3)
	defer cancelStorage()

	// mock coordinator
	cClient, cancelServer, err := mockCoordinator("8999", ports)
	if err != nil {
		t.Error(err)
	}
	defer cancelServer()

	resp, err := cClient.CreateTable(context.Background(), table)
	if err != nil {
		t.Errorf("CreateTable() error = %v", err)
	}

	t.Logf("response: %v", resp)
}

func TestCoordinatorServerImpl_DeleteLine(t *testing.T) {
	// mock storage cluster
	ports, cancelStorage, err := mockStorageCluster(3)
	defer cancelStorage()

	// mock coordinator
	cClient, cancelServer, err := mockCoordinator("8999", ports)
	if err != nil {
		t.Error(err)
	}
	defer cancelServer()

	resp, err := cClient.CreateTable(context.Background(), table)
	if err != nil {
		t.Errorf("CreateTable() error = %v", err)
	}

	t.Logf("response: %v", resp)
}

func TestCoordinatorServerImpl_DeleteTable(t *testing.T) {
	// mock storage cluster
	ports, cancelStorage, err := mockStorageCluster(3)
	defer cancelStorage()

	// mock coordinator
	cClient, cancelServer, err := mockCoordinator("8999", ports)
	if err != nil {
		t.Error(err)
	}
	defer cancelServer()

	resp, err := cClient.CreateTable(context.Background(), table)
	if err != nil {
		t.Errorf("CreateTable() error = %v", err)
	}

	resp, err = cClient.DeleteTable(context.Background(), table)
	if err != nil {
		t.Errorf("DeleteTable() error = %v", err)
	}

	t.Logf("response: %v", resp)
}

func TestCoordinatorServerImpl_GetLine(t *testing.T) {
	// mock storage cluster
	ports, cancelStorage, err := mockStorageCluster(3)
	defer cancelStorage()

	// mock coordinator
	cClient, cancelServer, err := mockCoordinator("8999", ports)
	if err != nil {
		t.Error(err)
	}
	defer cancelServer()

	_, err = cClient.CreateTable(context.Background(), table)
	if err != nil {
		t.Errorf("CreateTable() error = %v", err)
	}

	// insert line
	line := &Line{
		Table: table.Name,
		Line: map[string]string{
			"id":   "1",
			"name": "test",
		},
	}
	_, err = cClient.InsertLine(context.Background(), line)
	if err != nil {
		t.Errorf("InsertLine() error = %v", err)
	}

	// get line
	newLine, err := cClient.GetLine(context.Background(), &Line{
		Table: table.Name,
		Line: map[string]string{
			"id": "1",
		},
	})
	if err != nil {
		t.Errorf("GetLine() error = %v", err)
	}

	t.Logf("response: %v", newLine)
}

func TestCoordinatorServerImpl_InsertLine(t *testing.T) {
	// mock storage cluster
	ports, cancelStorage, err := mockStorageCluster(3)

	defer cancelStorage()

	// mock coordinator
	cClient, cancelServer, err := mockCoordinator("8999", ports)
	if err != nil {
		t.Error(err)
	}

	defer cancelServer()

	resp, err := cClient.CreateTable(context.Background(), table)
	if err != nil {
		t.Errorf("CreateTable() error = %v", err)
	}

	// insert line
	line := &Line{
		Table: table.Name,
		Line: map[string]string{
			"id":   "1",
			"name": "test",
		},
	}
	resp, err = cClient.InsertLine(context.Background(), line)
	if err != nil {
		t.Errorf("InsertLine() error = %v", err)
	}

	t.Logf("response: %v", resp)

}

func TestCoordinatorServerImpl_UpdateLine(t *testing.T) {
	// mock storage cluster
	ports, cancelStorage, err := mockStorageCluster(3)
	defer cancelStorage()

	// mock coordinator
	cClient, cancelServer, err := mockCoordinator("8999", ports)
	if err != nil {
		t.Error(err)
	}
	defer cancelServer()

	resp, err := cClient.CreateTable(context.Background(), table)
	if err != nil {
		t.Errorf("CreateTable() error = %v", err)
	}

	// insert line
	line := &Line{
		Table: table.Name,
		Line: map[string]string{
			"id":   "1",
			"name": "test",
		},
	}
	resp, err = cClient.InsertLine(context.Background(), line)
	if err != nil {
		t.Errorf("InsertLine() error = %v", err)
	}

	// update line
	resp, err = cClient.UpdateLine(context.Background(), &Line{
		Table: table.Name,
		Line: map[string]string{
			"id":   "1",
			"name": "test2",
		},
	})
	if err != nil {
		t.Errorf("UpdateLine() error = %v", err)
	}

	t.Logf("response: %v", resp)
}
