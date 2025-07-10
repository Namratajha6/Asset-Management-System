package handlers

import (
	"asset_management/database"
	"asset_management/database/dbHelper"
	"asset_management/models"
	"asset_management/utils"
	"log"
	"net/http"
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

	err = dbHelper.InsertAssetHistory(tx, assetID, req.AssetStatus, req.CreatedBy)
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

	err = utils.JSON.NewEncoder(w).Encode(map[string]string{
		"assetID": assetID,
	})
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
