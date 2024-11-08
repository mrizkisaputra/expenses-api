package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/mrizkisaputra/expenses-api/config"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"github.com/pkg/errors"
	"time"
)

// JWT Claims
type Claim struct {
	ID    string
	Email string
	jwt.RegisteredClaims
}

func GenerateJwtToken(user *model.User, cfg *config.Config, expire time.Duration) (string, error) {
	claims := Claim{
		ID:    user.Id.String(),
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			Issuer:    "jwt",
		},
	}

	// mendeklarasikan token dengan algoritma yang digunakan untuk penandatangan, dan tambahan claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// register the JWT string
	tokenString, err := token.SignedString([]byte(cfg.Server.JWTSecretKey))
	if err != nil {
		return "", errors.Wrap(err, "GenerateJWTTokenPair.SignedString")
	}
	return tokenString, nil
}

// generate JWT access token & refresh token
func GenerateTokenPair(user *model.User, cfg *config.Config) (accToken, refToken string, err error) {
	accToken, err = GenerateJwtToken(user, cfg, 15*time.Minute) // 15 minute
	if err != nil {
		return
	}
	refToken, err = GenerateJwtToken(user, cfg, 1*24*time.Hour) // 1 day
	return
}
