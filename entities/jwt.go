package entities

import "github.com/golang-jwt/jwt/v4"

type JWTCustomClaims struct {
	Dperms []string `json:"dperms"`
	jwt.StandardClaims
}
