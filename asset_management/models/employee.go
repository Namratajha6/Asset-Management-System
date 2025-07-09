package models

type EmployeeRequest struct {
	Email string `json:"email"`
}

type LoginRequest struct {
	Email string `json:"email"`
}

type EmployeeReq struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	PhoneNo string `json:"phoneNo"`
	Type    string `json:"type"`
	Role    string `json:"role"`
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
