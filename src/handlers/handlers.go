package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type authParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResp struct {
	Token string `json:"token"`
}

// login interface
type LoginRequest struct {
	Username string `json:"username" validate:"required" example:"email@email.com"`
	Password string `json:"password" validate:"required" example:"secure_password"`
}

// LoginHandler
// @Summary Login Method
// @Description Login by username and password
// @Tags Sign In
// @ID auth-login
// @Accept  json
// @Produce  json
// @Param loginRequest body LoginRequest true "Login Request"
// @Router /auth/login [post]
// @Success 200 {object} authResp
// @Failure 400 {object} string
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// do something with loginRequest
}

// HealthCheck godoc
// @Summary Health Check
// @Description Health Check
// @Tags health
// @ID health-check
// @Accept  json
// @Produce  json
// @Router /api/v1/health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	log.Println("healthCheck handler called...")

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
