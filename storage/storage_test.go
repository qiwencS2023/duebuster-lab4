package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"google.golang.org/grpc"
)

// call register
var database = &Database{
	Type:     "mysql",
	Host:     "localhost",
	Port:     3306,
	Database: "golab4",
	User:     "golab4",
	Password: "golab4",
}

var table = &Table{
	Name: "test_table",
	Columns: map[string]string{
		"id":          "int",
		"test_column": "varchar(255)",
	},
	PrimaryKey: "id",
}

func printPointsOnSuccess(testName string, points int) {
	fmt.Println("[TEST] Test ", testName, "succesfulPoints: ", points)
}

func TestStorageServerImpl_Register(t *testing.T) {
	// mock command line arguments
	client, cancelServer, err := startStorageServer("9000")
	t.Cleanup(cancelServer)

	// call register
	_, err = client.Register(context.Background(), &Database{
		Type:     "mysql",
		Host:     "localhost",
		Port:     3306,
		Database: "golab4",
		User:     "golab4",
		Password: "golab4",
	})

	if err != nil {
		t.Error(err)
	}

}

func startStorageServer(port string) (StorageClient, context.CancelFunc, error) {
	os.Args = []string{"storage", "-p", port}

	// run the server with a context
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		startStorageServerWithCtx(ctx)
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

func TestStorageServerImpl_CreateTable(t *testing.T) {
	client, cancelServer, err := startStorageServer("9000")
	t.Cleanup(cancelServer)
	_, err = client.Register(context.Background(), database)
	if err != nil {
		t.Error(err)
	}

	// call create table
	_, err = client.CreateTable(context.Background(), table)

	if err != nil {
		t.Error(err)
	}

	printPointsOnSuccess("TestStorageServerImpl_CreateTable", 10)

}

func TestStorageServerImpl_DeleteTable(t *testing.T) {
	client, cancelServer, err := startStorageServer("9000")
	t.Cleanup(cancelServer)
	_, err = client.Register(context.Background(), database)
	if err != nil {
		t.Error(err)
	}

	// call create table
	_, err = client.CreateTable(context.Background(), table)

	if err != nil {
		t.Error(err)
	}

	// call delete table
	_, err = client.DeleteTable(context.Background(), &Table{
		Name: "test_table",
	})
	if err != nil {
		t.Error(err)
	}

	printPointsOnSuccess("TestStorageServerImpl_DeleteTable", 10)
}

func TestStorageServerImpl_InsertLine(t *testing.T) {
	client, cancelServer, err := startStorageServer("9000")
	t.Cleanup(cancelServer)
	_, err = client.Register(context.Background(), database)
	if err != nil {
		t.Error(err)
	}

	// call create table
	_, err = client.CreateTable(context.Background(), table)
	if err != nil {
		t.Error(err)
	}

	// call insert line
	_, err = client.InsertLine(context.Background(), &Line{
		Table: "test_table",
		Line: map[string]string{
			"id":          "1",
			"test_column": "test_value",
		},
	})
	if err != nil {
		t.Error(err)
	}

	printPointsOnSuccess("TestStorageServerImpl_InsertLine", 10)
}

func TestStorageServerImpl_DeleteLine(t *testing.T) {
	client, cancelServer, err := startStorageServer("9000")
	t.Cleanup(cancelServer)
	_, err = client.Register(context.Background(), database)
	if err != nil {
		t.Error(err)
	}

	// call delete table
	_, err = client.DeleteTable(context.Background(), &Table{
		Name: "test_table",
	})
	if err != nil {
		t.Error(err)
	}

	// call create table
	_, err = client.CreateTable(context.Background(), table)
	if err != nil {
		t.Error(err)
	}

	// call insert line
	_, err = client.InsertLine(context.Background(), &Line{
		Table: "test_table",
		Line: map[string]string{
			"id":          "1",
			"test_column": "test_value",
		},
	})

	if err != nil {
		t.Error(err)
	}

	// call delete line
	_, err = client.DeleteLine(context.Background(), &Line{
		Table: "test_table",
		Line: map[string]string{
			"id": "1",
		},
		PrimaryKey: "id",
	})
	if err != nil {
		t.Error(err)
	}

	printPointsOnSuccess("TestStorageServerImpl_DeleteLine", 10)
}

func TestStorageServerImpl_UpdateLine(t *testing.T) {
	client, cancelServer, err := startStorageServer("9000")
	t.Cleanup(cancelServer)
	_, err = client.Register(context.Background(), database)
	if err != nil {
		t.Error(err)
	}

	// call delete table
	_, err = client.DeleteTable(context.Background(), &Table{
		Name: "test_table",
	})
	if err != nil {
		t.Error(err)
	}

	// call create table
	_, err = client.CreateTable(context.Background(), table)
	if err != nil {
		t.Error(err)
	}

	// call insert line
	_, err = client.InsertLine(context.Background(), &Line{
		Table: "test_table",
		Line: map[string]string{
			"id":          "1",
			"test_column": "test_value",
		},
	})

	if err != nil {
		t.Error(err)
	}

	// call update line
	_, err = client.UpdateLine(context.Background(), &Line{
		Table: "test_table",
		Line: map[string]string{
			"id":          "1",
			"test_column": "test_value_updated",
		},
		PrimaryKey: "id",
	})
	if err != nil {
		t.Error(err)
	}

	printPointsOnSuccess("TestStorageServerImpl_UpdateLine", 10)
}

func TestStorageServerImpl_GetLine(t *testing.T) {
	client, cancelServer, err := startStorageServer("9000")
	t.Cleanup(cancelServer)
	_, err = client.Register(context.Background(), database)
	if err != nil {
		t.Error(err)
	}

	// call delete table
	_, err = client.DeleteTable(context.Background(), &Table{
		Name: "test_table",
	})
	if err != nil {
		t.Error(err)
	}

	// call create table
	_, err = client.CreateTable(context.Background(), table)
	if err != nil {
		t.Error(err)
	}

	// call insert line
	_, err = client.InsertLine(context.Background(), &Line{
		Table: "test_table",
		Line: map[string]string{
			"id":          "1",
			"test_column": "test_value",
		},
	})

	if err != nil {
		t.Error(err)
	}

	fmt.Printf("table: %v\n", table)

	// call get line
	_, err = client.GetLine(context.Background(), &GetLineRequest{
		Table:           table,
		PrimaryKeyValue: "1",
	})

	if err != nil {
		cancelServer()
		t.Error(err)
	}

	printPointsOnSuccess("TestStorageServerImpl_GetLine", 10)
}
