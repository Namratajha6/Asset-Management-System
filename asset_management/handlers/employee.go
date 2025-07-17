package handlers

import (
	"asset_management/database"
	"asset_management/database/dbHelper"
	"asset_management/models"
	"asset_management/utils"
	"fmt"
	"log"
	"net/http"
)

func CreateEmployeeByEmployee(w http.ResponseWriter, r *http.Request) {
	var req models.EmployeeRequest
	if err := utils.JSON.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "missing input", http.StatusBadRequest)
		return
	}

	ok, err := dbHelper.IsEmployeeExists(req.Email)
	if err != nil {
		http.Error(w, "failed to check employee existence", http.StatusInternalServerError)
		return
	}
	if ok {
		http.Error(w, "employee already exists", http.StatusConflict)
		return
	}

	ok = utils.IsValidCompanyEmail(req.Email)
	if !ok {
		http.Error(w, "invalid email address", http.StatusBadRequest)
		return
	}

	name := utils.GetNameFromEmail(req.Email)
	if name == "" {
		http.Error(w, "can't fetch name from email", http.StatusBadRequest)
		return
	}

	employee := models.Employee{
		Name:  name,
		Email: req.Email,
	}
	err = dbHelper.CreateEmployeeByEmployee(employee)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to create employee", http.StatusInternalServerError)
		return
	}

	err = utils.JSON.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
	})
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var req models.Employee
	if err := utils.JSON.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Type == "" || req.Email == "" || req.PhoneNo == "" || req.Role == "" {
		http.Error(w, "missing input", http.StatusBadRequest)
		return
	}

	ok, err := dbHelper.IsEmployeeExists(req.Email)
	if err != nil {
		http.Error(w, "failed to check employee existence", http.StatusInternalServerError)
		return
	}
	if ok {
		http.Error(w, "employee already exists", http.StatusConflict)
		return
	}

	ok = utils.IsValidCompanyEmail(req.Email)
	if !ok {
		http.Error(w, "invalid email address", http.StatusBadRequest)
		return
	}

	claims, ok := utils.GetClaims(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	req.CreatedBy = claims.UserID

	err = dbHelper.CreateEmployee(req)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to create employee", http.StatusInternalServerError)
		return
	}

	err = utils.JSON.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
	})
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := utils.JSON.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "missing credentials", http.StatusBadRequest)
		return
	}

	emp, err := dbHelper.GetEmployeeByEmail(req.Email)
	if err != nil {
		log.Println("error:", err)
		http.Error(w, "invalid email", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(emp.ID, emp.Role)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(emp.ID, emp.Role)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	err = utils.JSON.NewEncoder(w).Encode(map[string]string{
		"token":         token,
		"refresh_token": refreshToken,
	})

	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func ListEmployees(w http.ResponseWriter, r *http.Request) {
	var req models.ListEmployees

	req.Types = utils.ParseCommaSeparatedParam(r, "types")
	req.Roles = utils.ParseCommaSeparatedParam(r, "roles")
	req.SearchText = r.URL.Query().Get("search")

	page, limit := parsePageLimit(r)

	emp, err := dbHelper.ListEmployees(req, page, limit)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "failed to fetch employees", http.StatusInternalServerError)
		return
	}

	_ = utils.JSON.NewEncoder(w).Encode(map[string]any{"employees": emp})
}

func EmployeeDetails(w http.ResponseWriter, r *http.Request) {
	employeeID := r.URL.Query().Get("id")
	if employeeID == "" {
		http.Error(w, "employeeId is required", http.StatusBadRequest)
		return
	}

	info, assets, err := dbHelper.EmployeeDetails(employeeID)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to fetch employee details", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = utils.JSON.NewEncoder(w).Encode(map[string]interface{}{
		"employee": info,
		"assets":   assets,
	})
}

func EmployeeTimeline(w http.ResponseWriter, r *http.Request) {
	empID := r.URL.Query().Get("id")
	if empID == "" {
		http.Error(w, "employeeId is required", http.StatusBadRequest)
		return
	}
	timeline, err := dbHelper.EmployeeTimeline(empID)
	if err != nil {
		http.Error(w, "failed to fetch employee timeline", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = utils.JSON.NewEncoder(w).Encode(timeline)
}

func ArchiveEmployee(w http.ResponseWriter, r *http.Request) {
	empID := r.URL.Query().Get("id")

	tx, err := database.Asset.Beginx()
	if err != nil {
		http.Error(w, "could not begin transaction", http.StatusInternalServerError)
		return
	}

	status, err := dbHelper.GetStatusByEmpID(tx, empID)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
		}
		log.Println(err)
		http.Error(w, "failed to fetch employee status", http.StatusInternalServerError)
		return
	}
	if status != "" {
		if status == "assigned" {
			http.Error(w, "employee is assigned, first retrieve the asset", http.StatusBadRequest)
			return
		}
	}

	err = dbHelper.ArchiveEmployee(tx, empID)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			http.Error(w, "could not rollback transaction", http.StatusInternalServerError)
			return
		}
		log.Println(err)
		http.Error(w, "failed to fetch employee archive", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = utils.JSON.NewEncoder(w).Encode(map[string]interface{}{
		"message": "archived successfully",
	})
}
