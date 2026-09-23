// Copyright (c) 2026 Ravi Sharma

package handler

import (
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// HealthCheck returns a simple JSON health response with CORS enabled for dev.
func HealthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// GetTile is retained for API compatibility and returns a JSON placeholder.
func GetTile(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "tile request",
		"z":       vars["z"],
		"x":       vars["x"],
		"y":       vars["y"],
	})
}

// GetTilePNG generates a simple 256x256 PNG tile for development demos.
// The tile color varies deterministically by z/x/y so tiles look different.
func GetTilePNG(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")

	z, _ := strconv.Atoi(vars["z"])
	x, _ := strconv.Atoi(vars["x"])
	y, _ := strconv.Atoi(vars["y"])

	// Create a simple image and fill with a color derived from z/x/y
	const size = 256
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	// Derive color components
	r := uint8((z*37 + x*13) % 256)
	g := uint8((x*17 + y*7) % 256)
	b := uint8((y*29 + z*11) % 256)
	col := color.RGBA{r, g, b, 0xff}
	draw.Draw(img, img.Bounds(), &image.Uniform{col}, image.Point{}, draw.Src)

	// Encode PNG to response
	_ = png.Encode(w, img)
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/health", HealthCheck).Methods(http.MethodGet)
	// JSON placeholder (kept for compatibility)
	// Dev tile endpoint for Leaflet and other web clients
	// Register the more specific PNG route before the generic JSON route
	r.HandleFunc("/tiles/{z}/{x}/{y}.png", GetTilePNG).Methods(http.MethodGet)
	// JSON placeholder (kept for compatibility)
	r.HandleFunc("/tiles/{z}/{x}/{y}", GetTile).Methods(http.MethodGet)
	return r
}
