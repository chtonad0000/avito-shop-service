package auth

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	EmployeeUsername string `json:"employee_username"`
	jwt.RegisteredClaims
}
