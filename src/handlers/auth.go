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
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

/*
reads privRSAKey from env var (which should be in base64), it then decodes it
after which, it is passed to pem.Decode, who'll try to find valid PEM data.
if all is well, we'll parse the private key, expecting it to be in the
PKCS #8, ASN.1 DER format. We then return a pointer to the private key.

ref https://stackoverflow.com/questions/44230634/how-to-read-an-rsa-key-from-file
*/
func readPrivRSAKeyFromEnv(env string) (*rsa.PrivateKey, error) {
	// base64
	var privRSAKey string = os.Getenv(env)
	if len(privRSAKey) == 0 {
		log.Fatalf("%v env var is not set...", env)
	}

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

func getTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header is missing\n")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("Invalid Auth header format.\n")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	return token, nil
}

func readPubRSAKeyFromEnv(env string) (*rsa.PublicKey, error) {
	var publicRSAKey string = os.Getenv(env)
	if len(publicRSAKey) == 0 {
		log.Fatalf("%v env var is not set...", env)
	}

	block, _ := pem.Decode([]byte(publicRSAKey))
	// block will be nil if no pem data is found
	if block == nil {
		log.Println("Error decoding publicRSAKey")
		return nil, errors.New("Invalid Public RSA key")
	}
	log.Println(block.Type)

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Printf("Error parsing PKIX publicRSAKey %v\n", err)
		return nil, err
	}

	return key.(*rsa.PublicKey), nil
}

// we validate the token and return an example custom field
func ValidateToken(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenFromHeader(r)
	if err != nil {
		log.Fatalf("Error getting token from auth header -> %v", err)
	}

	pubRSAKey, err := readPubRSAKeyFromEnv("ID_RSA_PUB")
	if err != nil {
		log.Fatalf("Error reading pubRSAKey -> %v\n", err)
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return pubRSAKey, nil
	})
	if err != nil {
		log.Println("Error parsing token ->", err)
		return
	}

	// this is one way to get the custom/sub claims from the jwt
	// ref https://stackoverflow.com/questions/61281636/how-to-access-jwt-sub-claims-using-go
	claims := parsedToken.Claims.(jwt.MapClaims)

	data := claims["email"].(string)

	if parsedToken.Valid {
		log.Println("Token is valid")
		log.Println(
			data,
		)
	} else {
		log.Println("Token is NOT valid")
	}
}
