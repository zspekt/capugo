package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GenerateToken(w http.ResponseWriter, r *http.Request) {
	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"email": "example@example.com",
	})

	priKey, err := os.ReadFile(os.Getenv("ID_RSA"))
	if err != nil {
		fmt.Println("Error al leer la clave privada:", err)
		return
	}

	b64, err := base64.StdEncoding.DecodeString(string(priKey))
	if err != nil {
		fmt.Println("Error al decodificar la clave privada:", err)
		return
	}

	token, err := claims.SignedString(b64)
	if err != nil {
		fmt.Println("Error al generar el token:", err)
		return
	}

	json.NewEncoder(w).Encode(token)
}

func ValidateToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	token = token[7:]

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ID_RSA_PUB")), nil
	})
	if err != nil {
		fmt.Println("Error al parsear el token:", err)
		return
	}

	if parsedToken.Valid {
		fmt.Println("Token válido")
	} else {
		fmt.Println("Token inválido")
	}

}