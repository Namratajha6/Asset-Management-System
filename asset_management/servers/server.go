package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func SetupRoutes() http.Handler {
	r := mux.NewRouter()

	// Health check route
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]bool{"ok": true}); err != nil {
			logrus.Errorf("Error encoding health response: %v", err)
		}
	}).Methods("GET")

	return r
}
