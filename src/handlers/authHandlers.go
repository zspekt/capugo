package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt"

	"github.com/zspekt/capugo/internal/utils"
	"github.com/zspekt/capugo/src/auth"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func GenerateToken(w http.ResponseWriter, r *http.Request) {
	slog.Debug("running GenerateToken")
	// ###########################################################################
	// ###########################################################################
	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"email": "example@example.com",
	})
	// ###########################################################################
	// ###########################################################################

	priKey, err := auth.ReadPrivRSAKeyFromEnv("ID_RSA")
	if err != nil {
		slog.Error("error reading private key", "error", err)
		return
	}

	token, err := claims.SignedString(priKey)
	if err != nil {
		slog.Error("error generating token", "error", err)
		return
	}

	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		utils.SlogFatal("error writing token to w", "error", err)
	}
}

// we validate the token and return an example custom field
func ValidateToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetTokenFromHeader(r)
	if err != nil {
		utils.SlogFatal("error getting token from auth header", "error", err)
	}

	pubRSAKey, err := auth.ReadPubRSAKeyFromEnv("ID_RSA_PUB")
	if err != nil {
		utils.SlogFatal("error reading pubRSAKey", "error", err)
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return pubRSAKey, nil
	})
	if err != nil {
		slog.Error("error parsing token", "error", err)
		return
	}

	// this is one way to get the custom/sub claims from the jwt
	// ref https://stackoverflow.com/questions/61281636/how-to-access-jwt-sub-claims-using-go
	claims := parsedToken.Claims.(jwt.MapClaims)

	data := claims["email"].(string)

	if parsedToken.Valid {
		slog.Info("token is valid", "data", data)
	} else {
		slog.Info("token is NOT valid")
	}
}
