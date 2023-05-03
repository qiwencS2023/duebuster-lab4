package main

import (
	"fmt"
	"golang.org/x/net/context"
)

// StorageServerImpl is an implementation of the storage server
type StorageServerImpl struct {
	dbConnector DBConnector
}

func (s *StorageServerImpl) mustEmbedUnimplementedStorageServer() {
	panic("implement me")
}

// DBConnector is an interface for database connectors, each database connector must implement this interface
type DBConnector interface {
	Connect(user, password, host, dbname string) error
	Disconnect() error
	CreateTable(table *Table) error
	DeleteTable(table *Table) error
	DeleteLine(line *Line) error
	InsertLine(line *Line) error
	UpdateLine(line *Line) error
	GetLine(request *GetLineRequest) (*Line, error)
}

// Register registers a database to the storage server, the database is identified by the type
func (s *StorageServerImpl) Register(ctx context.Context, db *Database) (*StorageResponse, error) {
	// db type
	switch db.Type {
	case "mysql":
		s.dbConnector = &MySQLConnector{}
	case "cassandra":
		s.dbConnector = &CassandraConnector{}
	default:
		return nil, fmt.Errorf("invalid database type, only mysql and cassandra are supported")
	}

	if err := s.dbConnector.Connect(db.User, db.Password, db.Host, db.Database); err != nil {
		return nil, err
	}

	fmt.Printf("[storage] registered database, db: %v with type %s\n", db, db.Type)

	return &StorageResponse{}, nil
}

// CreateTable creates a table in the database, the table is identified by the name
func (s *StorageServerImpl) CreateTable(ctx context.Context, table *Table) (*StorageResponse, error) {
	fmt.Printf("[storage] creating table, table: %v\n", table)
	if err := s.dbConnector.CreateTable(table); err != nil {
		return nil, err
	}

	return &StorageResponse{}, nil
}

// DeleteTable deletes a table from the database, the table is identified by the name
func (s *StorageServerImpl) DeleteTable(ctx context.Context, table *Table) (*StorageResponse, error) {
	fmt.Printf("[storage] deleting table, table: %v\n", table)
	if err := s.dbConnector.DeleteTable(table); err != nil {
		return nil, err
	}

	return &StorageResponse{}, nil
}

// InsertLine inserts a line into the database
func (s *StorageServerImpl) InsertLine(ctx context.Context, line *Line) (*StorageResponse, error) {
	fmt.Printf("[storage] inserting line, line: %v\n", line)
	if err := s.dbConnector.InsertLine(line); err != nil {
		return nil, err
	}

	return &StorageResponse{}, nil
}

// DeleteLine deletes a line from the database, the line is identified by the primary key
// the return value is a line
func (s *StorageServerImpl) DeleteLine(ctx context.Context, line *Line) (*StorageResponse, error) {
	fmt.Printf("[storage] deleting line, line: %v\n", line)
	if err := s.dbConnector.DeleteLine(line); err != nil {
		return nil, err
	}

	return &StorageResponse{}, nil
}

// UpdateLine updates a line in the database, the line is identified by the primary key
// the line is updated with the values in the line parameter
func (s *StorageServerImpl) UpdateLine(ctx context.Context, line *Line) (*StorageResponse, error) {
	fmt.Printf("[storage] updating line, line: %v\n", line)
	if err := s.dbConnector.UpdateLine(line); err != nil {
		return nil, err
	}

	return &StorageResponse{}, nil
}

// GetLine returns a line from the database, the line is identified by the primary key
// the return value is a line
func (s *StorageServerImpl) GetLine(ctx context.Context, request *GetLineRequest) (*Line, error) {
	fmt.Printf("[storage] getting line, request: %v\n", request)
	if line, err := s.dbConnector.GetLine(request); err != nil {
		return nil, err
	} else {
		return line, nil
	}
}
