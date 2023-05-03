package main

import (
	"context"
	"database/sql"
	"fmt"
	"google.golang.org/grpc"
	"os"
	"testing"
)

// call register
var database = &Database{
	Type:     "mysql",
	Host:     "localhost",
	Port:     3306,
	Database: "test",
	User:     "test",
	Password: "test",
}

var table = &Table{
	Name: "test_table",
	Columns: map[string]string{
		"id":          "int",
		"test_column": "varchar(255)",
	},
	PrimaryKey: "id",
}

func TestStorageServerImpl_Register(t *testing.T) {
	// mock command line arguments
	client, cancelServer, err := startStorageServer("9000")

	// call register
	_, err = client.Register(context.Background(), &Database{
		Type:     "mysql",
		Host:     "localhost",
		Port:     3306,
		Database: "test",
		User:     "test",
		Password: "test",
	})

	if err != nil {
		t.Error(err)
	}

	cancelServer()
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
	_, err = client.Register(context.Background(), database)
	if err != nil {
		t.Error(err)
	}

	// call create table
	_, err = client.CreateTable(context.Background(), table)

	if err != nil {
		t.Error(err)
	}

	cancelServer()
}

func TestStorageServerImpl_DeleteTable(t *testing.T) {
	client, cancelServer, err := startStorageServer("9000")
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

	cancelServer()
}

func TestStorageServerImpl_InsertLine(t *testing.T) {
	client, cancelServer, err := startStorageServer("9000")
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

	cancelServer()
}

func TestStorageServerImpl_DeleteLine(t *testing.T) {
	client, cancelServer, err := startStorageServer("9000")
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

	cancelServer()
}

func TestStorageServerImpl_UpdateLine(t *testing.T) {
	client, cancelServer, err := startStorageServer("9000")
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

	cancelServer()
}

func TestStorageServerImpl_GetLine(t *testing.T) {
	client, cancelServer, err := startStorageServer("9000")
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

	cancelServer()
}

func Test_GetLine(t *testing.T) {
	request := &GetLineRequest{
		Table: &Table{
			Name:       "test_table",
			PrimaryKey: "id",
		},
		PrimaryKeyValue: "1",
	}
	table := request.Table
	pk := request.PrimaryKeyValue

	// Create SQL command
	cmd := fmt.Sprintf("SELECT * FROM %s WHERE %s = %s", table.Name, table.PrimaryKey, pk)

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", "test", "test", "localhost", "test")
	db, err := sql.Open("mysql", dsn)

	// Execute SQL command
	rows, err := db.Query(cmd)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	// Get column names
	columns, err := rows.Columns()

	// Get column types
	types, err := rows.ColumnTypes()

	for rows.Next() {
		// Create column values
		columnValues := make(map[string]string)

		// Create column pointers
		columnPointers := make([]interface{}, len(columns))
		for i, _ := range columns {
			columnPointers[i] = new(interface{})
		}

		// Scan columns into pointers
		if err := rows.Scan(columnPointers...); err != nil {
			fmt.Printf("error: %v\n", err)
		}

		// Put column values into map
		for i, column := range columns {
			columnValue := columnPointers[i].(*interface{})
			columnType := types[i].DatabaseTypeName()
			fmt.Printf("column: %s, type: %s\n", column, columnType)
			switch columnType {
			case "VARCHAR", "TEXT":
				columnValues[column] = string((*columnValue).([]uint8))
			case "INT":
				columnValues[column] = string((*columnValue).([]uint8))
			case "FLOAT":
				columnValues[column] = string((*columnValue).([]uint8))
			default:
				return
			}
		}

		// Create line
		_ = &Line{
			Table:      table.Name,
			PrimaryKey: table.PrimaryKey,
			Line:       columnValues,
		}

		fmt.Printf("line: %v\n", columnValues)

	}

}
