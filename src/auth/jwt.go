package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/zspekt/capugo/internal/utils"
)

/*
reads privRSAKey from env var (which should be in base64), it then decodes it,
after which it is passed to pem.Decode, who'll try to find valid PEM data.
if all is well, it'll parse the private key, expecting it (the key) to be in the
PKCS #8, ASN.1 DER format. It then return a pointer to the private key.

ref https://stackoverflow.com/questions/44230634/how-to-read-an-rsa-key-from-file
*/
func ReadPrivRSAKeyFromEnv(env string) (*rsa.PrivateKey, error) {
	slog.Debug("running readPrivRSAKeyFromEnv")

	// base64
	var privRSAKey string = os.Getenv(env)
	if len(privRSAKey) == 0 {
		utils.SlogFatal("Env var is not set", "env", env)
	}

	decoded, err := base64.StdEncoding.DecodeString(string(privRSAKey))
	if err != nil {
		slog.Error("Error decoding privRSAKey from base64", "error", err)
		return nil, err
	}

	block, _ := pem.Decode([]byte(decoded))
	// block will be nil if no pem data is found
	if block == nil {
		err = errors.New("Invalid private RSA key")
		slog.Error("Error decoding privRSAKey pem", "error", err)
		return nil, err
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		slog.Error("Error parsing PKCS8 privRSAKey", "error", err)
		return nil, err
	}

	return key.(*rsa.PrivateKey), nil
}

// gets Token from header
func GetTokenFromHeader(r *http.Request) (string, error) {
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

func ReadPubRSAKeyFromEnv(env string) (*rsa.PublicKey, error) {
	slog.Debug("running readPubRSAKeyFromEnv")

	var publicRSAKey string = os.Getenv(env)
	if len(publicRSAKey) == 0 {
		utils.SlogFatal("env var is not set", "env", env)
	}

	block, _ := pem.Decode([]byte(publicRSAKey))
	// block will be nil if no pem data is found
	if block == nil {
		err := errors.New("Invalid Public RSA key")
		slog.Error("error decoding publicRSAKey", "error", err)
		return nil, err
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		slog.Error("error parsing PKIX publicRSAKey", "error", err)
		return nil, err
	}

	return key.(*rsa.PublicKey), nil
}
