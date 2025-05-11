package server

import (
	"encoding/json"
	"homelab-inventory/internal/logging"
	"net/http"
	"sync"

	"homelab-inventory/pkg/model"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	sysInfos []model.SystemInfo
	mu       sync.Mutex
)

func StartServer(port string) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		if err != nil {
			logging.Logger.Errorw("Error encoding JSON", "error", err)
			return
		}
	})
	r.Post("/sysinfo", handlePost)

	logging.Logger.Infow("Starting server on port", "port", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		return
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var info model.SystemInfo
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		logging.Logger.Errorw("Error decoding JSON", "error", err)
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	sysInfos = append(sysInfos, info)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"status": "received"})
	if err != nil {
		return
	}
}
