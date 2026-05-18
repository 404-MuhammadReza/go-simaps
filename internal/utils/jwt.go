package utils

import (
	"time"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshClaims struct {
	UserID string
	jwt.RegisteredClaims
}

type AccessClaims struct {
	UserID     string
	UserEmail  string
	UserRole   string
	jwt.RegisteredClaims
}

func GenerateRefreshToken(userID, issuer string, secret []byte, expiry int) (string, error) {
	claims := RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiry) * 24 * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: issuer,
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil { return "", err }
	return token, nil
}

func GenerateAccessToken(userID, userEmail, userRole, issuer string, secret []byte, expiry int) (string, error) {
	claims := AccessClaims{
		UserID: userID,
		UserEmail: userEmail,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiry) * time.Minute)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: issuer,
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil { return "", err }
	return token, nil
}

func ValidateRefreshToken(token string, secret []byte) (*RefreshClaims, error) {
	tokenClaims := &RefreshClaims{}
	secretValidation := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return secret, nil
	}

	parsedToken, err := jwt.ParseWithClaims(token, tokenClaims, secretValidation)
	if err != nil { return nil, err }

	claims, ok := parsedToken.Claims.(*RefreshClaims)
	if ok && parsedToken.Valid { return claims, nil }

	return nil, errors.New("invalid refresh token")
}

func ValidateAccessToken(token string, secret []byte) (*AccessClaims, error) {
	tokenClaims := &AccessClaims{}
	secretValidation := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }

		return secret, nil
	}

	parsedToken, err := jwt.ParseWithClaims(token, tokenClaims, secretValidation)
	if err != nil { return nil, err }

	claims, ok := parsedToken.Claims.(*AccessClaims)
	if ok && parsedToken.Valid { return claims, nil }

	return nil, errors.New("invalid access token")
}