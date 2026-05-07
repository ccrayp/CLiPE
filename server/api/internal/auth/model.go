package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type PrincipalType string

const (
	PrincipalHuman   PrincipalType = "human"
	PrincipalMachine PrincipalType = "machine"
)

type Principal struct {
	Type PrincipalType `json:"type"`
	ID   string        `json:"id"`
}

type Claims struct {
	Username string `json:"username"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Token     string    `gorm:"type:text;not null;uniqueIndex" json:"token"`
	Username  string    `gorm:"type:text;not null;index" json:"username"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
