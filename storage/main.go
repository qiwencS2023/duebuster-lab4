package main

import (
	"os"
)

func main() {
	// detect port argument
	if len(os.Args) < 2 {
		println("Usage: storage <port>")
		return
	}

}
