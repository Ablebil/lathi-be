package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/Ablebil/lathi-be/internal/config"
	"github.com/Ablebil/lathi-be/internal/domain/contract"
	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"github.com/Ablebil/lathi-be/internal/infra/redis"
	"github.com/Ablebil/lathi-be/pkg/bcrypt"
	"github.com/Ablebil/lathi-be/pkg/jwt"
	"github.com/Ablebil/lathi-be/pkg/mail"
	"github.com/Ablebil/lathi-be/pkg/response"
	"github.com/google/uuid"
)

type authUsecase struct {
	repo   contract.UserRepositoryItf
	bcrypt bcrypt.BcryptItf
	mail   mail.MailItf
	cache  redis.RedisItf
	jwt    jwt.JWTItf
	env    *config.Env
}

func NewAuthUsecase(userRepo contract.UserRepositoryItf, bcrypt bcrypt.BcryptItf, mail mail.MailItf, cache redis.RedisItf, jwt jwt.JWTItf, env *config.Env) contract.AuthUsecaseItf {
	return &authUsecase{
		repo:   userRepo,
		bcrypt: bcrypt,
		mail:   mail,
		cache:  cache,
		jwt:    jwt,
		env:    env,
	}
}

func (uc *authUsecase) Register(ctx context.Context, req *dto.RegisterRequest) *response.APIError {
	user, err := uc.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		slog.Error("failed to get user", "error", err)
		return response.ErrInternal("Coba lagi nanti ya!")
	}
	if user != nil {
		return response.ErrConflict("Email ini udah pernah didaftarin, coba email lain ya")
	}

	userByUsn, err := uc.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		slog.Error("failed to get user", "error", err)
		return response.ErrInternal("Coba lagi nanti ya!")
	}
	if userByUsn != nil {
		return response.ErrConflict("Username ini udah dipake, coba yang lain ya")
	}

	hashed, err := uc.bcrypt.Hash(req.Password)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return response.ErrInternal("Coba lagi nanti ya!")
	}

	newUser := &entity.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashed,
		AvatarURL: uc.env.DefaultAvatarURL,
	}
	if err := uc.repo.CreateUser(ctx, newUser); err != nil {
		slog.Error("failed to create user", "error", err)
		return response.ErrInternal("Coba lagi nanti ya!")
	}

	// generate verif token
	token := uuid.NewString()
	cacheKey := fmt.Sprintf("verify:%s", token)
	if err := uc.cache.Set(ctx, cacheKey, newUser.Email, uc.env.VerifTokenTTL); err != nil {
		slog.Error("failed to store verification token", "error", err)
		return response.ErrInternal("Coba lagi nanti ya!")
	}

	u, _ := url.Parse(uc.env.VerifURL)
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()

	verifyURL := u.String()
	mailData := map[string]any{
		"Username":      newUser.Username,
		"VerifyURL":     verifyURL,
		"ExpireMinutes": uc.env.VerifTokenTTL.Minutes(),
	}
	if err := uc.mail.Send(newUser.Email, "Verifikasi Email", "verification.html", mailData); err != nil {
		slog.Error("failed to send verification email", "error", err)
		return response.ErrInternal("Coba lagi nanti ya!")
	}

	return nil
}

func (uc *authUsecase) Verify(ctx context.Context, req *dto.VerifyRequest) *response.APIError {
	cacheKey := fmt.Sprintf("verify:%s", req.Token)
	var email string
	if err := uc.cache.Get(ctx, cacheKey, &email); err != nil {
		return response.ErrBadRequest("Token verifikasi ga valid atau udah kadaluarsa, coba daftar lagi ya")
	}

	user, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		slog.Error("failed to get user", "error", err)
		return response.ErrInternal("Coba lagi nanti ya!")
	}
	if user == nil {
		return response.ErrNotFound("Akun ga ditemukan, coba daftar dulu ya")
	}
	if user.IsVerified {
		return response.ErrBadRequest("Akunmu udah terverifikasi sebelumnya")
	}

	user.IsVerified = true
	if err := uc.repo.UpdateUser(ctx, user); err != nil {
		slog.Error("failed to update user", "error", err)
		return response.ErrInternal("Coba lagi nanti ya!")
	}

	_ = uc.cache.Del(ctx, cacheKey)
	return nil
}

func (uc *authUsecase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.TokenResponse, *response.APIError) {
	user, err := uc.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		slog.Error("failed to get user", "error", err)
		return nil, response.ErrInternal("Coba lagi nanti ya!")
	}
	if user == nil || !uc.bcrypt.Compare(req.Password, user.Password) {
		return nil, response.ErrUnauthorized("Email atau password kamu salah")
	}
	if !user.IsVerified {
		return nil, response.ErrUnauthorized("Akunmu belum terverifikasi, cek email kamu ya")
	}

	accessToken, err := uc.jwt.CreateAccessToken(user.ID, user.Username, user.Email, uc.env.AccessTTL)
	if err != nil {
		slog.Error("failed to create access token", "error", err)
		return nil, response.ErrInternal("Coba lagi nanti ya!")
	}

	refreshToken, err := uc.jwt.CreateRefreshToken(user.ID, uc.env.RefreshTTL)
	if err != nil {
		slog.Error("failed to create refresh token", "error", err)
		return nil, response.ErrInternal("Coba lagi nanti ya!")
	}

	cacheKey := fmt.Sprintf("refresh:%s", refreshToken)
	if err := uc.cache.Set(ctx, cacheKey, user.ID.String(), uc.env.RefreshTTL); err != nil {
		slog.Error("failed to store refresh token", "error", err)
		return nil, response.ErrInternal("Coba lagi nanti ya!")
	}

	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (uc *authUsecase) Refresh(ctx context.Context, refreshToken string) (*dto.TokenResponse, *response.APIError) {
	cacheKey := fmt.Sprintf("refresh:%s", refreshToken)
	var userID string
	if err := uc.cache.Get(ctx, cacheKey, &userID); err != nil {
		return nil, response.ErrUnauthorized("Sesi kamu udah habis, coba login lagi ya")
	}

	user, err := uc.repo.GetUserByID(ctx, uuid.MustParse(userID))
	if err != nil || user == nil {
		return nil, response.ErrUnauthorized("Sesi kamu udah habis, coba login lagi ya")
	}

	newAccessToken, err := uc.jwt.CreateAccessToken(user.ID, user.Username, user.Email, uc.env.AccessTTL)
	if err != nil {
		slog.Error("failed to create access token", "error", err)
		return nil, response.ErrInternal("Coba lagi nanti ya!")
	}

	newRefreshToken, err := uc.jwt.CreateRefreshToken(user.ID, uc.env.RefreshTTL)
	if err != nil {
		slog.Error("failed to create refresh token", "error", err)
		return nil, response.ErrInternal("Coba lagi nanti ya!")
	}

	_ = uc.cache.Del(ctx, cacheKey)
	newCacheKey := fmt.Sprintf("refresh:%s", newRefreshToken)
	if err := uc.cache.Set(ctx, newCacheKey, user.ID.String(), uc.env.RefreshTTL); err != nil {
		slog.Error("failed to store refresh token", "error", err)
		return nil, response.ErrInternal("Coba lagi nanti ya!")
	}

	return &dto.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (uc *authUsecase) Logout(ctx context.Context, refreshToken string) *response.APIError {
	cacheKey := fmt.Sprintf("refresh:%s", refreshToken)
	if err := uc.cache.Del(ctx, cacheKey); err != nil {
		slog.Error("failed to delete refresh token", "error", err)
		return response.ErrInternal("Coba lagi nanti ya!")
	}

	return nil
}
