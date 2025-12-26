package jwt

import (
	"time"

	"github.com/Ablebil/lathi-be/internal/config"
	j "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtItf interface {
	CreateAccessToken(userID uuid.UUID, username, email string, exp time.Duration) (string, error)
	CreateRefreshToken(userID uuid.UUID, exp time.Duration) (string, error)
	ParseAccessToken(tokenStr string) (*AccessClaims, error)
	ParseRefreshToken(tokenStr string) (*RefreshClaims, error)
}

type AccessClaims struct {
	j.RegisteredClaims
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RefreshClaims struct {
	j.RegisteredClaims
}

type jwt struct {
	accessSecret  []byte
	refreshSecret []byte
}

func NewJwt(env *config.Env) JwtItf {
	return &jwt{
		accessSecret:  []byte(env.AccessSecret),
		refreshSecret: []byte(env.RefreshSecret),
	}
}

func (jw *jwt) CreateAccessToken(userID uuid.UUID, username, email string, exp time.Duration) (string, error) {
	claims := &AccessClaims{
		RegisteredClaims: j.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: j.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  j.NewNumericDate(time.Now()),
		},
		Username: username,
		Email:    email,
	}

	token := j.NewWithClaims(j.SigningMethodHS256, claims)
	return token.SignedString(jw.accessSecret)
}

func (jw *jwt) CreateRefreshToken(userID uuid.UUID, exp time.Duration) (string, error) {
	claims := &RefreshClaims{
		RegisteredClaims: j.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: j.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  j.NewNumericDate(time.Now()),
		},
	}

	token := j.NewWithClaims(j.SigningMethodHS256, claims)
	return token.SignedString(jw.refreshSecret)
}

func (jw *jwt) ParseAccessToken(tokenStr string) (*AccessClaims, error) {
	claims := &AccessClaims{}
	token, err := j.ParseWithClaims(tokenStr, claims, func(token *j.Token) (any, error) {
		return jw.accessSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func (jw *jwt) ParseRefreshToken(tokenStr string) (*RefreshClaims, error) {
	claims := &RefreshClaims{}
	token, err := j.ParseWithClaims(tokenStr, claims, func(token *j.Token) (any, error) {
		return jw.refreshSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
