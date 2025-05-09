package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateToken(email, secret, ttl string) (string, error) {
	d, err := time.ParseDuration(ttl)
	if err != nil {
		return "", err
	}

	claims := jwt.RegisteredClaims{
		Subject:   email,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(d)),
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tkn.SignedString([]byte(secret))
}

func ParseToken(token, secret string) (string, error) {
	parsed, err := jwt.ParseWithClaims(
		token,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (any, error) { return []byte(secret), nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return "", err
	}

	claims, ok := parsed.Claims.(*jwt.RegisteredClaims)
	if !ok || !parsed.Valid {
		return "", errors.New("invalid token")
	}
	return claims.Subject, nil
}
