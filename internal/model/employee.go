package model

import (
	"time"
)

type Employee struct {
	ID        int       `json:"id" db:"id"`
	FullName  string    `json:"full_name" db:"full_name" binding:"required"`
	Phone     string    `json:"phone" db:"phone" binding:"required"`
	City      string    `json:"city" db:"city" binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateEmployeeRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=100"`
	Phone    string `json:"phone" binding:"required,min=10,max=20"`
	City     string `json:"city" binding:"required,min=2,max=50"`
}

type EmployeeResponse struct {
	ID        int       `json:"id"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
