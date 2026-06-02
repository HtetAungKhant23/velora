package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/HtetAungKhant23/velora/internal/core/ports"
	"github.com/golang-jwt/jwt/v5"
)

const defaultExpiry = 24 * time.Hour

type JWTTokenService struct {
	secret []byte
	expiry time.Duration
}

type jwtClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewJWTTokenService(secret string, expiry time.Duration) *JWTTokenService {
	if expiry == 0 {
		expiry = defaultExpiry
	}
	return &JWTTokenService{
		secret: []byte(secret),
		expiry: expiry,
	}
}

func (s *JWTTokenService) Generate(userID, email string) (string, error) {
	expiresAt := time.Now().UTC().Add(s.expiry)
	claims := jwtClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("jwt sign: %w", err)
	}

	return signed, nil
}

func (s *JWTTokenService) Validate(tokenStr string) (ports.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		return ports.Claims{}, fmt.Errorf("jwt validate: %w", err)
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return ports.Claims{}, errors.New("jwt: invalid token")
	}

	return ports.Claims{
		UserID: claims.UserID,
		Email:  claims.Email,
	}, nil
}
