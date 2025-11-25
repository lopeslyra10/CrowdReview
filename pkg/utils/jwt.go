package utils

import (
	"time"

	"crowdreview/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenClaims extends registered claims with a role.
type TokenClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateTokens issues access and refresh JWTs.
func GenerateTokens(userID uuid.UUID, role string, cfg config.Config) (string, string, error) {
	now := time.Now()
	claims := TokenClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refreshClaims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(now.Add(cfg.RefreshTTL)),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	access, err := accessToken.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", "", err
	}
	refresh, err := refreshToken.SignedString([]byte(cfg.RefreshSecret))
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

// ParseToken validates a JWT with the provided secret.
func ParseToken(tokenString, secret string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
