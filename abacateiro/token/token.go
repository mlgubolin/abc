package token

import (
	"application"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("r5g1er65g1er65g1e6r5g1er65g1e6r51ge6r51ge65rg1er651g65er1") // Secret key for signing the token

type TokenService struct{}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (s *TokenService) GenerateToken(authInfo *application.AuthInfo) (*application.Token, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"userId":     authInfo.AuthId,
		"email":      authInfo.Email,
		"exp":        time.Now().Add(15 * time.Minute).Unix(), // Expires in 15 minutes
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	return &application.Token{Token: tokenString, Expiry: time.Now().Add(15 * time.Minute)}, nil
}

func (s *TokenService) ValidateToken(tokenString string) error {
	if len(strings.Split(tokenString, ".")) != 3 {
		return fmt.Errorf("invalid token format")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return fmt.Errorf("invalid token: %v", err)
	}

	if expClaim, ok := token.Claims.(jwt.MapClaims)["exp"].(float64); ok && time.Now().Before(time.Unix(int64(expClaim), 0)) {
		return nil // Success
	}

	return fmt.Errorf("token expired or missing exp claim")
}
