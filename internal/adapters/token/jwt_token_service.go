package token

import (
	"fmt"
	"time"

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

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("jwt sign: %w", err)
	}

	return signed, nil
}
