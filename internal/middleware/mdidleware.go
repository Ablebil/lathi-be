package middleware

import "github.com/Ablebil/lathi-be/pkg/jwt"

type middleware struct {
	jwt jwt.JwtItf
}

func NewMiddleware(jwt jwt.JwtItf) *middleware {
	return &middleware{
		jwt: jwt,
	}
}
