package coordinator

import (
	"github.com/gorilla/mux"
	"net/http"
)

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
