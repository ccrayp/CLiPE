package auth

import (
	"clipe/pkg/utils"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		Type:     "human",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(header string) (*Claims, error) {
	if !strings.HasPrefix(header, "Bearer ") {
		return nil, errors.New("invalid authorization header")
	}

	tokenStr := strings.TrimPrefix(header, "Bearer ")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func GenerateRefreshToken(repo *Repository, username string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	token := hex.EncodeToString(b)

	if err := repo.SaveRefreshToken(token, username); err != nil {
		return "", err
	}

	return token, nil
}

func RefreshAccessToken(repo *Repository, refreshToken string) (string, string, error) {
	rt, err := repo.GetRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	if err := repo.DeleteRefreshToken(refreshToken); err != nil {
		return "", "", err
	}

	accessToken, err := GenerateToken(rt.Username)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := GenerateRefreshToken(repo, rt.Username)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func BuildPrincipalFromClaims(claims *Claims) *Principal {
	return &Principal{
		Type: PrincipalHuman,
		ID:   claims.Username,
	}
}

func CheckPassword(raw, hash string) bool {
	return utils.CheckPasswordHash(raw, hash)
}
