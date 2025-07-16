package handlers

import (
	"asset_management/database"
	"asset_management/database/dbHelper"
	"asset_management/models"
	"asset_management/utils"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func CreateAsset(w http.ResponseWriter, r *http.Request) {
	var req models.AssetRequest
	if err := utils.JSON.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Brand == "" || req.Model == "" || req.SerialNo == "" || req.AssetType == "" || req.AssetStatus == "" {
		http.Error(w, "missing input", http.StatusBadRequest)
		return
	}

	exists, err := dbHelper.IsAssetExists(req.SerialNo)
	if err != nil {
		http.Error(w, "failed to check asset existence", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "asset already exists", http.StatusConflict)
		return
	}

	claims, ok := utils.GetClaims(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	req.CreatedBy = claims.UserID
	req.UpdatedBy = claims.UserID

	tx, err := database.Asset.Beginx()
	if err != nil {
		http.Error(w, "could not begin transaction", http.StatusInternalServerError)
		return
	}

	// Insert into assets table
	assetID, err := dbHelper.CreateAsset(tx, req)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
			return
		}

		http.Error(w, "failed to insert asset", http.StatusInternalServerError)
		return
	}

	// Insert into specific asset type table
	switch strings.ToLower(req.AssetType) {
	case "laptop":
		err = dbHelper.InsertLaptop(tx, models.Laptop{
			AssetID:   assetID,
			OS:        req.OS,
			RAM:       req.RAM,
			Storage:   req.Storage,
			Processor: req.Processor,
			CreatedBy: req.CreatedBy,
			UpdatedBy: req.UpdatedBy,
		})
	case "mouse":
		err = dbHelper.InsertMouse(tx, models.Mouse{
			AssetID:          assetID,
			ConnectivityType: req.ConnectivityType,
			CreatedBy:        req.CreatedBy,
			UpdatedBy:        req.UpdatedBy,
		})
	case "hard_disk":
		err = dbHelper.InsertHardDisk(tx, models.HardDisk{
			AssetID:         assetID,
			StorageCapacity: req.StorageCapacity,
			CreatedBy:       req.CreatedBy,
			UpdatedBy:       req.UpdatedBy,
		})
	case "pendrive":
		err = dbHelper.InsertPendrive(tx, models.Pendrive{
			AssetID:         assetID,
			StorageCapacity: req.StorageCapacity,
			CreatedBy:       req.CreatedBy,
			UpdatedBy:       req.UpdatedBy,
		})
	case "mobile":
		err = dbHelper.InsertMobile(tx, models.Mobile{
			AssetID:   assetID,
			IMEI1:     req.IMEI1,
			IMEI2:     req.IMEI2,
			OS:        req.OS,
			RAM:       req.RAM,
			Storage:   req.Storage,
			CreatedBy: req.CreatedBy,
			UpdatedBy: req.UpdatedBy,
		})
	case "sim":
		err = dbHelper.InsertSIM(tx, models.SIM{
			AssetID:         assetID,
			MobileNumber:    req.MobileNumber,
			NetworkProvider: req.NetworkProvider,
			CreatedBy:       req.CreatedBy,
			UpdatedBy:       req.UpdatedBy,
		})
	default:
		err := tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
		}
		http.Error(w, "unsupported asset type", http.StatusBadRequest)
		return
	}

	if err != nil {
		err := tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
		}
		log.Printf("Error inserting asset type %s: %v\n", req.AssetType, err)
		http.Error(w, "failed to insert asset details", http.StatusInternalServerError)
		return
	}

	err = dbHelper.CreateAssetHistory(tx, assetID, req.AssetStatus, req.CreatedBy)
	if err != nil {
		_ = tx.Rollback()
		log.Println("Error inserting asset history:", err)
		http.Error(w, "failed to insert asset history", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = utils.JSON.NewEncoder(w).Encode(map[string]string{
		"assetID": assetID,
	})
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func AssignAsset(w http.ResponseWriter, r *http.Request) {
	var req models.ChangeAssetStatusRequest

	if err := utils.JSON.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	claims, ok := utils.GetClaims(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tx, err := database.Asset.Beginx()
	if err != nil {
		http.Error(w, "could not begin transaction", http.StatusInternalServerError)
		return
	}

	status, err := dbHelper.GetStatus(tx, req.AssetID)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
		}
		http.Error(w, "error fetching status", http.StatusInternalServerError)
		return
	}

	if status != "available" {
		err = tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
		}
		http.Error(w, "asset is not available for assignment", http.StatusConflict)
		return
	}

	req.PerformedBy = claims.UserID
	err = dbHelper.AssignAsset(tx, req)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to assign asset", http.StatusInternalServerError)
		return
	}

	err = dbHelper.InsertAssetHistory(tx, status, req)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to insert into asset history", http.StatusInternalServerError)
		return
	}

	err = dbHelper.UpdateAssetStatus(tx, req.AssetID, req.Status)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to update asset status", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = utils.JSON.NewEncoder(w).Encode(map[string]string{
		"message": "asset assigned",
	})

}

