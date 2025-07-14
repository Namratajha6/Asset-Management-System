package dbHelper

import (
	"asset_management/database"
	"asset_management/models"
	"github.com/jmoiron/sqlx"
)

func GetAssetStats() (models.AssetStatsResponse, error) {
	const q = `
        SELECT
            COUNT(*)                                                       AS total,
            COUNT(*) FILTER (WHERE asset_status = 'available')            AS available,
            COUNT(*) FILTER (WHERE asset_status = 'assigned')             AS assigned,
            COUNT(*) FILTER (WHERE asset_status = 'waiting_for_repair')   AS waiting_for_repair,
            COUNT(*) FILTER (WHERE asset_status = 'service')              AS service,
            COUNT(*) FILTER (WHERE asset_status = 'damaged')              AS damaged
        FROM assets
        WHERE archived_at IS NULL;
    `
	var stats models.AssetStatsResponse
	err := database.Asset.Get(&stats, q)
	return stats, err
}

func GetEmployeeNameTx(tx *sqlx.Tx, id string) (string, error) {
	var name string
	err := tx.Get(&name, `SELECT name FROM employees WHERE id = $1`, id)
	return name, err
}

func GetActiveAssetCountTx(tx *sqlx.Tx, employeeID string) (int, error) {
	var count int
	err := tx.Get(&count, `
		SELECT COUNT(*) FROM asset_employee_history
		WHERE employee_id = $1
		  AND status      = 'assigned'
		  AND return_date IS NULL`, employeeID)
	return count, err
}

func GetAssignedAssetsTx(tx *sqlx.Tx, employeeID string) ([]models.AssignedAsset, error) {
	var list []models.AssignedAsset
	err := tx.Select(&list, `
		SELECT a.brand,
		       a.model,
		       a.serial_no,
		       TO_CHAR(h.assigned_date, 'DD/MM/YYYY') AS assigned_date,
		       h.status
		FROM   asset_employee_history h
		JOIN   assets a ON a.id = h.asset_id
		WHERE  h.employee_id = $1
		  AND  h.status      = 'assigned'
		  AND  h.return_date IS NULL
		ORDER  BY h.assigned_date DESC`, employeeID)
	return list, err
}
