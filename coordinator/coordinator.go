package main

import (
	context "context"
	"google.golang.org/grpc"
	"math/rand"
	"time"
)

type StorageServerImpl struct {
	StorageClient
	stop func()
}

type StorageServerMap struct {
	storageServers map[string]*StorageServerImpl // port -> storage server
}

type TablePartition struct {
	StorageServer *StorageServerImpl
	Name          string
	RowCount      int
}

type CoordinatorTable struct {
	TablePartitions       []TablePartition
	ReplicationPartitions []TablePartition
}

type CoordinatorServerImpl struct {
	SSM                 StorageServerMap
	CoordinatorTableMap map[string]*CoordinatorTable
}

func (c *CoordinatorServerImpl) CreateTable(ctx context.Context, request *CreateTableRequest) (*CoordinatorResponse, error) {
	table := request.Table
	partitionCount := int(request.PartitionCount)
	tableName := table.Name

	// create a table in coordinator's view
	// a new table will be paritioned into 2 parts and replicate 1 times.
	ct := &CoordinatorTable{}
	// we do a default replication factor of 2
	ct.TablePartitions = make([]TablePartition, partitionCount)
	ct.ReplicationPartitions = make([]TablePartition, partitionCount)

	c.CoordinatorTableMap[tableName] = ct

	// randomly assign storage servers to the table partitions
	paritions := RandomPartitions(c.SSM.storageServers, partitionCount*2)

	for i := 0; i < partitionCount; i++ {
		TablePartition := TablePartition{
			StorageServer: paritions[i],
			Name:          tableName + "_partition_" + string(i),
			RowCount:      0,
		}
		ct.TablePartitions = append(ct.TablePartitions, TablePartition)
		// create actual table in storage server
		table := &Table{
			Name:    TablePartition.Name,
			Columns: table.Columns,
		}
		TablePartition.StorageServer.CreateTable(context.Background(), table)
	}

	for i := 0; i < partitionCount; i++ {
		TablePartition := TablePartition{
			StorageServer: paritions[i+partitionCount],
			Name:          tableName + "_partition_replica_" + string(i),
			RowCount:      0,
		}
		ct.ReplicationPartitions = append(ct.ReplicationPartitions, TablePartition)
		// create actual table in storage server
		table := &Table{
			Name:    TablePartition.Name,
			Columns: table.Columns,
		}
		TablePartition.StorageServer.CreateTable(context.Background(), table)
	}

	// return success
	return &CoordinatorResponse{
		Message: "success",
	}, nil
}

func (c *CoordinatorServerImpl) DeleteTable(ctx context.Context, table *Table) (*CoordinatorResponse, error) {
	tableName := table.Name
	// delete a table in coordinator's view
	CoordinatorTable := c.CoordinatorTableMap[tableName]
	// delete the underlying table partitions
	for _, tablePartition := range CoordinatorTable.TablePartitions {
		table := &Table{
			Name: tablePartition.Name,
		}
		tablePartition.StorageServer.DeleteTable(context.Background(), table)
	}

	// delete the underlying replication partitions
	for _, tablePartition := range CoordinatorTable.ReplicationPartitions {
		table := &Table{
			Name: tablePartition.Name,
		}
		tablePartition.StorageServer.DeleteTable(context.Background(), table)
	}

	delete(c.CoordinatorTableMap, tableName)
	return &CoordinatorResponse{
		Message: "success",
	}, nil
}

func (c *CoordinatorServerImpl) InsertLine(ctx context.Context, line *Line) (*CoordinatorResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CoordinatorServerImpl) DeleteLine(ctx context.Context, line *Line) (*CoordinatorResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CoordinatorServerImpl) GetLine(ctx context.Context, line *Line) (*Line, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CoordinatorServerImpl) UpdateLine(ctx context.Context, line *Line) (*CoordinatorResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CoordinatorServerImpl) mustEmbedUnimplementedCoordinatorServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (c *CoordinatorServerImpl) RegisterStorageServer() {
	// read from config file
	// implement me
	panic("implement me")
}

// randomly get 4 storage servers
func RandomPartitions(m map[string]*StorageServerImpl, n int) []*StorageServerImpl {
	keys := make([]*StorageServerImpl, 0, len(m))
	for _, value := range m {
		keys = append(keys, value)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	return keys[:n]
}

func (c *CoordinatorServerImpl) GetCoordinatorTable(tableName string) *CoordinatorTable {
	return c.CoordinatorTableMap[tableName]
}

func (c *CoordinatorServerImpl) init() {
	c.SSM.storageServers = make(map[string]*StorageServerImpl)
	c.CoordinatorTableMap = make(map[string]*CoordinatorTable)
}

func NewCoordinatorServerImpl(storagePorts ...string) *CoordinatorServerImpl {
	coordinatorServer := &CoordinatorServerImpl{}
	coordinatorServer.init()
	for _, storagePort := range storagePorts {
		// create client stub
		conn, err := grpc.Dial("localhost"+storagePort, grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		coordinatorServer.SSM.storageServers[storagePort] = &StorageServerImpl{NewStorageClient(conn), nil}
	}
	return coordinatorServer
}
