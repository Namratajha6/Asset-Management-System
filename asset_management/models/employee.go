package models

import "time"

const (
	RoleAdmin           = "admin"
	RoleEmployee        = "employee"
	RoleAssetManager    = "asset_manager"
	RoleEmployeeManager = "employee_manager"
)

type EmployeeRequest struct {
	Email string `json:"email"`
}

type LoginRequest struct {
	Email string `json:"email"`
}

type Employee struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Email     string `json:"email" db:"email"`
	PhoneNo   string `json:"phoneNo" db:"phone_no"`
	Type      string `json:"type" db:"type"`
	Role      string `json:"role" db:"role"`
	CreatedBy string `json:"createdBy" db:"created_by"`
}

type EmployeeListResponse struct {
	ID      string  `json:"id" db:"id"`
	Name    string  `json:"name" db:"name"`
	Email   string  `json:"email" db:"email"`
	PhoneNo *string `json:"phoneNo" db:"phone_no"`
	Type    string  `json:"type" db:"type"`
	Role    string  `json:"role" db:"role"`
}

type EmployeeDetailsResponse struct {
	Name    string  `json:"name" db:"name"`
	Email   string  `json:"email" db:"email"`
	PhoneNo *string `json:"phoneNo" db:"phone_no"`
	Type    string  `json:"type" db:"type"`
	Role    string  `json:"role" db:"role"`
}

type EmployeeAssetResponse struct {
	Model     string `json:"model" db:"model"`
	AssetType string `json:"assetType" db:"asset_type"`
	SerialNo  string `json:"serialNo" db:"serial_no"`
	AssetID   string `json:"assetId" db:"asset_id"`
}

type EmployeeTimelineResponse struct {
	Model      string     `json:"model" db:"model"`
	SerialNo   string     `json:"serialNo" db:"serial_no"`
	AssetType  string     `json:"assetType" db:"asset_type"`
	Status     string     `json:"status" db:"status"`
	AssignedAt *time.Time `json:"assignedAt" db:"assigned_date"`
	ReturnedAt *time.Time `json:"returnedAt" db:"return_date"`
}

type ListEmployees struct {
	SearchText string   `json:"searchText"`
	Roles      []string `json:"roles"`
	Types      []string `json:"types"`
}
