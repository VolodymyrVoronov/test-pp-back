package models

import "github.com/golang-jwt/jwt"

var JWTKey = []byte("your_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var Users = map[string]string{} // A simple in-memory user store
