package interview

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/deepawasthi/careercopilot/pkg/config"
	"github.com/deepawasthi/careercopilot/pkg/middleware"
	"github.com/deepawasthi/careercopilot/pkg/response"
	"github.com/deepawasthi/careercopilot/pkg/validator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ---- Entities ----

type Stage string
type Result string

const (
	StageApplied          Stage = "applied"
	StageRecruiterCall    Stage = "recruiter_call"
	StageOnlineAssessment Stage = "online_assessment"
	StageTechnical1       Stage = "technical_round_1"
	StageTechnical2       Stage = "technical_round_2"
	StageSystemDesign     Stage = "system_design"
	StageManager          Stage = "manager_round"
	StageHR               Stage = "hr_round"
	StageOffer            Stage = "offer"
	StageRejected         Stage = "rejected"
)

const (
	ResultPending   Result = "pending"
	ResultPassed    Result = "passed"
	ResultFailed    Result = "failed"
	ResultCancelled Result = "cancelled"
)

type Interview struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	ApplicationID uint           `gorm:"uniqueIndex;not null" json:"application_id"`
	CurrentStage  Stage          `gorm:"type:interview_stage;default:'applied'" json:"current_stage"`
	Notes         string         `gorm:"type:text" json:"notes"`
	Rounds        []Round        `gorm:"foreignKey:InterviewID" json:"rounds,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Interview) TableName() string { return "interviews" }

type Round struct {
	ID                 uint           `gorm:"primarykey" json:"id"`
	InterviewID        uint           `gorm:"not null;index" json:"interview_id"`
	Stage              Stage          `gorm:"type:interview_stage;not null" json:"stage"`
	ScheduledAt        *time.Time     `json:"scheduled_at"`
	DurationMinutes    int            `json:"duration_minutes"`
	InterviewerName    string         `json:"interviewer_name"`
	InterviewerRole    string         `json:"interviewer_role"`
	InterviewerLinkedin string        `json:"interviewer_linkedin"`
	Feedback           string         `gorm:"type:text" json:"feedback"`
	Notes              string         `gorm:"type:text" json:"notes"`
	Result             Result         `gorm:"type:interview_result;default:'pending'" json:"result"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
}

func (Round) TableName() string { return "interview_rounds" }

// ---- DTOs ----

type CreateInterviewRequest struct {
	ApplicationID uint   `json:"application_id" validate:"required"`
	Notes         string `json:"notes"`
}

type AddRoundRequest struct {
	Stage              string    `json:"stage" validate:"required"`
	ScheduledAt        *time.Time `json:"scheduled_at"`
	DurationMinutes    int       `json:"duration_minutes"`
	InterviewerName    string    `json:"interviewer_name"`
	InterviewerRole    string    `json:"interviewer_role"`
	InterviewerLinkedin string   `json:"interviewer_linkedin"`
	Feedback           string    `json:"feedback"`
	Notes              string    `json:"notes"`
	Result             string    `json:"result"`
}

type UpdateRoundRequest struct {
	ScheduledAt        *time.Time `json:"scheduled_at"`
	DurationMinutes    int        `json:"duration_minutes"`
	InterviewerName    string     `json:"interviewer_name"`
	Feedback           string     `json:"feedback"`
	Notes              string     `json:"notes"`
	Result             string     `json:"result"`
}

// ---- Repository ----

type Repository interface {
	Create(ctx context.Context, interview *Interview) error
	FindByApplicationID(ctx context.Context, appID, userID uint) (*Interview, error)
	FindByID(ctx context.Context, id, userID uint) (*Interview, error)
	ListByUser(ctx context.Context, userID uint) ([]*Interview, error)
	Update(ctx context.Context, interview *Interview) error
	AddRound(ctx context.Context, round *Round) error
	UpdateRound(ctx context.Context, round *Round) error
	FindUpcoming(ctx context.Context, userID uint, limit int) ([]*Round, error)
}

type repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repository{db: db} }

func (r *repository) Create(ctx context.Context, interview *Interview) error {
	return r.db.WithContext(ctx).Create(interview).Error
}

func (r *repository) FindByApplicationID(ctx context.Context, appID, userID uint) (*Interview, error) {
	var iv Interview
	result := r.db.WithContext(ctx).
		Where("application_id = ? AND user_id = ?", appID, userID).
		Preload("Rounds").First(&iv)
	if result.Error != nil {
		return nil, fmt.Errorf("interview not found")
	}
	return &iv, nil
}

func (r *repository) FindByID(ctx context.Context, id, userID uint) (*Interview, error) {
	var iv Interview
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Preload("Rounds").First(&iv)
	if result.Error != nil {
		return nil, fmt.Errorf("interview not found")
	}
	return &iv, nil
}

func (r *repository) ListByUser(ctx context.Context, userID uint) ([]*Interview, error) {
	var ivs []*Interview
	r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Rounds").Find(&ivs)
	return ivs, nil
}

func (r *repository) Update(ctx context.Context, interview *Interview) error {
	return r.db.WithContext(ctx).Save(interview).Error
}

func (r *repository) AddRound(ctx context.Context, round *Round) error {
	return r.db.WithContext(ctx).Create(round).Error
}

func (r *repository) UpdateRound(ctx context.Context, round *Round) error {
	return r.db.WithContext(ctx).Save(round).Error
}

func (r *repository) FindUpcoming(ctx context.Context, userID uint, limit int) ([]*Round, error) {
	var rounds []*Round
	r.db.WithContext(ctx).
		Joins("JOIN interviews iv ON iv.id = interview_rounds.interview_id").
		Where("iv.user_id = ? AND interview_rounds.scheduled_at > NOW()", userID).
		Order("interview_rounds.scheduled_at ASC").
		Limit(limit).Find(&rounds)
	return rounds, nil
}

