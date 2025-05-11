package server

import (
	"encoding/json"
	"homelab-inventory/internal/logging"
	"homelab-inventory/internal/storage"
	"homelab-inventory/internal/version"
	"net/http"

	"homelab-inventory/pkg/model"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func StartServer(port string) {

	store, err := storage.NewSQLiteStorage("./homelab.db")
	if err != nil {
		logging.Logger.Fatalw("failed to start sqlite", "error", err)
	}

	r := chi.NewRouter()
	r.Use(ZapLogger(logging.Logger))
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		v := version.Get()
		err := json.NewEncoder(w).Encode(model.HealthResponse{
			Status:    "ok",
			Version:   v.Version,
			Commit:    v.Commit,
			BuildTime: v.BuildTime,
			GoVersion: v.GoVersion,
		})
		if err != nil {
			logging.Logger.Errorw("Error encoding JSON", "error", err)
			return
		}
	})
	r.Post("/sysinfo", handlePost(store))

	logging.Logger.Infow("Starting server on port", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logging.Logger.Fatalw("Error starting server", "error", err)
		return
	}
}

func handlePost(store *storage.SQLiteStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var info model.SystemInfo
		if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := store.SaveSystemInfo(&info); err != nil {
			logging.Logger.Errorw("failed to save system info", "error", err)
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]string{"status": "received"})
		if err != nil {
			logging.Logger.Errorw("Error encoding JSON", "error", err)
			return
		}
	}
}
