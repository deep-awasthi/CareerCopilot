package auth

import "context"

// Repository interface
type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uint) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUUID(ctx context.Context, uuid string) (*User, error)
	FindByResetToken(ctx context.Context, token string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// Service interface
type Service interface {
	Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*AuthResponse, error)
	RequestPasswordReset(ctx context.Context, req *PasswordResetRequest) error
	ResetPassword(ctx context.Context, req *PasswordResetConfirm) error
	ChangePassword(ctx context.Context, userID uint, req *ChangePasswordRequest) error
	GetProfile(ctx context.Context, userID uint) (*UserDTO, error)
	Logout(ctx context.Context, userID uint, refreshToken string) error
}
