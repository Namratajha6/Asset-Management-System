package dbHelper

import (
	"asset_management/database"
	"asset_management/models"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
)

func IsEmployeeExists(email string) (bool, error) {
	var id string
	err := database.Asset.Get(&id, `SELECT id FROM employees WHERE email = $1 AND archived_at IS NULL`, email)
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
	var emp models.EmployeeDetailsResponse
	var assets []models.EmployeeAssetResponse

	query1 := `SELECT name, email, phone_no, type, role FROM employees WHERE id = $1`
	if err := database.Asset.Get(&emp, query1, employeeID); err != nil {
		return emp, nil, err
	}

	query2 := `
		SELECT a.id AS asset_id, a.model, a.asset_type, a.serial_no
		FROM assets a
		INNER JOIN asset_employee_history h ON h.asset_id = a.id
		WHERE h.employee_id = $1
		  AND h.status = 'assigned'
		  AND h.return_date IS NULL
	`
	err := database.Asset.Select(&assets, query2, employeeID)
	return emp, assets, err
}

func ListEmployees(emp models.ListEmployees, page int, limit int) ([]models.EmployeeListResponse, error) {
	query := `
		SELECT id, name, email, phone_no, type, role
		FROM employees
		WHERE archived_at IS NULL
	`
	idx := 1
	var args []any

	if emp.SearchText != "" {
		pattern := "%" + emp.SearchText + "%"
		query += fmt.Sprintf(` AND (
             name ILIKE $%[1]d OR 
             phone_no ILIKE $%[1]d OR 
             email ILIKE $%[1]d)`, idx)
		args = append(args, pattern)
		idx++
	}

	if len(emp.Types) > 0 {
		query += fmt.Sprintf(" AND type::text = ANY($%d)", idx)
		args = append(args, pq.Array(emp.Types))
		idx++
	}

	if len(emp.Roles) > 0 {
		query += fmt.Sprintf(" AND role::text = ANY($%d)", idx)
		args = append(args, pq.Array(emp.Roles))
		idx++
	}

	offset := (page - 1) * limit
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", idx, idx+1)
	args = append(args, limit, offset)

	log.Println("SQL:", query)
	log.Println("ARGS:", args)

	var employees []models.EmployeeListResponse
	err := database.Asset.Select(&employees, query, args...)
	return employees, err
}

func ArchiveEmployee(tx *sqlx.Tx, employeeID string) error {
	const q = `UPDATE employees
		SET    archived_at = NOW()
		WHERE  id = $1;
		`
	_, err := tx.Exec(q, employeeID)
	return err
}
