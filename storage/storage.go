package storage

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
)

// Define the Database struct
type Database struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// Define the Table struct (same as before)
type Table struct {
	Name       string            `json:"name"`
	Columns    map[string]string `json:"columns"`
	PrimaryKey string            `json:"primary_key"`
}

// Define the Line struct (same as before)
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

	port, _ := strconv.Atoi(os.Args[1])
	http.ListenAndServe(":"+strconv.Itoa(port), router)
}

func register(w http.ResponseWriter, r *http.Request) {
	var db Database
	json.NewDecoder(r.Body).Decode(&db)
	// implement your logic to register database here
	json.NewEncoder(w).Encode(db) // placeholder response
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
	json.NewEncoder(w).Encode(table) // placeholder
}

func deleteLine(w http.ResponseWriter, r *http.Request) {
	var line Line
	json.NewDecoder(r.Body).Decode(&line)
	// implement your logic to delete line here
	json.NewEncoder(w).Encode(line) // placeholder
}

func insertLine(w http.ResponseWriter, r *http.Request) {
	var line Line
	json.NewDecoder(r.Body).Decode(&line)
	// implement your logic to insert line here
	json.NewEncoder(w).Encode(line) // placeholder
}

func updateLine(w http.ResponseWriter, r *http.Request) {
	var line Line
	json.NewDecoder(r.Body).Decode(&line)
	// implement your logic to update line here
	json.NewEncoder(w).Encode(line) // placeholder
}
