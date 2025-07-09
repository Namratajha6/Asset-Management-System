package handlers

import (
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
	var req models.EmployeeReq
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

	employee := models.Employee{
		Name:      req.Name,
		Email:     req.Email,
		PhoneNo:   req.PhoneNo,
		Role:      req.Role,
		Type:      req.Type,
		CreatedBy: claims.UserID,
	}

	err = dbHelper.CreateEmployee(employee)
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
