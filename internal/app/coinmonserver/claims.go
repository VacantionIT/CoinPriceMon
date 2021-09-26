package coinmonserver

import "github.com/dgrijalva/jwt-go"

// Claims - Структура данных токена
type Claims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}