func RetrieveAsset(w http.ResponseWriter, r *http.Request) {
	var req models.ChangeAssetStatusRequest

	if err := utils.JSON.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	claims, ok := utils.GetClaims(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tx, err := database.Asset.Beginx()
	if err != nil {
		http.Error(w, "could not begin transaction", http.StatusInternalServerError)
		return
	}

	status, err := dbHelper.GetStatus(tx, req.AssetID)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
		}
		http.Error(w, "error fetching status", http.StatusInternalServerError)
		return
	}

	if status != "assigned" {
		err = tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
		}
		http.Error(w, "asset is not assigned to anyone", http.StatusConflict)
		return
	}

	req.PerformedBy = claims.UserID
	err = dbHelper.RetrieveAsset(tx, req)
	log.Println("Error assigning asset:", err)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to assign asset", http.StatusInternalServerError)
		return
	}

	req.Status = "available"

	err = dbHelper.InsertAssetHistory(tx, status, req)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to insert into asset history", http.StatusInternalServerError)
		return
	}

	err = dbHelper.UpdateAssetStatus(tx, req.AssetID, req.Status)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to update asset status", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	err = utils.JSON.NewEncoder(w).Encode(map[string]string{
		"message": "asset retrieved",
	})
}

func ChangeAssetStatus(w http.ResponseWriter, r *http.Request) {
	var req models.ChangeAssetStatusRequest
	err := utils.JSON.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	claims, ok := utils.GetClaims(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tx, err := database.Asset.Beginx()
	if err != nil {
		http.Error(w, "could not begin transaction", http.StatusInternalServerError)
		return
	}

	status, err := dbHelper.GetStatus(tx, req.AssetID)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
		}
		http.Error(w, "error fetching status", http.StatusInternalServerError)
		return
	}

	if status == "assigned" {
		err = tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
			return
		}
		http.Error(w, "asset is already assigned please retrieve it", http.StatusConflict)
		return
	}

	req.PerformedBy = claims.UserID
	req.EmployeeID = ""

	err = dbHelper.InsertAssetHistory(tx, status, req)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
			return
		}
		log.Println("Error inserting asset history:", err)
		http.Error(w, "failed to insert into asset history", http.StatusInternalServerError)
		return
	}

	err = dbHelper.UpdateAssetStatus(tx, req.AssetID, req.Status)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
			return
		}
		http.Error(w, "failed to update asset status", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = utils.JSON.NewEncoder(w).Encode(map[string]string{
		"message": "asset status changed",
	})

}

func ListAllAssets(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePageLimit(r)
	var req models.ListAssets
	err := utils.JSON.NewDecoder(r.Body).Decode(&req)
	if err != nil {

		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	log.Printf("DEBUG  AssetTypes = %#v\n", req.AssetTypes)

	assets, err := dbHelper.ListAssets(req, page, limit)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to list assets", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = utils.JSON.NewEncoder(w).Encode(map[string]interface{}{
		"assets": assets,
	})
}

func AssetDetails(w http.ResponseWriter, r *http.Request) {
	assetID := r.URL.Query().Get("id")
	if assetID == "" {
		http.Error(w, "assetId is required", http.StatusBadRequest)
		return
	}

	asset, err := dbHelper.AssetDetails(assetID)
	if err != nil {
		http.Error(w, "failed to fetch asset", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = utils.JSON.NewEncoder(w).Encode(asset)
}

func AssetTimeline(w http.ResponseWriter, r *http.Request) {
	assetID := r.URL.Query().Get("id")
	if assetID == "" {
		http.Error(w, "assetId is required", http.StatusBadRequest)
		return
	}

	timeline, err := dbHelper.AssetTimeline(assetID)
	if err != nil {
		http.Error(w, "failed to fetch asset timeline", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = utils.JSON.NewEncoder(w).Encode(timeline)
}

func ArchiveAsset(w http.ResponseWriter, r *http.Request) {
	assetID := r.URL.Query().Get("id")

	tx, err := database.Asset.Beginx()
	if err != nil {
		http.Error(w, "could not begin transaction", http.StatusInternalServerError)
		return
	}

	status, err := dbHelper.GetStatus(tx, assetID)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
		}
		log.Println(err)
		http.Error(w, "failed to fetch employee status", http.StatusInternalServerError)
		return
	}
	if status == "assigned" {
		http.Error(w, "employee is assigned, first retrieve the asset", http.StatusBadRequest)
		return
	}

	err = dbHelper.ArchiveAsset(tx, assetID)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
			return
		}
		log.Println(err)
		http.Error(w, "failed to fetch employee archive", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = utils.JSON.NewEncoder(w).Encode(map[string]interface{}{
		"message": "archived successfully",
	})
}

func parsePageLimit(r *http.Request) (int, int) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return page, limit
}
