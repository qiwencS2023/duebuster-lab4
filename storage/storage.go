package main

import (
	"fmt"
	"golang.org/x/net/context"
)

type StorageServerImpl struct {
	dbConnector DBConnector
}

func (s *StorageServerImpl) mustEmbedUnimplementedStorageServer() {
	panic("implement me")
}

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

	return &StorageResponse{}, nil
}

func (s *StorageServerImpl) CreateTable(ctx context.Context, table *Table) (*StorageResponse, error) {
	if err := s.dbConnector.CreateTable(table); err != nil {
		return nil, err
	}

	return &StorageResponse{}, nil
}

func (s *StorageServerImpl) DeleteTable(ctx context.Context, table *Table) (*StorageResponse, error) {
	if err := s.dbConnector.DeleteTable(table); err != nil {
		return nil, err
	}

	return &StorageResponse{}, nil
}

func (s *StorageServerImpl) InsertLine(ctx context.Context, line *Line) (*StorageResponse, error) {
	if err := s.dbConnector.InsertLine(line); err != nil {
		return nil, err
	}

	return &StorageResponse{}, nil
}

func (s *StorageServerImpl) DeleteLine(ctx context.Context, line *Line) (*StorageResponse, error) {
	if err := s.dbConnector.DeleteLine(line); err != nil {
		return nil, err
	}

	return &StorageResponse{}, nil
}

func (s *StorageServerImpl) UpdateLine(ctx context.Context, line *Line) (*StorageResponse, error) {
	if err := s.dbConnector.UpdateLine(line); err != nil {
		return nil, err
	}

	return &StorageResponse{}, nil
}

func (s *StorageServerImpl) GetLine(ctx context.Context, request *GetLineRequest) (*Line, error) {
	if line, err := s.dbConnector.GetLine(request); err != nil {
		return nil, err
	} else {
		return line, nil
	}
}
