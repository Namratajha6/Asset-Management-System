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
