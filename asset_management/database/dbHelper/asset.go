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

func InsertAssetHistory(tx *sqlx.Tx, assetID, newStatus, performedBy string) error {
	_, err := tx.Exec(`
		INSERT INTO asset_history (asset_id, old_status, new_status, performed_by)
		VALUES ($1, NULL, $2, $3)
	`, assetID, newStatus, performedBy)
	return err
}
