package dbHelper

import (
	"asset_management/database"
	"asset_management/models"
	"database/sql"
)

func IsEmployeeExists(email string) (bool, error) {
	var id string
	err := database.Asset.Get(&id, `SELECT id FROM employees WHERE email = $1`, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CreateEmployeeByEmployee(emp models.Employee) error {
	emp.Type = "employee"
	emp.Role = "employee"

	_, err := database.Asset.NamedExec(`
		INSERT INTO employees (name, email, type, role)
		VALUES (:name, :email, :type, :role)
	`, &emp)
	return err
}

func CreateEmployee(emp models.Employee) error {

	_, err := database.Asset.NamedExec(`
		INSERT INTO employees (name, email,phone_no, type, role, created_by)
		VALUES (:name, :email, :phone_no, :type, :role, :created_by)
	`, &emp)
	return err
}

func GetEmployeeByEmail(email string) (models.Employee, error) {
	var emp models.Employee
	err := database.Asset.Get(&emp,
		`SELECT id, name, type, role FROM employees WHERE email = $1`, email)
	return emp, err
}

func EmployeeTimeline(employeeID string) ([]models.EmployeeTimelineResponse, error) {
	const query = `
		SELECT a.model, a.serial_no, a.asset_type, h.status, h.assigned_date, h.return_date
		FROM asset_employee_history h
		JOIN assets a ON a.id = h.asset_id
		WHERE h.employee_id = $1
		ORDER BY h.performed_at DESC
	`

	var timeline []models.EmployeeTimelineResponse
	err := database.Asset.Select(&timeline, query, employeeID)
	return timeline, err
}

func EmployeeDetails(employeeID string) (models.EmployeeDetailsResponse, []models.EmployeeAssetResponse, error) {
	const infoQuery = `SELECT name, email, phone_no, type, role FROM employees WHERE id = $1`

	var info models.EmployeeDetailsResponse
	if err := database.Asset.Get(&info, infoQuery, employeeID); err != nil {
		return info, nil, err
	}

	const assetQuery = `
		SELECT a.id as asset_id, a.model, a.asset_type, a.serial_no
		FROM assets a
		JOIN asset_employee_history h ON a.id = h.asset_id
		WHERE h.employee_id = $1 AND h.status = 'assigned' AND h.return_date IS NULL
	`

	var assets []models.EmployeeAssetResponse
	err := database.Asset.Select(&assets, assetQuery, employeeID)
	return info, assets, err
}

func ListAllEmployees(page int, limit int) ([]models.EmployeeListResponse, error) {
	const query = `
		SELECT id, name, email, phone_no, type, role
		FROM employees
		WHERE archived_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	offset := (page - 1) * limit
	var employees []models.EmployeeListResponse
	err := database.Asset.Select(&employees, query, limit, offset)
	return employees, err
}
