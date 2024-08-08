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

type key interface {
	Size() int
}

func DecodeAndParse(keyString string, isPrivate bool) (key, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(string(keyString))
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(decodedKey))
	// block will be nil if no pem data is found
	if block == nil {
		return nil, errors.New("pem.Decode couldn't find any pem data")
	}

	switch {
	case isPrivate && block.Type == "RSA PRIVATE KEY":
		k, err := x509.ParsePKCS8PrivateKey(block.Bytes) // good case key is private
		if err != nil {
			return nil, err
		}
		key := k.(*rsa.PrivateKey)
		return key, nil

	case !isPrivate && block.Type == "RSA PUBLIC KEY":
		k, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		key := k.(*rsa.PublicKey)
		return key, nil

	default:
		return nil, fmt.Errorf(
			"key block type <%v> does not match isPrivate argument <%v>",
			block.Type,
			isPrivate,
		)
	}
}

func readFromEnv(envVar string) (string, error) {
	if envVar == "" {
		return "", errors.New("caller passed an empty string as envVar parameter")
	}

	// os.LookupEnv returns false if var is not set
	v, ok := os.LookupEnv(envVar)
	switch {
	case !ok:
		fmt.Printf("\nhere is the shit that should be empty -> %v\n", envVar)
		return "", fmt.Errorf("%v var missing from environment...", envVar)
	case v == "":
		return "", fmt.Errorf("%v var set but empty (how did this happen?)...", envVar)
	}
	return v, nil
}

/*
gets Secret RSA Key from env var (which should be in base64), it then decodes it,
after which it is passed to pem.Decode, who'll try to find valid PEM data.
if all is well, it'll parse the private key, expecting it (the key) to be in the
PKCS #8, ASN.1 DER format. It then return a pointer to the private key.

ref https://stackoverflow.com/questions/44230634/how-to-read-an-rsa-key-from-file
*/
func GetSecKey(env string) (any, error) {
	slog.Debug("running GetSecKey")

	// base64
	privRSAKey, err := readFromEnv(env)
	if err != nil {
		return nil, err
	}

	k, err := DecodeAndParse(privRSAKey, true)
	if err != nil {
		return nil, err
	}

	return k, nil
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
