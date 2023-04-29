package main

import "os"

func main() {
	// usage: coordinator <addr1> <addr2> ...
	// where addr is the address of a storage server
	// parse arguments
	args := os.Args
	println(args)
	if len(args) < 2 {
		panic("usage: coordinator <addr1> <addr2> ...")
	}

	args = args[1:]

	storageServers := make([]string, len(args))
	for i, arg := range args {
		storageServers[i] = arg
	}

}