// ---- Service ----

type Service interface {
	Create(ctx context.Context, userID uint, req *CreateInterviewRequest) (*Interview, error)
	Get(ctx context.Context, id, userID uint) (*Interview, error)
	ListByUser(ctx context.Context, userID uint) ([]*Interview, error)
	AddRound(ctx context.Context, interviewID, userID uint, req *AddRoundRequest) (*Round, error)
	UpdateRound(ctx context.Context, roundID, userID uint, req *UpdateRoundRequest) (*Round, error)
	GetUpcoming(ctx context.Context, userID uint) ([]*Round, error)
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Create(ctx context.Context, userID uint, req *CreateInterviewRequest) (*Interview, error) {
	iv := &Interview{
		UserID:        userID,
		ApplicationID: req.ApplicationID,
		CurrentStage:  StageApplied,
		Notes:         req.Notes,
	}
	return iv, s.repo.Create(ctx, iv)
}

func (s *service) Get(ctx context.Context, id, userID uint) (*Interview, error) {
	return s.repo.FindByID(ctx, id, userID)
}

func (s *service) ListByUser(ctx context.Context, userID uint) ([]*Interview, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *service) AddRound(ctx context.Context, interviewID, userID uint, req *AddRoundRequest) (*Round, error) {
	iv, err := s.repo.FindByID(ctx, interviewID, userID)
	if err != nil { return nil, err }

	round := &Round{
		InterviewID:         iv.ID,
		Stage:               Stage(req.Stage),
		ScheduledAt:         req.ScheduledAt,
		DurationMinutes:     req.DurationMinutes,
		InterviewerName:     req.InterviewerName,
		InterviewerRole:     req.InterviewerRole,
		InterviewerLinkedin: req.InterviewerLinkedin,
		Feedback:            req.Feedback,
		Notes:               req.Notes,
		Result:              Result(req.Result),
	}
	if round.Result == "" { round.Result = ResultPending }

	if err := s.repo.AddRound(ctx, round); err != nil { return nil, err }

	// Update current stage
	iv.CurrentStage = round.Stage
	_ = s.repo.Update(ctx, iv)

	return round, nil
}

func (s *service) UpdateRound(ctx context.Context, roundID, userID uint, req *UpdateRoundRequest) (*Round, error) {
	var round Round
	round.ID = roundID
	if req.ScheduledAt != nil { round.ScheduledAt = req.ScheduledAt }
	if req.DurationMinutes > 0 { round.DurationMinutes = req.DurationMinutes }
	if req.InterviewerName != "" { round.InterviewerName = req.InterviewerName }
	if req.Feedback != "" { round.Feedback = req.Feedback }
	if req.Notes != "" { round.Notes = req.Notes }
	if req.Result != "" { round.Result = Result(req.Result) }
	return &round, s.repo.UpdateRound(ctx, &round)
}

func (s *service) GetUpcoming(ctx context.Context, userID uint) ([]*Round, error) {
	return s.repo.FindUpcoming(ctx, userID, 10)
}

// ---- Controller ----

type Controller struct{ svc Service }

func NewController(svc Service) *Controller { return &Controller{svc: svc} }

func (c *Controller) Create(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	var req CreateInterviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.BadRequest(ctx, err.Error()); return }
	if err := validator.Validate(&req); err != nil { response.BadRequest(ctx, err.Error()); return }
	result, err := c.svc.Create(ctx.Request.Context(), userID, &req)
	if err != nil { response.Error(ctx, err); return }
	response.Created(ctx, "interview created", result)
}

func (c *Controller) List(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	ivs, err := c.svc.ListByUser(ctx.Request.Context(), userID)
	if err != nil { response.Error(ctx, err); return }
	response.Success(ctx, "interviews retrieved", ivs)
}

func (c *Controller) Get(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	iv, err := c.svc.Get(ctx.Request.Context(), uint(id), userID)
	if err != nil { response.NotFound(ctx, "interview not found"); return }
	response.Success(ctx, "interview retrieved", iv)
}

func (c *Controller) AddRound(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	var req AddRoundRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.BadRequest(ctx, err.Error()); return }
	round, err := c.svc.AddRound(ctx.Request.Context(), uint(id), userID, &req)
	if err != nil { response.Error(ctx, err); return }
	response.Created(ctx, "round added", round)
}

func (c *Controller) UpdateRound(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	roundID, _ := strconv.ParseUint(ctx.Param("round_id"), 10, 64)
	var req UpdateRoundRequest
	if err := ctx.ShouldBindJSON(&req); err != nil { response.BadRequest(ctx, err.Error()); return }
	round, err := c.svc.UpdateRound(ctx.Request.Context(), uint(roundID), userID, &req)
	if err != nil { response.Error(ctx, err); return }
	response.Success(ctx, "round updated", round)
}

func (c *Controller) Upcoming(ctx *gin.Context) {
	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok { response.Unauthorized(ctx, "unauthorized"); return }
	rounds, err := c.svc.GetUpcoming(ctx.Request.Context(), userID)
	if err != nil { response.Error(ctx, err); return }
	response.Success(ctx, "upcoming interviews", rounds)
}

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller, cfg *config.Config) {
	iv := r.Group("/interviews")
	iv.Use(middleware.JWTAuth(&cfg.JWT))
	{
		iv.POST("", ctrl.Create)
		iv.GET("", ctrl.List)
		iv.GET("/upcoming", ctrl.Upcoming)
		iv.GET("/:id", ctrl.Get)
		iv.POST("/:id/rounds", ctrl.AddRound)
		iv.PUT("/:id/rounds/:round_id", ctrl.UpdateRound)
	}
}
