package models

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

type AssignRequest struct {
	EmployeeID  string `json:"employeeId" db:"employee_id"`
	AssetID     string `json:"assetId" db:"asset_id"`
	Status      string `json:"status" db:"status"`
	PerformedBy string `json:"performedBy" db:"performed_by"`
}
