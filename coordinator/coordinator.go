package main

import (
	context "context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	sync "sync"
	"time"

	"google.golang.org/grpc"
)

// StorageCommand is the command line arguments for the storage server
type StorageServerImpl struct {
	StorageClient
	stop func()
}

// StorageServerMap is a map of storage servers, indexed by port
type StorageServerMap struct {
	storageServers map[string]*StorageServerImpl // port -> storage server
}

// TablePartition is a partition of a table, stored on a storage server
type TablePartition struct {
	StorageServer *StorageServerImpl
	Name          string
	RowCount      int
}

// CoordinatorTable is a table in the coordinator's view, which is partitioned and replicated
type CoordinatorTable struct {
	TablePartitions       []*TablePartition
	ReplicationPartitions []*TablePartition
	PrimaryKey            string
	lock                  sync.Mutex
}

// CoordinatorServerImpl is the coordinator server, which is responsible for managing the storage servers
type CoordinatorServerImpl struct {
	CoordinatorCache
	SSM                 StorageServerMap
	CoordinatorTableMap map[string]*CoordinatorTable
}

// CreateTable creates a table in the coordinator's view
// It will randomly assign storage servers to the table partitions
func (c *CoordinatorServerImpl) CreateTable(ctx context.Context, request *CreateTableRequest) (*CoordinatorResponse, error) {
	fmt.Println("[Coordinator] Received CreateTable request: ", request)

	table := request.Table
	partitionCount := int(request.PartitionCount)
	tableName := table.Name

	// create a table in coordinator's view
	// a new table will be paritioned into 2 parts and replicate 1 times.
	ct := &CoordinatorTable{}

	// acquire lock
	ct.lock.Lock()
	defer ct.lock.Unlock()

	// we do a default replication factor of 2
	ct.TablePartitions = []*TablePartition{}
	ct.ReplicationPartitions = []*TablePartition{}
	ct.PrimaryKey = table.PrimaryKey

	c.CoordinatorTableMap[tableName] = ct

	// randomly assign storage servers to the table partitions
	paritions := RandomPartitions(c.SSM.storageServers, partitionCount*2)

	fmt.Printf("[Coordinator] Randomly assigned partitions: %v\n", paritions)

	for i := 0; i < partitionCount; i++ {
		TablePartition := TablePartition{
			StorageServer: paritions[i],
			Name:          tableName + "_partition_" + strconv.Itoa(i),
			RowCount:      0,
		}
		ct.TablePartitions = append(ct.TablePartitions, &TablePartition)

		// create actual table in storage server
		table := &Table{
			Name:       TablePartition.Name,
			Columns:    table.Columns,
			PrimaryKey: table.PrimaryKey,
		}

		_, err := TablePartition.StorageServer.CreateTable(context.Background(), table)
		if err != nil {
			return nil, fmt.Errorf("[coordinator] error creating table: %v", err)
		}
	}

	for i := 0; i < partitionCount; i++ {
		TablePartition := TablePartition{
			StorageServer: paritions[i+partitionCount],
			Name:          tableName + "_partition_replica_" + strconv.Itoa(i),
			RowCount:      0,
		}
		ct.ReplicationPartitions = append(ct.ReplicationPartitions, &TablePartition)
		// create actual table in storage server
		table := &Table{
			Name:       TablePartition.Name,
			Columns:    table.Columns,
			PrimaryKey: table.PrimaryKey,
		}
		_, err := TablePartition.StorageServer.CreateTable(context.Background(), table)
		if err != nil {
			return nil, err
		}
	}

	// put ct to the coordinator table map
	c.CoordinatorTableMap[tableName] = ct
	log.Printf("[coordinator] coordinator table map: %v\n", c.CoordinatorTableMap[tableName])

	// return success
	return &CoordinatorResponse{
		Message: "success",
	}, nil
}

