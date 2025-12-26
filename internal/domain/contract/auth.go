package contract

import (
	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/Ablebil/lathi-be/pkg/response"
)

type AuthUsecaseItf interface {
	Register(req *dto.RegisterRequest) *response.APIError
	Verify(req *dto.VerifyRequest) *response.APIError
	Login(req *dto.LoginRequest) (string, string, *response.APIError)
	Refresh(req *dto.RefreshRequest) (string, string, *response.APIError)
	Logout(req *dto.LogoutRequest) *response.APIError
}
