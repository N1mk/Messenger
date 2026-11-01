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

type accesTokenClaims struct {
	ID       int
	Username string
	jwt.RegisteredClaims
}

type refreshTokenClaims struct {
	ID int
	jwt.RegisteredClaims
}

type JWT struct{}

func NewJWT() *JWT {
	return &JWT{}
}

func (j *JWT) GenerateTokenPair(id int, username string) (accesToken string, refreshToken string, err error) {
	signingKey := os.Getenv("JWT_SIG_KEY")
	if signingKey == "" {
		return "", "", fmt.Errorf("no jwt signing key in env variables")
	}

	unsignedAccesToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &accesTokenClaims{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})

	unsignedRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &refreshTokenClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 720)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	},
	)

	accesToken, err = unsignedAccesToken.SignedString([]byte(signingKey))
	if err != nil {
		return "", "", err
	}

	refreshToken, err = unsignedRefreshToken.SignedString([]byte(signingKey))
	if err != nil {
		return "", "", err
	}

	return accesToken, refreshToken, nil
}

func (j *JWT) ParseAccesToken(accesToken string) (id int, username string, err error) {
	signingKey := os.Getenv("JWT_SIG_KEY")
	if signingKey == "" {
		return 0, "", fmt.Errorf("no jwt signing key in env variables")
	}

	token, err := jwt.ParseWithClaims(string(accesToken), &accesTokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, "", err
	}

	tokenClaims, ok := token.Claims.(*accesTokenClaims)
	if !ok {
		return 0, "", fmt.Errorf("invalid token claims")
	}

	return tokenClaims.ID, tokenClaims.Username, nil
}

func (j *JWT) ParseRefreshToken(refreshToken string) (id int, err error) {
	signingKey := os.Getenv("JWT_SIG_KEY")
	if signingKey == "" {
		return 0, fmt.Errorf("no jwt signing key in env variables")
	}

	token, err := jwt.ParseWithClaims(string(refreshToken), &refreshTokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	tokenClaims, ok := token.Claims.(*refreshTokenClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	return tokenClaims.ID, err
}

func (j *JWT) IdentifyMiddleware(next http.Handler) http.Handler {
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

		userID, username, err := j.ParseAccesToken(headerParts[1])
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "is_identified", true))
		r = r.WithContext(context.WithValue(r.Context(), "user_id", userID))
		r = r.WithContext(context.WithValue(r.Context(), "username", username))

		next.ServeHTTP(w, r)
	})
}
