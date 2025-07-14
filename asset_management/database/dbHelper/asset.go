package dbHelper

import (
	"asset_management/database"
	"asset_management/models"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func IsAssetExists(serialNo string) (bool, error) {
	var id string
	err := database.Asset.Get(&id, `SELECT id FROM assets WHERE serial_no = $1`, serialNo)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CreateAsset(tx *sqlx.Tx, asset models.AssetRequest) (string, error) {
	var id string
	err := tx.QueryRowx(`
		INSERT INTO assets (brand, model, serial_no, asset_type, asset_status, owned_by, purchased_date, warranty_start_date, warranty_end_date, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`,
		asset.Brand, asset.Model, asset.SerialNo, asset.AssetType, asset.AssetStatus, asset.OwnedBy, asset.PurchasedDate, asset.WarrantyStartDate, asset.WarrantyEndDate, asset.CreatedBy,
	).Scan(&id)

	return id, err
}

func InsertLaptop(tx *sqlx.Tx, laptop models.Laptop) error {
	res, err := tx.NamedExec(`
		INSERT INTO laptops (asset_id, os, ram, storage, processor, created_by, updated_by)
		VALUES (:asset_id, :os, :ram, :storage, :processor, :created_by, :updated_by)
	`, &laptop)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	fmt.Println("Rows inserted into laptops table:", count)
	return nil
}

func InsertMouse(tx *sqlx.Tx, mouse models.Mouse) error {
	res, err := tx.NamedExec(`
		INSERT INTO mouse (asset_id, connectivity_type, created_by, updated_by)
		VALUES (:asset_id, :connectivity_type, :created_by, :updated_by)
	`, &mouse)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	fmt.Println("Rows inserted into mouse table:", count)
	return nil
}

func InsertHardDisk(tx *sqlx.Tx, hd models.HardDisk) error {
	res, err := tx.NamedExec(`
		INSERT INTO hard_disks (asset_id, storage_capacity, created_by, updated_by)
		VALUES (:asset_id, :storage_capacity, :created_by, :updated_by)
	`, &hd)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	fmt.Println("Rows inserted into hard_disks table:", count)
	return nil
}

func InsertPendrive(tx *sqlx.Tx, pd models.Pendrive) error {
	res, err := tx.NamedExec(`
		INSERT INTO pendrives (asset_id, storage_capacity, created_by, updated_by)
		VALUES (:asset_id, :storage_capacity, :created_by, :updated_by)
	`, &pd)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	fmt.Println("Rows inserted into pendrives table:", count)
	return nil
}

func InsertMobile(tx *sqlx.Tx, mobile models.Mobile) error {
	res, err := tx.NamedExec(`
		INSERT INTO mobiles (asset_id, imei1, imei2, os, ram, storage, created_by, updated_by)
		VALUES (:asset_id, :imei1, :imei2, :os, :ram, :storage, :created_by, :updated_by)
	`, &mobile)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	fmt.Println("Rows inserted into mobiles table:", count)
	return nil
}

func InsertSIM(tx *sqlx.Tx, sim models.SIM) error {
	res, err := tx.NamedExec(`
		INSERT INTO sims (asset_id, mobile_number, network_provider, created_by, updated_by)
		VALUES (:asset_id, :mobile_number, :network_provider, :created_by, :updated_by)
	`, &sim)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	fmt.Println("Rows inserted into sims table:", count)
	return nil
}

func CreateAssetHistory(tx *sqlx.Tx, assetID, newStatus, performedBy string) error {
	_, err := tx.Exec(`
		INSERT INTO asset_history (asset_id, old_status, new_status, performed_by)
		VALUES ($1, NULL, $2, $3)
	`, assetID, newStatus, performedBy)
	return err
}

func InsertAssetHistory(tx *sqlx.Tx, oldStatus string, req models.ChangeAssetStatusRequest) error {
	var employeeID interface{} = nil
	if req.EmployeeID != "" {
		employeeID = req.EmployeeID
	}

	_, err := tx.Exec(`
	INSERT INTO asset_history (employee_id, asset_id, old_status, new_status, performed_by)
	VALUES ($1, $2, $3, $4, $5)
	`, employeeID, req.AssetID, oldStatus, req.Status, req.PerformedBy)
	return err

}

func GetStatus(tx *sqlx.Tx, assetID string) (string, error) {
	var status string
	err := tx.Get(&status, `SELECT asset_status FROM assets WHERE id = $1`, assetID)
	if err != nil {
		return "", err
	}
	return status, nil
}

func UpdateAssetStatus(tx *sqlx.Tx, assetID, newStatus string) error {
	_, err := tx.Exec(`
		UPDATE assets
		SET asset_status = $1
		WHERE id = $2
	`, newStatus, assetID)
	return err
}

func AssignAsset(tx *sqlx.Tx, assign models.ChangeAssetStatusRequest) error {
	_, err := tx.NamedExec(`
			INSERT INTO asset_employee_history(employee_id, asset_id, status, performed_by, assigned_date)
			VALUES (:employee_id, :asset_id, :status, :performed_by, now())`, &assign)
	return err
}

func RetrieveAsset(tx *sqlx.Tx, assign models.ChangeAssetStatusRequest) error {
	res, err := tx.Exec(`
		UPDATE asset_employee_history
		SET    status      = 'retrieved',
		       return_date = NOW(),
		       updated_by  = $1
		WHERE  asset_id = $2
		  AND  status   = 'assigned'
		  AND  return_date IS NULL
	`, assign.PerformedBy, assign.AssetID)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no open assignment to retrieve")
	}
	return nil

}

func ListAllAssets(page, limit int) ([]models.AssetResponse, error) {
	const query = `
		SELECT
			a.brand,
			a.model,
			a.asset_type,
			a.asset_status,
			a.serial_no,
			a.owned_by,
			a.purchased_date,
			h.employee_id                 
		FROM assets a
		LEFT JOIN asset_employee_history h             
		       ON h.asset_id   = a.id
		      AND h.status     = 'assigned'
		      AND h.return_date IS NULL
		WHERE a.archived_at IS NULL
		ORDER BY 
			  CASE WHEN a.asset_status = 'assigned' THEN 0 ELSE 1 END,
			  a.created_at DESC                                         
		LIMIT  $1
		OFFSET $2;
	`

	assets := make([]models.AssetResponse, 0)
	err := database.Asset.Select(&assets, query, limit, (page-1)*limit)
	return assets, err
}

func AssetTimeline(assetID string) ([]models.AssetTimelineResponse, error) {
	const query = `
		SELECT 
			e.name AS employee_name,
			h.status ,
			h.assigned_date,
			h.return_date,
			h.performed_at,
			h.performed_by
		FROM asset_employee_history h
		LEFT JOIN employees e ON h.employee_id = e.id
		WHERE h.asset_id = $1
		ORDER BY h.performed_at DESC;
	`

	var timeline []models.AssetTimelineResponse
	err := database.Asset.Select(&timeline, query, assetID)
	return timeline, err
}

func AssetDetails(assetID string) (models.AssetDetailsResponse, error) {
	const q = `
		SELECT
			a.id,
			a.brand,
			a.model,
			a.serial_no,
			a.asset_type,
			a.asset_status,
			a.owned_by,
			a.purchased_date,
			a.warranty_start_date,
			a.warranty_end_date,

			ass.employee_id,
			e.name AS employee_name
		FROM assets a

		LEFT JOIN LATERAL (
			SELECT h.employee_id
			FROM   asset_employee_history h
			WHERE  h.asset_id   = a.id
			  AND  h.status     = 'assigned'
			  AND  h.return_date IS NULL
			ORDER BY h.performed_at DESC     
			LIMIT 1
		) ass ON TRUE

		LEFT JOIN employees e ON e.id = ass.employee_id

		WHERE a.id = $1
		LIMIT 1;
	`

	var out models.AssetDetailsResponse
	err := database.Asset.Get(&out, q, assetID)
	return out, err

}
