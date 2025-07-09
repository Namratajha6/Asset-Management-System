package server

import (
	"asset_management/handlers"
	"asset_management/middleware"
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

	r.HandleFunc("/api/v1/public/employee", handlers.CreateEmployeeByEmployee).Methods("POST")
	r.HandleFunc("/api/v1/auth/login", handlers.Login).Methods("POST")

	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	roleProtected := protected.NewRoute().Subrouter()
	roleProtected.Use(middleware.RoleMiddleware("admin", "asset_manager"))
	roleProtected.HandleFunc("/admin/employee", handlers.CreateEmployee).Methods("POST")

	return r
}
