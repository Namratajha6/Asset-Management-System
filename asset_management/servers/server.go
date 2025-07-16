package server

import (
	"asset_management/handlers"
	"asset_management/middleware"
	"asset_management/models"
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
	protected.HandleFunc("/dashboard/employee", handlers.EmployeeDashboard).Methods("GET")

	roleAssetProtected := protected.NewRoute().Subrouter()
	roleAssetProtected.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleAssetManager))
	roleAssetProtected.HandleFunc("/admin/dashboard/asset", handlers.AssetStats).Methods("GET")
	roleAssetProtected.HandleFunc("/admin/asset", handlers.CreateAsset).Methods("POST")
	roleAssetProtected.HandleFunc("/admin/asset", handlers.AssetDetails).Methods("GET")
	roleAssetProtected.HandleFunc("/admin/assets", handlers.ListAllAssets).Methods("GET")
	roleAssetProtected.HandleFunc("/admin/assign/asset", handlers.AssignAsset).Methods("POST")
	roleAssetProtected.HandleFunc("/admin/retrieve/asset", handlers.RetrieveAsset).Methods("POST")
	roleAssetProtected.HandleFunc("/admin/asset/status", handlers.ChangeAssetStatus).Methods("POST")
	roleAssetProtected.HandleFunc("/admin/asset/timeline", handlers.AssetTimeline).Methods("GET")
	roleAssetProtected.HandleFunc("/admin/asset/archive", handlers.ArchiveAsset).Methods("DELETE")

	roleEmployeeProtected := protected.NewRoute().Subrouter()
	roleEmployeeProtected.Use(middleware.RoleMiddleware(models.RoleAdmin, models.RoleEmployeeManager))
	roleEmployeeProtected.HandleFunc("/admin/employee", handlers.CreateEmployee).Methods("POST")
	roleEmployeeProtected.HandleFunc("/admin/employee", handlers.EmployeeDetails).Methods("GET")
	roleEmployeeProtected.HandleFunc("/admin/employees", handlers.ListEmployees).Methods("GET")
	roleEmployeeProtected.HandleFunc("/admin/employee/timeline", handlers.EmployeeTimeline).Methods("GET")
	roleEmployeeProtected.HandleFunc("/admin/employee/archive", handlers.ArchiveEmployee).Methods("DELETE")

	return r
}
