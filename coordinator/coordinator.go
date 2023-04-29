package main

type StorageServerImpl struct {
	StorageServer
}

type CoordinatorServerImpl struct {
	storageServers []*StorageServerImpl
}
