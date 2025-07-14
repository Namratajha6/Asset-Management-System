package models

type AssetStatsResponse struct {
	Total            int `json:"total"`
	Available        int `json:"available"`
	Assigned         int `json:"assigned"`
	WaitingForRepair int `json:"waitingForRepair" db:"waiting_for_repair"`
	Service          int `json:"service"`
	Damaged          int `json:"damaged"`
}

type AssignedAsset struct {
	Brand        string `db:"brand" json:"brand"`
	Model        string `db:"model" json:"model"`
	SerialNumber string `db:"serial_no" json:"serialNumber"`
	AssignedDate string `db:"assigned_date" json:"assignedDate"`
	Status       string `db:"status" json:"status"`
}

type DashboardData struct {
	EmployeeName   string          `json:"employeeName"`
	ActiveAssets   int             `json:"activeAssets"`
	AssignedAssets []AssignedAsset `json:"assignedAssets"`
}