// DeleteTable deletes a table in the coordinator's view,
// and also deletes the underlying table partitions and replication partitions
func (c *CoordinatorServerImpl) DeleteTable(ctx context.Context, table *Table) (*CoordinatorResponse, error) {
	tableName := table.Name
	// delete a table in coordinator's view
	CoordinatorTable := c.CoordinatorTableMap[tableName]
	// acquire lock
	CoordinatorTable.lock.Lock()
	defer CoordinatorTable.lock.Unlock()
	// delete the underlying table partitions
	for _, tablePartition := range CoordinatorTable.TablePartitions {
		table := &Table{
			Name:       tablePartition.Name,
			Columns:    map[string]string{},
			PrimaryKey: "",
		}
		log.Printf("[coordinator] deleting table partition %v\n", tablePartition)
		_, err := tablePartition.StorageServer.DeleteTable(context.Background(), table)
		if err != nil {
			return nil, err
		}
	}

	// delete the underlying replication partitions
	for _, tablePartition := range CoordinatorTable.ReplicationPartitions {
		table := &Table{
			Name: tablePartition.Name,
		}
		_, err := tablePartition.StorageServer.DeleteTable(context.Background(), table)
		if err != nil {
			return nil, err
		}
	}

	delete(c.CoordinatorTableMap, tableName)
	return &CoordinatorResponse{
		Message: "success",
	}, nil
}

// InsertLine inserts a line into a table in the coordinator's view
// It will randomly and by default insert into the partition with lower count
// to achieve load balancing
func (c *CoordinatorServerImpl) InsertLine(ctx context.Context, line *Line) (*CoordinatorResponse, error) {
	// find the table partition
	tableName := line.Table
	CoordinatorTable := c.CoordinatorTableMap[tableName]
	// acquire lock
	CoordinatorTable.lock.Lock()
	defer CoordinatorTable.lock.Unlock()

	// select the partition with lower count
	var minRowCountTablePartition TablePartition
	var partitionIdx int
	for idx, tablePartition := range CoordinatorTable.TablePartitions {
		if idx == 0 {
			minRowCountTablePartition = *tablePartition
			partitionIdx = idx
		} else {
			if tablePartition.RowCount < minRowCountTablePartition.RowCount {
				minRowCountTablePartition = *tablePartition
				partitionIdx = idx
			}
		}
	}

	// insert the line into the partition
	paritionLine := &Line{
		Table:      minRowCountTablePartition.Name,
		PrimaryKey: line.PrimaryKey,
		Line:       line.Line,
	}
	_, err := minRowCountTablePartition.StorageServer.InsertLine(context.Background(), paritionLine)
	if err != nil {
		return nil, err
	}

	// insert the line into the replication
	replicationLine := &Line{
		Table:      CoordinatorTable.ReplicationPartitions[partitionIdx].Name,
		PrimaryKey: line.PrimaryKey,
		Line:       line.Line,
	}
	_, err = CoordinatorTable.ReplicationPartitions[partitionIdx].StorageServer.InsertLine(context.Background(), replicationLine)
	if err != nil {
		return nil, err
	}

	// update the row count
	CoordinatorTable.TablePartitions[partitionIdx].RowCount++
	CoordinatorTable.ReplicationPartitions[partitionIdx].RowCount++

	// cache the result
	c.PutCache(&GetLineRequest{
		Table: &Table{
			Name:       tableName,
			PrimaryKey: CoordinatorTable.PrimaryKey,
		},
		PrimaryKeyValue: line.Line[CoordinatorTable.PrimaryKey],
	}, replicationLine)

	return &CoordinatorResponse{
		Message: "success",
	}, nil

}

// DeleteLine deletes a line from a table in the coordinator's view
// It will delete the line from all the partitions and replication partitions
func (c *CoordinatorServerImpl) DeleteLine(ctx context.Context, line *Line) (*CoordinatorResponse, error) {
	tableName := line.Table
	CoordinatorTable := c.CoordinatorTableMap[tableName]
	// acquire lock
	CoordinatorTable.lock.Lock()
	defer CoordinatorTable.lock.Unlock()

	// delete the line from all the partition
	for _, tablePartition := range CoordinatorTable.TablePartitions {
		request := &Line{
			Table:      tablePartition.Name,
			PrimaryKey: line.PrimaryKey,
			Line:       line.Line,
		}
		log.Printf("[coordinator] deleting line from table partition %v\n", tablePartition)
		_, err := tablePartition.StorageServer.DeleteLine(context.Background(), request)
		if err != nil {
			return nil, err
		}
	}

	for _, tablePartition := range CoordinatorTable.ReplicationPartitions {
		request := &Line{
			Table:      tablePartition.Name,
			PrimaryKey: line.PrimaryKey,
			Line:       line.Line,
		}
		_, err := tablePartition.StorageServer.DeleteLine(context.Background(), request)
		if err != nil {
			return nil, err
		}
	}

	// invalidate the cache
	c.InvalidateCache(tableName, line.Line[line.PrimaryKey])

	return &CoordinatorResponse{
		Message: "success",
	}, nil

}

