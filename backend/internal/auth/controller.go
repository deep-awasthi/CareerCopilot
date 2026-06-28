package auth

import (
	"github.com/deepawasthi/careercopilot/pkg/middleware"
	"github.com/deepawasthi/careercopilot/pkg/response"
	"github.com/deepawasthi/careercopilot/pkg/validator"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	svc Service
}

func NewController(svc Service) *Controller {
	return &Controller{svc: svc}
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "Registration payload"
// @Success 201 {object} response.Response
// @Router /api/v1/auth/register [post]
func (c *Controller) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "invalid request body")
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	result, err := c.svc.Register(ctx.Request.Context(), &req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Created(ctx, "user registered successfully", result)
}

// Login godoc
// @Summary Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login payload"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/login [post]
func (c *Controller) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "invalid request body")
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}
	result, err := c.svc.Login(ctx.Request.Context(), &req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "login successful", result)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/refresh [post]
func (c *Controller) RefreshToken(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "invalid request body")
		return
	}
	result, err := c.svc.RefreshToken(ctx.Request.Context(), &req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "token refreshed", result)
}

// RequestPasswordReset godoc
// @Summary Request password reset email
// @Tags auth
// @Accept json
// @Produce json
// @Param body body PasswordResetRequest true "Email"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/password-reset [post]
func (c *Controller) RequestPasswordReset(ctx *gin.Context) {
	var req PasswordResetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "invalid request body")
		return
	}
	_ = c.svc.RequestPasswordReset(ctx.Request.Context(), &req)
	response.Success(ctx, "if the email exists, a reset link will be sent", nil)
}

// ResetPassword godoc
// @Summary Reset password with token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body PasswordResetConfirm true "Reset payload"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/password-reset/confirm [post]
func (c *Controller) ResetPassword(ctx *gin.Context) {
	var req PasswordResetConfirm
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "invalid request body")
		return
	}
	if err := c.svc.ResetPassword(ctx.Request.Context(), &req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "password reset successfully", nil)
}

// ChangePassword godoc
// @Summary Change password (authenticated)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body ChangePasswordRequest true "Change password payload"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/change-password [post]
func (c *Controller) ChangePassword(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		response.Unauthorized(ctx, "unauthorized")
		return
	}
	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "invalid request body")
		return
	}
	if err := c.svc.ChangePassword(ctx.Request.Context(), userID, &req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "password changed successfully", nil)
}

// Me godoc
// @Summary Get current authenticated user
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/v1/auth/me [get]
func (c *Controller) Me(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		response.Unauthorized(ctx, "unauthorized")
		return
	}
	user, err := c.svc.GetProfile(ctx.Request.Context(), userID)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "user profile", user)
}

// Logout godoc
// @Summary Logout (invalidate refresh token)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Router /api/v1/auth/logout [post]
func (c *Controller) Logout(ctx *gin.Context) {
	userID, _ := middleware.GetUserIDFromContext(ctx)
	var req RefreshTokenRequest
	_ = ctx.ShouldBindJSON(&req)
	_ = c.svc.Logout(ctx.Request.Context(), userID, req.RefreshToken)
	response.Success(ctx, "logged out successfully", nil)
}
