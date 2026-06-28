package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/deepawasthi/careercopilot/pkg/config"
	"github.com/deepawasthi/careercopilot/pkg/email"
	"github.com/deepawasthi/careercopilot/pkg/errors"
	"github.com/deepawasthi/careercopilot/pkg/logger"
	"github.com/deepawasthi/careercopilot/pkg/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo        Repository
	cfg         *config.Config
	emailClient *email.Client
}

func NewService(repo Repository, cfg *config.Config, emailClient *email.Client) Service {
	return &service{repo: repo, cfg: cfg, emailClient: emailClient}
}

func (s *service) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	exists, err := s.repo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.InternalWrap("failed to check existing user", err)
	}
	if exists {
		return nil, errors.ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.InternalWrap("failed to hash password", err)
	}

	user := &User{
		UUID:         uuid.New().String(),
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, errors.InternalWrap("failed to create user", err)
	}

	logger.Info("User registered", zap.String("email", user.Email), zap.Uint("id", user.ID))
	return s.buildAuthResponse(user)
}

func (s *service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, errors.Unauthorized("account is deactivated")
	}

	now := time.Now()
	user.LastLoginAt = &now
	_ = s.repo.Update(ctx, user)

	logger.Info("User logged in", zap.String("email", user.Email))
	return s.buildAuthResponse(user)
}

func (s *service) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*AuthResponse, error) {
	claims, err := middleware.ParseRefreshToken(req.RefreshToken, &s.cfg.JWT)
	if err != nil {
		return nil, errors.ErrTokenInvalid
	}

	user, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	return s.buildAuthResponse(user)
}

func (s *service) RequestPasswordReset(ctx context.Context, req *PasswordResetRequest) error {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		// Return nil to not reveal user existence
		return nil
	}

	token, err := generateSecureToken(32)
	if err != nil {
		return errors.Internal("failed to generate reset token")
	}

	expires := time.Now().Add(1 * time.Hour)
	user.PasswordResetToken = token
	user.PasswordResetExpires = &expires

	if err := s.repo.Update(ctx, user); err != nil {
		return errors.InternalWrap("failed to save reset token", err)
	}

	resetURL := fmt.Sprintf("%s/password-reset/confirm?token=%s", s.cfg.Frontend.URL, token)
	if err := s.emailClient.SendPasswordReset(user.Email, resetURL); err != nil {
		return errors.InternalWrap("failed to send password reset email", err)
	}

	logger.Info("Password reset requested", zap.String("email", req.Email))
	return nil
}

func (s *service) ResetPassword(ctx context.Context, req *PasswordResetConfirm) error {
	user, err := s.repo.FindByResetToken(ctx, req.Token)
	if err != nil {
		return errors.BadRequest("invalid or expired reset token")
	}

	if user.PasswordResetExpires == nil || time.Now().After(*user.PasswordResetExpires) {
		return errors.BadRequest("reset token has expired")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Internal("failed to hash password")
	}

	user.PasswordHash = string(hash)
	user.PasswordResetToken = ""
	user.PasswordResetExpires = nil

	return s.repo.Update(ctx, user)
}

func (s *service) ChangePassword(ctx context.Context, userID uint, req *ChangePasswordRequest) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return errors.BadRequest("current password is incorrect")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.Internal("failed to hash new password")
	}

	user.PasswordHash = string(hash)
	return s.repo.Update(ctx, user)
}

func (s *service) GetProfile(ctx context.Context, userID uint) (*UserDTO, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	dto := ToUserDTO(user)
	return &dto, nil
}

func (s *service) Logout(ctx context.Context, userID uint, refreshToken string) error {
	// Token invalidation is handled client-side; optionally add to Redis blacklist
	logger.Info("User logged out", zap.Uint("user_id", userID))
	return nil
}

func (s *service) buildAuthResponse(user *User) (*AuthResponse, error) {
	accessToken, err := middleware.GenerateAccessToken(user.ID, user.Email, &s.cfg.JWT)
	if err != nil {
		return nil, errors.InternalWrap("failed to generate access token", err)
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.ID, user.Email, &s.cfg.JWT)
	if err != nil {
		return nil, errors.InternalWrap("failed to generate refresh token", err)
	}

	expiresIn := int(s.cfg.JWT.AccessExpiry.Seconds())
	if expiresIn == 0 {
		expiresIn = 900
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		User:         ToUserDTO(user),
	}, nil
}

func generateSecureToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return hex.EncodeToString(b), nil
}
