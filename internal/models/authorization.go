package models

// LoginRequest swagger::model
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
