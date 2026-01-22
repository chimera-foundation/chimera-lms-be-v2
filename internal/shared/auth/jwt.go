package auth

import (
	"time"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtProvider struct {
	secretKey []byte
	expiryDuration time.Duration
}

func NewJWTProvider(secret string, expiry time.Duration) *jwtProvider {
	return &jwtProvider{
		secretKey: []byte(secret),
		expiryDuration: expiry,
	}
}

type CustomClaims struct {
	UserID         uuid.UUID `json:"user_id"`
	OrganizationID uuid.UUID `json:"org_id"`
	jwt.RegisteredClaims
}

func (j *jwtProvider) GenerateToken(userID uuid.UUID, orgID uuid.UUID) (string, error) {
	claims := CustomClaims {
		UserID: userID,
		OrganizationID: orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiryDuration)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

func (j *jwtProvider) ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return j.secretKey, nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid claims")
	}

	return claims.UserID, nil
}