// TODO:
//      - func that verifies token

package handlers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// reads privRSAKey from env var (which should be in base64), it then decodes it
// after which, it is passed to pem.Decode, who'll try to find valid PEM data.
// if all is well, we'll parse the private key, expecting it to be in the
// PKCS #8, ASN.1 DER format. We then return a pointer to the private key.

func readPrivRSAKeyFromEnv(env string) (*rsa.PrivateKey, error) {
	// base64
	var privRSAKey string = os.Getenv(env)

	decoded, err := base64.StdEncoding.DecodeString(string(privRSAKey))
	if err != nil {
		log.Println("Error decoding private key:", err)
		return nil, err
	}

	block, _ := pem.Decode([]byte(decoded))
	// block will be nil if no pem data is found
	if block == nil {
		log.Println("Error decoding privRSAKey")
		return nil, errors.New("Invalid private RSA key")
	}
	log.Println(block.Type)

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Printf("Error parsing PKCS8 privRSAKey %v\n", err)
		return nil, err
	}

	return key.(*rsa.PrivateKey), nil
}

func GenerateToken(w http.ResponseWriter, r *http.Request) {
	// ###########################################################################
	// ###########################################################################
	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"email": "example@example.com",
	})
	// ###########################################################################
	// ###########################################################################

	priKey, err := readPrivRSAKeyFromEnv("ID_RSA")
	if err != nil {
		log.Println("Error reading private key ->", err)
		return
	}

	token, err := claims.SignedString(priKey)
	if err != nil {
		log.Println("Error al generar el token ->", err)
		return
	}

	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		log.Fatalf("Error writing token to w -> %v\n", err)
	}
}

func ValidateToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	token = token[7:]

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ID_RSA_PUB")), nil
	})
	if err != nil {
		log.Println("Error al parsear el token:", err)
		return
	}

	if parsedToken.Valid {
		log.Println("Token válido")
	} else {
		log.Println("Token inválido")
	}
}
