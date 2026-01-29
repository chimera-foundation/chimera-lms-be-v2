package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
)

type jwtProvider struct {
	secretKey []byte
	expiryDuration time.Duration
	redis *redis.Client
}

func NewJWTProvider(secret string, expiry time.Duration, redisClient *redis.Client) TokenProvider {
	return &jwtProvider{
		secretKey: []byte(secret),
		expiryDuration: expiry,
		redis: redisClient,
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

func (j *jwtProvider) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
        return j.secretKey, nil
	})

	if err != nil {
        return nil, fmt.Errorf("token parsing failed: %w", err)
    }

    if !token.Valid {
        return nil, errors.New("invalid token")
    }

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

func (j *jwtProvider) BlacklistToken(ctx context.Context, token string, expiration time.Duration) error {
    key := "blacklist:" + token
    return j.redis.Set(ctx, key, "true", expiration).Err()
}

func (j *jwtProvider) IsBlacklisted(ctx context.Context, token string) (bool, error) {
    key := "blacklist:" + token
    exists, err := j.redis.Exists(ctx, key).Result()
    return exists > 0, err
}