package coordinator

import (
	"encoding/json"
	"net/http"
)

import (
	"github.com/gorilla/mux"
)

// Define the server struct
type Server struct {
	Host        string `json:"host"`
	CommandPort int    `json:"command_port"`
}

// Define the table struct
type Table struct {
	Name       string            `json:"name"`
	Columns    map[string]string `json:"columns"`
	PrimaryKey string            `json:"primary_key"`
}

// Define the line struct
type Line struct {
	Table      string            `json:"table"`
	PrimaryKey string            `json:"primary_key"`
	Line       map[string]string `json:"line"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/createTable", createTable).Methods("POST")
	router.HandleFunc("/deleteTable", deleteTable).Methods("DELETE")
	router.HandleFunc("/deleteLine", deleteLine).Methods("DELETE")
	router.HandleFunc("/insertLine", insertLine).Methods("POST")
	router.HandleFunc("/updateLine", updateLine).Methods("PUT")
	http.ListenAndServe(":8000", router)
}

func register(w http.ResponseWriter, r *http.Request) {
	var servers []Server
	json.NewDecoder(r.Body).Decode(&servers)
	// implement your logic to register servers here
	json.NewEncoder(w).Encode(servers) // placeholder response
}

func createTable(w http.ResponseWriter, r *http.Request) {
	var table Table
	json.NewDecoder(r.Body).Decode(&table)
	// implement your logic to create table here
	json.NewEncoder(w).Encode(table) // placeholder response
}

func deleteTable(w http.ResponseWriter, r *http.Request) {
	var table Table
	json.NewDecoder(r.Body).Decode(&table)
	// implement your logic to delete table here
	json.NewEncoder(w).Encode(table) // placeholder response
}

func deleteLine(w http.ResponseWriter, r *http.Request) {
	var line Line
	json.NewDecoder(r.Body).Decode(&line)
	// implement your logic to delete line here
	json.NewEncoder(w).Encode(line) // placeholder response
}

func insertLine(w http.ResponseWriter, r *http.Request) {
	var line Line
	json.NewDecoder(r.Body).Decode(&line)
	// implement your logic to insert line here
	json.NewEncoder(w).Encode(line) // placeholder response
}

func updateLine(w http.ResponseWriter, r *http.Request) {
	var line Line
	json.NewDecoder(r.Body).Decode(&line)
	// implement your logic to update line here
	json.NewEncoder(w).Encode(line) // placeholder response
}
