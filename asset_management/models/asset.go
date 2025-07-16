package models

import (
	"time"
)

type AssetRequest struct {
	Brand             string `json:"brand"`
	Model             string `json:"model"`
	SerialNo          string `json:"serialNo"`
	AssetType         string `json:"assetType"`
	AssetStatus       string `json:"assetStatus"`
	OwnedBy           string `json:"ownedBy"`
	PurchasedDate     string `json:"purchasedDate"`
	WarrantyStartDate string `json:"warrantyStartDate"`
	WarrantyEndDate   string `json:"warrantyEndDate"`
	OS                string `json:"os"`
	RAM               string `json:"ram"`
	Storage           string `json:"storage"`
	Processor         string `json:"processor"`
	ConnectivityType  string `json:"connectivityType"`
	StorageCapacity   string `json:"storageCapacity"`
	IMEI1             string `json:"imei1"`
	IMEI2             string `json:"imei2"`
	MobileNumber      string `json:"mobileNumber"`
	NetworkProvider   string `json:"networkProvider"`
	CreatedBy         string `json:"createdBy"`
	UpdatedBy         string `json:"updatedBy"`
}

type AssetsResponse struct {
	Brand         string    `json:"brand" db:"brand"`
	Model         string    `json:"model" db:"model"`
	SerialNo      string    `json:"serialNo" db:"serial_no"`
	AssetType     string    `json:"assetType" db:"asset_type"`
	AssetStatus   string    `json:"assetStatus" db:"asset_status"`
	OwnedBy       string    `json:"ownedBy" db:"owned_by"`
	PurchasedDate time.Time `json:"purchasedDate" db:"purchased_date"`
	AssignedTo    *string   `json:"assignedTo" db:"employee_id"`
}

type AssetTimelineResponse struct {
	EmployeeName *string    `json:"employeeName" db:"employee_name"` // nullable
	Status       *string    `json:"status" db:"status"`              // assigned/retrieved/etc.
	AssignedDate *time.Time `json:"assignedDate" db:"assigned_date"` // nullable
	ReturnDate   *time.Time `json:"returnDate" db:"return_date"`     // nullable
	PerformedAt  time.Time  `json:"performedAt" db:"performed_at"`   // not null
	PerformedBy  string     `json:"performedBy" db:"performed_by"`   // who performed action
}

type AssetDetailsResponse struct {
	ID                string     `json:"id"                db:"id"`
	Brand             string     `json:"brand"             db:"brand"`
	Model             string     `json:"model"             db:"model"`
	SerialNo          string     `json:"serialNo"          db:"serial_no"`
	AssetType         string     `json:"assetType"         db:"asset_type"`
	AssetStatus       string     `json:"assetStatus"       db:"asset_status"`
	OwnedBy           string     `json:"ownedBy"           db:"owned_by"`
	PurchasedDate     *time.Time `json:"purchasedDate"     db:"purchased_date"`
	WarrantyStartDate *time.Time `json:"warrantyStartDate" db:"warranty_start_date"`
	WarrantyEndDate   *time.Time `json:"warrantyEndDate"   db:"warranty_end_date"`
	EmployeeID        *string    `json:"employeeId"   db:"employee_id"`
	EmployeeName      *string    `json:"employeeName" db:"employee_name"`
}

type ChangeAssetStatusRequest struct {
	EmployeeID  string `json:"employeeId" db:"employee_id"`
	AssetID     string `json:"assetId" db:"asset_id"`
	Status      string `json:"status" db:"status"`
	PerformedBy string `json:"performedBy" db:"performed_by"`
}

type Laptop struct {
	ID         string     `json:"id" db:"id"`
	AssetID    string     `json:"assetId" db:"asset_id"`
	OS         string     `json:"os" db:"os"`
	RAM        string     `json:"ram" db:"ram"`
	Storage    string     `json:"storage" db:"storage"`
	Processor  string     `json:"processor" db:"processor"`
	CreatedAt  time.Time  `json:"createdAt" db:"created_at"`
	ArchivedAt *time.Time `json:"archivedAt" db:"archived_at"`
	CreatedBy  string     `json:"createdBy" db:"created_by"`
	UpdatedBy  string     `json:"updatedBy" db:"updated_by"`
}

type Mouse struct {
	ID               string     `json:"id" db:"id"`
	AssetID          string     `json:"assetId" db:"asset_id"`
	ConnectivityType string     `json:"connectivityType" db:"connectivity_type"`
	CreatedAt        time.Time  `json:"createdAt" db:"created_at"`
	ArchivedAt       *time.Time `json:"archivedAt" db:"archived_at"`
	CreatedBy        string     `json:"createdBy" db:"created_by"`
	UpdatedBy        string     `json:"updatedBy" db:"updated_by"`
}

type HardDisk struct {
	ID              string     `json:"id" db:"id"`
	AssetID         string     `json:"assetId" db:"asset_id"`
	StorageCapacity string     `json:"storageCapacity" db:"storage_capacity"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
	ArchivedAt      *time.Time `json:"archivedAt" db:"archived_at"`
	CreatedBy       string     `json:"createdBy" db:"created_by"`
	UpdatedBy       string     `json:"updatedBy" db:"updated_by"`
}

type Pendrive struct {
	ID              string     `json:"id" db:"id"`
	AssetID         string     `json:"assetId" db:"asset_id"`
	StorageCapacity string     `json:"storageCapacity" db:"storage_capacity"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
	ArchivedAt      *time.Time `json:"archivedAt" db:"archived_at"`
	CreatedBy       string     `json:"createdBy" db:"created_by"`
	UpdatedBy       string     `json:"updatedBy" db:"updated_by"`
}

type Mobile struct {
	ID         string     `json:"id" db:"id"`
	AssetID    string     `json:"assetId" db:"asset_id"`
	IMEI1      string     `json:"imei1" db:"imei1"`
	IMEI2      string     `json:"imei2" db:"imei2"`
	OS         string     `json:"os" db:"os"`
	RAM        string     `json:"ram" db:"ram"`
	Storage    string     `json:"storage" db:"storage"`
	CreatedAt  time.Time  `json:"createdAt" db:"created_at"`
	ArchivedAt *time.Time `json:"archivedAt" db:"archived_at"`
	CreatedBy  string     `json:"createdBy" db:"created_by"`
	UpdatedBy  string     `json:"updatedBy" db:"updated_by"`
}

type SIM struct {
	ID              string     `json:"id" db:"id"`
	AssetID         string     `json:"assetId" db:"asset_id"`
	MobileNumber    string     `json:"mobileNumber" db:"mobile_number"`
	NetworkProvider string     `json:"networkProvider" db:"network_provider"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
	ArchivedAt      *time.Time `json:"archivedAt" db:"archived_at"`
	CreatedBy       string     `json:"createdBy" db:"created_by"`
	UpdatedBy       string     `json:"updatedBy" db:"updated_by"`
}

type AssetSearchResponse struct {
	Brand         string `json:"brand"`
	Model         string `json:"model"`
	AssetType     string `json:"assetType" db:"asset_type"`
	SerialNo      string `json:"serialNo" db:"serial_no"`
	AssetStatus   string `json:"assetStatus" db:"asset_status"`
	OwnedBy       string `json:"ownedBt" db:"owned_by"`
	PurchasedDate string `json:"purchasedDate" db:"purchased_date"`
}

type ListAssets struct {
	SearchText  string   `json:"searchText"`
	AssetTypes  []string `json:"assetTypes"`
	AssetStatus []string `json:"assetStatus"`
	OwnedBy     string   `json:"ownedBy"`
}
