package main

import (
	"encoding/json"
	"fmt"
	"log"
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

type DBConnector interface {
	Connect(user, password, host, dbname string) error
	Disconnect() error
	CreateTable(table Table) error
	DeleteTable(table Table) error
	DeleteLine(line Line) error
	InsertLine(line Line) error
	UpdateLine(line Line) error
}

var dbConnector DBConnector

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

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var db Database
	if err := json.NewDecoder(r.Body).Decode(&db); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch db.Type {
	case "mysql":
		dbConnector = &MySQLConnector{}
	case "Cassandra":
		dbConnector = &CassandraConnector{}
	default:
		http.Error(w, "Invalid database type", http.StatusBadRequest)
		return
	}

	if err := dbConnector.Connect(db.User, db.Password, db.Host, db.Database); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleCreateTable(w http.ResponseWriter, r *http.Request) {
	var table Table
	if err := json.NewDecoder(r.Body).Decode(&table); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := dbConnector.CreateTable(table); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleDeleteTable(w http.ResponseWriter, r *http.Request) {
	var table Table
	if err := json.NewDecoder(r.Body).Decode(&table); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := dbConnector.DeleteTable(table); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleInsertLine(w http.ResponseWriter, r *http.Request) {
	var line Line
	if err := json.NewDecoder(r.Body).Decode(&line); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := dbConnector.InsertLine(line); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleDeleteLine(w http.ResponseWriter, r *http.Request) {
	var line Line
	if err := json.NewDecoder(r.Body).Decode(&line); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := dbConnector.DeleteLine(line); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleUpdateLine(w http.ResponseWriter, r *http.Request) {
	var line Line
	if err := json.NewDecoder(r.Body).Decode(&line); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := dbConnector.UpdateLine(line); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
