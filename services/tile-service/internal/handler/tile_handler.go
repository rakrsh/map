// Copyright (c) 2026 Ravi Sharma

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func HealthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func GetTile(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "tile request",
		"z":       vars["z"],
		"x":       vars["x"],
		"y":       vars["y"],
	})
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", HealthCheck).Methods(http.MethodGet)
	r.HandleFunc("/tiles/{z}/{x}/{y}", GetTile).Methods(http.MethodGet)
	return r
}
