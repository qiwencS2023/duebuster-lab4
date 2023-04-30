package main

import (
	context "context"
	"errors"
	"math/rand"
	"time"

	"google.golang.org/grpc"
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
	// find the table partition
	tableName := line.Table
	CoordinatorTable := c.CoordinatorTableMap[tableName]

	// select the partition with lower count
	var minRowCountTablePartition TablePartition
	var partitionIdx int
	for idx, tablePartition := range CoordinatorTable.TablePartitions {
		if idx == 0 {
			minRowCountTablePartition = tablePartition
			partitionIdx = idx
		} else {
			if tablePartition.RowCount < minRowCountTablePartition.RowCount {
				minRowCountTablePartition = tablePartition
				partitionIdx = idx
			}
		}
	}

	// insert the line into the partition
	paritionLine := &Line{
		Table: minRowCountTablePartition.Name,
		PrimaryKey: line.PrimaryKey,
		Line: line.Line,
	}
	minRowCountTablePartition.StorageServer.InsertLine(context.Background(), paritionLine)

	// insert the line into the replication
	replicationLine := &Line{
		Table: CoordinatorTable.ReplicationPartitions[partitionIdx].Name,
		PrimaryKey: line.PrimaryKey,
		Line: line.Line,
	}
	CoordinatorTable.ReplicationPartitions[partitionIdx].StorageServer.InsertLine(context.Background(), replicationLine);

	// update the row count
	CoordinatorTable.TablePartitions[partitionIdx].RowCount++
	CoordinatorTable.ReplicationPartitions[partitionIdx].RowCount++

	return &CoordinatorResponse{
		Message: "success",
	}, nil

}

func (c *CoordinatorServerImpl) DeleteLine(ctx context.Context, line *Line) (*CoordinatorResponse, error) {
	tableName := line.Table
	CoordinatorTable := c.CoordinatorTableMap[tableName]

	// delete the line from all the partition
	for _, tablePartition := range CoordinatorTable.TablePartitions {
		request := &Line{
			Table: tablePartition.Name,
			PrimaryKey: line.PrimaryKey,
		}
		tablePartition.StorageServer.DeleteLine(context.Background(), request)
	}

	for _, tablePartition := range CoordinatorTable.ReplicationPartitions {
		request := &Line{
			Table: tablePartition.Name,
			PrimaryKey: line.PrimaryKey,
		}
		tablePartition.StorageServer.DeleteLine(context.Background(), request)
	}

	return &CoordinatorResponse{
		Message: "success",
	}, nil

}

func (c *CoordinatorServerImpl) GetLine(ctx context.Context, lineRequest *GetLineRequest) (*Line, error) {
	// find the table partition
	tableName := lineRequest.Table.Name
	CoordinatorTable := c.CoordinatorTableMap[tableName]

	// send the query to all the partitions
	for _, tablePartition := range CoordinatorTable.TablePartitions {
		request := &GetLineRequest{
			Table : &Table{
				Name: tablePartition.Name,
			},
			PrimaryKeyValue: lineRequest.PrimaryKeyValue,
		}
		line, err := tablePartition.StorageServer.GetLine(context.Background(), request)
		if err == nil {
			return line, nil
		}
	}

	return nil, errors.New("line not found")
}

func (c *CoordinatorServerImpl) UpdateLine(ctx context.Context, line *Line) (*CoordinatorResponse, error) {
	tableName := line.Table
	CoordinatorTable := c.CoordinatorTableMap[tableName]

	// delete the line from all the partition
	for _, tablePartition := range CoordinatorTable.TablePartitions {
		request := &Line{
			Table: tablePartition.Name,
			PrimaryKey: line.PrimaryKey,
		}
		tablePartition.StorageServer.UpdateLine(context.Background(), request)
	}

	for _, tablePartition := range CoordinatorTable.ReplicationPartitions {
		request := &Line{
			Table: tablePartition.Name,
			PrimaryKey: line.PrimaryKey,
		}
		tablePartition.StorageServer.UpdateLine(context.Background(), request)
	}

	return &CoordinatorResponse{
		Message: "success",
	}, nil
}

func (c *CoordinatorServerImpl) mustEmbedUnimplementedCoordinatorServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (c *CoordinatorServerImpl) RegisterStorageServer(request *RegisterRequest) {

	conn, err := grpc.Dial(request.StorageAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	storageServer := &StorageServerImpl{NewStorageClient(conn), nil}
	c.SSM.storageServers[request.StorageAddr] = storageServer
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
