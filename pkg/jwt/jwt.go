package jwt

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type claims struct {
	ID       uint
	Username string
	jwt.RegisteredClaims
}

func GenerateToken(id uint, username string) (string, error) {
	signingKey := os.Getenv("JWT_SIG_KEY")
	if signingKey == "" {
		return "", fmt.Errorf("no jwt signing key in env variables")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})

	return token.SignedString([]byte(signingKey))
}

func ParseToken(accesToken string) (id uint, err error) {
	signingKey := os.Getenv("JWT_SIG_KEY")
	if signingKey == "" {
		return 0, fmt.Errorf("no jwt signing key in env variables")
	}

	token, err := jwt.ParseWithClaims(accesToken, &claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	user_claims, ok := token.Claims.(*claims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	return user_claims.ID, nil
}

func IdentifyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			r = r.WithContext(context.WithValue(r.Context(), "is_identified", false))
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		userId, err := ParseToken(headerParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "is_identified", true))
		r = r.WithContext(context.WithValue(r.Context(), "user_id", userId))

		next.ServeHTTP(w, r)
	})
}
