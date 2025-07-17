package handlers

import (
	"asset_management/database"
	"asset_management/database/dbHelper"
	"asset_management/models"
	"asset_management/utils"
	"log"
	"net/http"
)

func AssetStats(w http.ResponseWriter, _ *http.Request) {
	stats, err := dbHelper.GetAssetStats()
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to fetch dashboard stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = utils.JSON.NewEncoder(w).Encode(stats)
}

func EmployeeDashboard(w http.ResponseWriter, r *http.Request) {
	claims, ok := utils.GetClaims(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	empID := claims.UserID

	tx, err := database.Asset.Beginx()
	if err != nil {
		http.Error(w, "failed to begin txn", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var data models.DashboardData

	if data.EmployeeName, err = dbHelper.GetEmployeeNameTx(tx, empID); err != nil {
		http.Error(w, "failed to fetch name", http.StatusInternalServerError)
		return
	}

	if data.ActiveAssets, err = dbHelper.GetActiveAssetCountTx(tx, empID); err != nil {
		http.Error(w, "failed to count assets", http.StatusInternalServerError)
		return
	}

	if data.AssignedAssets, err = dbHelper.GetAssignedAssetsTx(tx, empID); err != nil {
		http.Error(w, "failed to list assets", http.StatusInternalServerError)
		return
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, "commit failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = utils.JSON.NewEncoder(w).Encode(data)
}
