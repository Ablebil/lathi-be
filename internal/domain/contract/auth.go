package contract

import (
	"context"

	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/Ablebil/lathi-be/pkg/response"
)

type AuthUsecaseItf interface {
	Register(ctx context.Context, req *dto.RegisterRequest) *response.APIError
	Verify(ctx context.Context, req *dto.VerifyRequest) *response.APIError
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.TokenResponse, *response.APIError)
	Refresh(ctx context.Context, req *dto.RefreshRequest) (*dto.TokenResponse, *response.APIError)
	Logout(ctx context.Context, req *dto.LogoutRequest) *response.APIError
}
