package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"google.golang.org/grpc"
)

var table = &Table{
	Name: "test_table",
	Columns: map[string]string{
		"id":          "int",
		"test_column": "varchar(255)",
	},
	PrimaryKey: "id",
}

var createTableRequest = &CreateTableRequest{
	Table:          table,
	PartitionCount: 2,
}

func cleanup() {
	// kill process on port 8999 to 9004
	for i := 9001; i < 9005; i++ {
		port := strconv.Itoa(i)
		cmd := exec.Command("sh", "-c", "lsof -i :"+port+" | grep LISTEN | awk '{print $2}' | xargs kill")
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func printPointsOnSuccess(testName string, points int) {
	fmt.Println("[TEST] Test ", testName, "succesfulPoints: ", points)
}

func startStorageServer(port string) (StorageClient, context.CancelFunc, error) {
	target := "../dist/storage -p " + port

	stopChan := make(chan bool)
	go func() {
		cmd := exec.CommandContext(context.Background(), "sh", "-c", target)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		done := make(chan error, 1)

		go func() {
			done <- cmd.Run()
		}()

		select {
		case <-stopChan:
			fmt.Println("[subroutine] Stopping storage server on port " + port)
			if err := cmd.Process.Kill(); err != nil {
				fmt.Println("failed to kill process: ", err)
			}
		case err := <-done:
			log.Printf("[subroutine] process done with error = %v", err)
			os.Exit(0)
		}
	}()

	time.Sleep(200 * time.Millisecond)

	// create a client
	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	// create a storage client
	client := NewStorageClient(conn)

	// register database
	_, err = client.Register(context.Background(), &Database{
		Type:     "mysql",
		Host:     "localhost",
		Port:     3306,
		Database: "golab4",
		Password: "golab4",
		User:     "golab4",
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Started storage server on port %s\n", port)
	cancel := func() {
		fmt.Printf("[startStorageServer] Stopping storage server on port %s\n", port)
		close(stopChan)
	}
	return client, cancel, err
}

func mockStorageCluster(numStorage int) ([]string, context.CancelFunc, error) {
	storageServers := make([]*StorageServerImpl, numStorage)
	for i := 0; i < numStorage; i++ {
		port := strconv.Itoa(9001 + i)
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
		ports[i] = strconv.Itoa(9001 + i)
	}
	return ports, func() {
		for _, server := range storageServers {
			server.stop()
		}
	}, nil
}

func mockCoordinator(port string, storagePorts []string) (context.CancelFunc, error) {
	os.Args = []string{"", "-p", port, "-s"}
	os.Args = append(os.Args, storagePorts...)

	// run the server with a context
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		startCoordinatorServer(ctx)
	}(ctx)

	// create a client
	// conn, err := grpc.Dial("localhost:8999", grpc.WithInsecure())
	// if err != nil {
	// 	panic(err)
	// }

	// create a sub context

	// close connection when context is done
	go func(ctx context.Context) {
		<-ctx.Done()
		// conn.Close()
	}(ctx)

	// create a storage client
	// client := NewCoordinatorServiceClient(conn)
	return cancel, nil
}

func ClientFactory(addr string) (CoordinatorServiceClient, error) {
	// create a client
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	// create a storage client
	client := NewCoordinatorServiceClient(conn)
	return client, err
}

func TestCoordinatorServerImpl_CreateTable(t *testing.T) {
	t.Cleanup(cleanup)

	// mock storage cluster
	ports, cancelStorage, err := mockStorageCluster(4)
	// mock coordinator
	cancelServer, err := mockCoordinator("8999", ports)
	cClient, err := ClientFactory("localhost:8999")
	if err != nil {
		t.Error(err)
	}

	resp, err := cClient.CreateTable(context.Background(), &CreateTableRequest{
		Table:          table,
		PartitionCount: 2,
	})

	if err != nil {
		t.Errorf("CreateTable() error = %v", err)
	}

	t.Logf("response: %v", resp)
	cancelServer()
	cancelStorage()
	printPointsOnSuccess("TestCoordinatorServerImpl_CreateTable", 20)
}

func TestCoordinatorServerImpl_DeleteLine(t *testing.T) {
	t.Cleanup(cleanup)

	// mock storage cluster
	ports, cancelStorage, err := mockStorageCluster(4)
	defer cancelStorage()

	// mock coordinator
	cancelServer, err := mockCoordinator("8999", ports)
	cClient, err := ClientFactory("localhost:8999")
	if err != nil {
		t.Error(err)
	}
	defer cancelServer()

	resp, err := cClient.CreateTable(context.Background(), createTableRequest)
	if err != nil {
		t.Errorf("CreateTable() error = %v", err)
	}

	line := &Line{
		Table: table.Name,
		Line: map[string]string{
			"id":          "1",
			"test_column": "test_1",
		},
	}
	resp, err = cClient.InsertLine(context.Background(), line)
	if err != nil {
		t.Errorf("InsertLine() error = %v", err)
	}

	resp, err = cClient.DeleteLine(context.Background(), &Line{
		PrimaryKey: "id",
		Line: map[string]string{
			"id": "1",
		},
		Table: table.Name,
	})
	if err != nil {
		t.Errorf("DeleteLine() error = %v", err)
	}

	// cleanup
	resp, err = cClient.DeleteTable(context.Background(), table)
	if err != nil {
		t.Errorf("DeleteTable() error = %v", err)
	}

	t.Logf("response: %v", resp)
	printPointsOnSuccess("TestCoordinatorServerImpl_DeleteLine", 20)
}

func TestCoordinatorServerImpl_DeleteTable(t *testing.T) {
	t.Cleanup(cleanup)

	// mock storage cluster
	ports, cancelStorage, err := mockStorageCluster(4)
	defer cancelStorage()

	// mock coordinator
	cancelServer, err := mockCoordinator("8999", ports)
	cClient, err := ClientFactory("localhost:8999")
	if err != nil {
		t.Error(err)
	}
	defer cancelServer()

	resp, err := cClient.CreateTable(context.Background(), createTableRequest)
	if err != nil {
		t.Errorf("CreateTable() error = %v", err)
	}

	resp, err = cClient.DeleteTable(context.Background(), table)
	if err != nil {
		t.Errorf("DeleteTable() error = %v", err)
	}

	t.Logf("response: %v", resp)
	printPointsOnSuccess("TestCoordinatorServerImpl_DeleteTable", 20)
}

func TestCoordinatorServerImpl_GetLine(t *testing.T) {
	t.Cleanup(cleanup)

	// mock storage cluster
	ports, cancelStorage, err := mockStorageCluster(4)
	defer cancelStorage()

	// mock coordinator
	cancelServer, err := mockCoordinator("8999", ports)
	cClient, err := ClientFactory("localhost:8999")
	if err != nil {
		t.Error(err)
	}
	defer cancelServer()

	_, err = cClient.CreateTable(context.Background(), createTableRequest)
	if err != nil {
		t.Errorf("CreateTable() error = %v", err)
	}

	// insert line
	line := &Line{
		Table: table.Name,
		Line: map[string]string{
			"id":          "1",
			"test_column": "test",
		},
	}
	_, err = cClient.InsertLine(context.Background(), line)
	if err != nil {
		t.Errorf("InsertLine() error = %v", err)
	}

	// get line
	newLine, err := cClient.GetLine(context.Background(), &GetLineRequest{
		Table: &Table{
			Name:       table.Name,
			PrimaryKey: "id",
		},
		PrimaryKeyValue: "1",
	})
	if err != nil {
		t.Errorf("GetLine() error = %v", err)
	}

	// cleanup
	_, err = cClient.DeleteTable(context.Background(), table)
	if err != nil {
		t.Errorf("DeleteTable() error = %v", err)
	}

	t.Logf("response: %v", newLine)
	t.Cleanup(cleanup)
	printPointsOnSuccess("TestCoordinatorServerImpl_GetLine", 20)
}

func TestCoordinatorServerImpl_InsertLine(t *testing.T) {
	t.Cleanup(cleanup)

	// mock storage cluster
	ports, cancelStorage, err := mockStorageCluster(4)

	defer cancelStorage()

	// mock coordinator
	cancelServer, err := mockCoordinator("8999", ports)
	cClient, err := ClientFactory("localhost:8999")
	if err != nil {
		t.Error(err)
	}

	defer cancelServer()

	resp, err := cClient.CreateTable(context.Background(), createTableRequest)
	if err != nil {
		t.Errorf("CreateTable() error = %v", err)
	}

	// insert line
	for i := 0; i < 100; i++ {
		line := &Line{
			Table: table.Name,
			Line: map[string]string{
				"id":          fmt.Sprintf("%d", i),
				"test_column": "test_" + fmt.Sprintf("%d", i),
			},
		}
		resp, err = cClient.InsertLine(context.Background(), line)
		if err != nil {
			t.Errorf("InsertLine() error = %v", err)
		}
	}

	// cleanup
	resp, err = cClient.DeleteTable(context.Background(), table)
	if err != nil {
		t.Errorf("DeleteTable() error = %v", err)
	}

	t.Logf("response: %v", resp)
	printPointsOnSuccess("TestCoordinatorServerImpl_InsertLine", 20)
}

func TestCoordinatorServerImpl_UpdateLine(t *testing.T) {
	t.Cleanup(cleanup)

	// mock storage cluster
	ports, cancelStorage, err := mockStorageCluster(4)
	defer cancelStorage()

	// mock coordinator
	cancelServer, err := mockCoordinator("8999", ports)
	cClient, err := ClientFactory("localhost:8999")
	if err != nil {
		t.Error(err)
	}
	defer cancelServer()

	resp, err := cClient.CreateTable(context.Background(), createTableRequest)
	if err != nil {
		t.Errorf("CreateTable() error = %v", err)
	}

	// insert line
	line := &Line{
		Table: table.Name,
		Line: map[string]string{
			"id":          "1",
			"test_column": "test",
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
			"id":          "1",
			"test_column": "test2",
		},
		PrimaryKey: "id",
	})
	if err != nil {
		t.Errorf("UpdateLine() error = %v", err)
	}

	// cleanup
	resp, err = cClient.DeleteTable(context.Background(), table)
	if err != nil {
		t.Errorf("DeleteTable() error = %v", err)
	}

	t.Logf("response: %v", resp)
	printPointsOnSuccess("TestCoordinatorServerImpl_UpdateLine", 20)
}

func TestMultipleClientCanInsertIntoDB(t *testing.T) {
	t.Cleanup(cleanup)

	// mock storage cluster
	ports, cancelStorage, err := mockStorageCluster(4)
	defer cancelStorage()

	// mock coordinator
	cancelServer, err := mockCoordinator("8999", ports)
	cClient, err := ClientFactory("localhost:8999")
	cClient2, err := ClientFactory("localhost:8999")
	if err != nil {
		t.Error(err)
	}
	defer cancelServer()

	resp, err := cClient.CreateTable(context.Background(), createTableRequest)
	if err != nil {
		t.Errorf("CreateTable() error = %v", err)
	}

	// insert line
	for i := 0; i < 100; i++ {
		line := &Line{
			Table: table.Name,
			Line: map[string]string{
				"id":          fmt.Sprintf("%d", i),
				"test_column": "test_" + fmt.Sprintf("%d", i),
			},
		}
		if i%2 == 0 {
			resp, err = cClient.InsertLine(context.Background(), line)
			if err != nil {
				t.Errorf("InsertLine() error = %v", err)
			}
		} else {
			resp, err = cClient2.InsertLine(context.Background(), line)
			if err != nil {
				t.Errorf("InsertLine() error = %v", err)
			}
		}
	}

	// check lines are presented from both client perspective
	for i := 0; i < 100; i++ {
		line, err := cClient.GetLine(context.Background(), &GetLineRequest{
			Table: &Table{
				Name:       table.Name,
				PrimaryKey: "id",
			},
			PrimaryKeyValue: fmt.Sprintf("%d", i),
		})
		if err != nil {
			t.Errorf("GetLine() error = %v", err)
		}

		if line.Line["test_column"] != "test_"+fmt.Sprintf("%d", i) {
			t.Errorf("GetLine() error = %v", err)
		}
	}

	// cleanup
	resp, err = cClient.DeleteTable(context.Background(), table)
	if err != nil {
		t.Errorf("DeleteTable() error = %v", err)
	}

	t.Logf("response: %v", resp)
	printPointsOnSuccess("TestMultipleClientCanInsertIntoDB", 20)
}
