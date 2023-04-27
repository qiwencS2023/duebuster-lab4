package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: storage <port>")
	}
	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Invalid command line argument: ", err)
	}

	http.HandleFunc("/register", handleRegister)
	http.HandleFunc("/createTable", handleCreateTable)
	http.HandleFunc("/deleteTable", handleDeleteTable)
	http.HandleFunc("/insertLine", handleInsertLine)
	http.HandleFunc("/deleteLine", handleDeleteLine)
	http.HandleFunc("/updateLine", handleUpdateLine)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
