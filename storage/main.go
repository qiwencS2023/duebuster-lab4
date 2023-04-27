package storage

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
)

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
