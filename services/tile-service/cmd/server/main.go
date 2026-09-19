package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}).Methods(http.MethodGet)

	r.HandleFunc("/tiles/{z}/{x}/{y}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(fmt.Sprintf("{\"message\":\"tile request\",\"z\":\"%s\",\"x\":\"%s\",\"y\":\"%s\"}", vars["z"], vars["x"], vars["y"])))
	}).Methods(http.MethodGet)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("tile-service listening on %s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Printf("server exited: %v", err)
		os.Exit(1)
	}
}