// GetLine gets a line from a table in the coordinator's view
// It will first check the cache, if not found, it will send the query to all the partitions
func (c *CoordinatorServerImpl) GetLine(ctx context.Context, lineRequest *GetLineRequest) (*Line, error) {
	// find the table partition
	tableName := lineRequest.Table.Name
	CoordinatorTable := c.CoordinatorTableMap[tableName]
	// acquire lock
	CoordinatorTable.lock.Lock()
	defer CoordinatorTable.lock.Unlock()

	// check if the line is in the cache
	if line, ok := c.GetCache(lineRequest); ok {
		log.Printf("[coordinator] cache hit for %v\n", lineRequest)
		return line, nil
	}

	// send the query to all the partitions
	for _, tablePartition := range CoordinatorTable.TablePartitions {
		request := &GetLineRequest{
			Table: &Table{
				Name:       tablePartition.Name,
				PrimaryKey: lineRequest.Table.PrimaryKey,
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

// UpdateLine updates a line in a table in the coordinator's view
// It will update the line in all the partitions and replication partitions
func (c *CoordinatorServerImpl) UpdateLine(ctx context.Context, line *Line) (*CoordinatorResponse, error) {
	tableName := line.Table
	CoordinatorTable := c.CoordinatorTableMap[tableName]

	// delete the line from all the partition
	for _, tablePartition := range CoordinatorTable.TablePartitions {
		request := &Line{
			Table:      tablePartition.Name,
			PrimaryKey: line.PrimaryKey,
			Line:       line.Line,
		}
		_, err := tablePartition.StorageServer.UpdateLine(context.Background(), request)
		if err != nil {
			return nil, err
		}
	}

	for _, tablePartition := range CoordinatorTable.ReplicationPartitions {
		request := &Line{
			Table:      tablePartition.Name,
			PrimaryKey: line.PrimaryKey,
			Line:       line.Line,
		}
		_, err := tablePartition.StorageServer.UpdateLine(context.Background(), request)
		if err != nil {
			return nil, err
		}
	}

	c.PutCache(&GetLineRequest{
		Table: &Table{
			Name:       tableName,
			PrimaryKey: CoordinatorTable.PrimaryKey,
		},
		PrimaryKeyValue: line.Line[CoordinatorTable.PrimaryKey],
	}, line)

	return &CoordinatorResponse{
		Message: "success",
	}, nil
}

// mustEmbedUnimplementedCoordinatorServiceServer is a "must implementation" function for the gRPC server
func (c *CoordinatorServerImpl) mustEmbedUnimplementedCoordinatorServiceServer() {
	//TODO implement me
	panic("implement me")
}

// RegisterStorageServer registers a storage server to the coordinator
// It will add the storage server to the storage server map
func (c *CoordinatorServerImpl) RegisterStorageServer(ctx context.Context, request *RegisterRequest) (*CoordinatorResponse, error) {

	conn, err := grpc.Dial(request.StorageAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	storageServer := &StorageServerImpl{NewStorageClient(conn), nil}
	c.SSM.storageServers[request.StorageAddr] = storageServer

	return &CoordinatorResponse{
		Message: "success",
	}, nil
}

// RandomPartitions returns n random partitions from the storage server map
// It will be used to select n random partitions for replication
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

// GetCoordinatorTable returns a coordinator table from the coordinator table map
// It will be used to get the coordinator table from the coordinator table map
func (c *CoordinatorServerImpl) GetCoordinatorTable(tableName string) *CoordinatorTable {
	return c.CoordinatorTableMap[tableName]
}

// NewCoordinatorServerImpl creates a new coordinator server
// It will create a new coordinator server with a new coordinator cache and a new storage server map
func NewCoordinatorServerImpl(storagePorts ...string) *CoordinatorServerImpl {
	coordinatorServer := &CoordinatorServerImpl{
		NewCoordinatorCache(),
		StorageServerMap{make(map[string]*StorageServerImpl)},
		make(map[string]*CoordinatorTable),
	}
	for _, storagePort := range storagePorts {
		// create client stub
		conn, err := grpc.Dial("localhost:"+storagePort, grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		coordinatorServer.SSM.storageServers[storagePort] = &StorageServerImpl{
			NewStorageClient(conn),
			nil,
		}
	}
	return coordinatorServer
}
