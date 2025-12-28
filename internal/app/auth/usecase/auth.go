package usecase

import (
	"context"
	"fmt"
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
	jwt    jwt.JwtItf
	env    *config.Env
}

func NewAuthUsecase(userRepo contract.UserRepositoryItf, bcrypt bcrypt.BcryptItf, mail mail.MailItf, cache redis.RedisItf, jwt jwt.JwtItf, env *config.Env) contract.AuthUsecaseItf {
	return &authUsecase{
		repo:   userRepo,
		bcrypt: bcrypt,
		mail:   mail,
		cache:  cache,
		jwt:    jwt,
		env:    env,
	}
}

func (uc *authUsecase) Register(req *dto.RegisterRequest) *response.APIError {
	user, err := uc.repo.GetUserByEmail(req.Email)
	if err != nil {
		return response.ErrInternal("failed to find user")
	}
	if user != nil {
		return response.ErrConflict("email already registered")
	}

	hashed, err := uc.bcrypt.Hash(req.Password)
	if err != nil {
		return response.ErrInternal("failed to hash password")
	}

	newUser := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashed,
	}
	if err := uc.repo.CreateUser(newUser); err != nil {
		return response.ErrInternal("failed to create user")
	}

	// generate verif token
	token := uuid.NewString()
	cacheKey := fmt.Sprintf("verify:%s", token)
	ctx := context.Background()
	if err := uc.cache.Set(ctx, cacheKey, newUser.Email, uc.env.VerifTokenTtl); err != nil {
		return response.ErrInternal("failed to store verification token")
	}

	u, _ := url.Parse(uc.env.VerifUrl)
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()

	verifyUrl := u.String()
	mailData := map[string]any{
		"Username":      newUser.Username,
		"VerifyURL":     verifyUrl,
		"ExpireMinutes": uc.env.VerifTokenTtl.Minutes(),
	}
	if err := uc.mail.Send(newUser.Email, "Verifikasi Email", "verification.html", mailData); err != nil {
		return response.ErrInternal("failed to send verification email")
	}

	return nil
}

func (uc *authUsecase) Verify(req *dto.VerifyRequest) *response.APIError {
	cacheKey := fmt.Sprintf("verify:%s", req.Token)
	ctx := context.Background()
	var email string
	if err := uc.cache.Get(ctx, cacheKey, &email); err != nil {
		return response.ErrBadRequest("invalid or expired token")
	}

	user, err := uc.repo.GetUserByEmail(email)
	if err != nil {
		return response.ErrNotFound("failed to find user")
	}
	if user == nil {
		return response.ErrNotFound("user not found")
	}
	if user.IsVerified {
		return response.ErrBadRequest("user already verified")
	}

	user.IsVerified = true
	if err := uc.repo.UpdateUser(user); err != nil {
		return response.ErrInternal("failed to update user verification")
	}

	_ = uc.cache.Del(ctx, cacheKey)
	return nil
}

func (uc *authUsecase) Login(req *dto.LoginRequest) (string, string, *response.APIError) {
	user, err := uc.repo.GetUserByEmail(req.Email)
	if err != nil {
		return "", "", response.ErrInternal("failed to find user")
	}
	if user == nil || !uc.bcrypt.Compare(req.Password, user.Password) {
		return "", "", response.ErrUnauthorized("invalid email or password")
	}
	if !user.IsVerified {
		return "", "", response.ErrUnauthorized("email not verified")
	}

	accessToken, err := uc.jwt.CreateAccessToken(user.ID, user.Username, user.Email, uc.env.AccessTtl)
	if err != nil {
		return "", "", response.ErrInternal("failed to create access token")
	}

	refreshToken, err := uc.jwt.CreateRefreshToken(user.ID, uc.env.RefreshTtl)
	if err != nil {
		return "", "", response.ErrInternal("failed to create refresh token")
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("refresh:%s", refreshToken)
	if err := uc.cache.Set(ctx, cacheKey, user.ID.String(), uc.env.RefreshTtl); err != nil {
		return "", "", response.ErrInternal("failed to store refresh token")
	}

	return accessToken, refreshToken, nil
}

func (uc *authUsecase) Refresh(req *dto.RefreshRequest) (string, string, *response.APIError) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("refresh:%s", req.RefreshToken)
	var userID string
	if err := uc.cache.Get(ctx, cacheKey, &userID); err != nil {
		return "", "", response.ErrUnauthorized("invalid or expired refresh token")
	}

	user, err := uc.repo.GetUserByID(uuid.MustParse(userID))
	if err != nil || user == nil {
		return "", "", response.ErrUnauthorized("invalid refresh or expired refresh token")
	}

	newAccessToken, err := uc.jwt.CreateAccessToken(user.ID, user.Username, user.Email, uc.env.AccessTtl)
	if err != nil {
		return "", "", response.ErrInternal("failed to create access token")
	}

	newRefreshToken, err := uc.jwt.CreateRefreshToken(user.ID, uc.env.RefreshTtl)
	if err != nil {
		return "", "", response.ErrInternal("failed to create refresh token")
	}

	_ = uc.cache.Del(ctx, cacheKey)
	newCacheKey := fmt.Sprintf("refresh:%s", newRefreshToken)
	if err := uc.cache.Set(ctx, newCacheKey, user.ID.String(), uc.env.RefreshTtl); err != nil {
		return "", "", response.ErrInternal("failed to store refresh token")
	}

	return newAccessToken, newRefreshToken, nil
}

func (uc *authUsecase) Logout(req *dto.LogoutRequest) *response.APIError {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("refresh:%s", req.RefreshToken)
	if err := uc.cache.Del(ctx, cacheKey); err != nil {
		return response.ErrInternal("failed to delete refresh token")
	}

	return nil
}
