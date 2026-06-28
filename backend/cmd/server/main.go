// @title CareerCopilot API
// @version 1.0
// @description Production-grade career management platform API
// @termsOfService https://careercopilot.io/terms

// @contact.name CareerCopilot Support
// @contact.url https://careercopilot.io/support
// @contact.email support@careercopilot.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deepawasthi/careercopilot/internal/analytics"
	"github.com/deepawasthi/careercopilot/internal/application"
	"github.com/deepawasthi/careercopilot/internal/auth"
	"github.com/deepawasthi/careercopilot/internal/bookmark"
	"github.com/deepawasthi/careercopilot/internal/company"
	"github.com/deepawasthi/careercopilot/internal/interview"
	"github.com/deepawasthi/careercopilot/internal/job"
	"github.com/deepawasthi/careercopilot/internal/keyword_alert"
	"github.com/deepawasthi/careercopilot/internal/notification"
	"github.com/deepawasthi/careercopilot/internal/provider"
	"github.com/deepawasthi/careercopilot/internal/referral"
	"github.com/deepawasthi/careercopilot/internal/resume"
	"github.com/deepawasthi/careercopilot/internal/scheduler"
	"github.com/deepawasthi/careercopilot/internal/search"
	search_profile "github.com/deepawasthi/careercopilot/internal/search_profile"
	"github.com/deepawasthi/careercopilot/internal/user"
	"github.com/deepawasthi/careercopilot/pkg/config"
	"github.com/deepawasthi/careercopilot/pkg/database"
	"github.com/deepawasthi/careercopilot/pkg/email"
	"github.com/deepawasthi/careercopilot/pkg/logger"
	"github.com/deepawasthi/careercopilot/pkg/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Load config
	cfg := config.Load()

	// Initialize logger
	logger.Init(cfg.App.Env)
	defer logger.Sync()

	logger.Info("Starting CareerCopilot API",
		zap.String("env", cfg.App.Env),
		zap.String("port", cfg.App.Port),
	)

	// Connect to PostgreSQL
	db, err := database.InitPostgres(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
	}

	// Connect to Redis (optional — don't fatal if not available)
	_, redisErr := database.InitRedis(&cfg.Redis)
	if redisErr != nil {
		logger.Warn("Redis not available, continuing without cache", zap.Error(redisErr))
	}

	// Connect to Elasticsearch (optional — graceful fallback)
	esClient, esErr := database.InitElasticsearch(&cfg.Elasticsearch)
	if esErr != nil {
		logger.Warn("Elasticsearch not available, global search will fallback to DB", zap.Error(esErr))
	}

	// Auto-migrate GORM models
	if err := autoMigrate(db); err != nil {
		logger.Fatal("Database migration failed", zap.Error(err))
	}

	// Gin router setup
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS([]string{cfg.Frontend.URL, "http://localhost:3001", "http://localhost:5173"}))
	router.Use(middleware.RateLimiter(cfg.RateLimit.Requests, cfg.RateLimit.Duration))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "CareerCopilot API",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	})

	// API v1 group
	v1 := router.Group("/api/v1")

	// Initialize Email Client
	emailClient := email.NewClient(&cfg.SMTP)

	// ---- Wire up all modules ----

	// Auth
	authRepo := auth.NewRepository(db)
	authSvc := auth.NewService(authRepo, cfg, emailClient)
	authCtrl := auth.NewController(authSvc)
	auth.RegisterRoutes(v1, authCtrl, cfg)

	// User Profile
	userRepo := user.NewRepository(db)
	userSvc := user.NewService(userRepo)
	userCtrl := user.NewController(userSvc)
	user.RegisterRoutes(v1, userCtrl, cfg)

	// Resume
	resumeRepo := resume.NewRepository(db)
	resumeSvc := resume.NewService(resumeRepo)
	resumeCtrl := resume.NewController(resumeSvc)
	resume.RegisterRoutes(v1, resumeCtrl, cfg)

	// Search Profiles
	spRepo := search_profile.NewRepository(db)
	spSvc := search_profile.NewService(spRepo)
	spCtrl := search_profile.NewController(spSvc)
	search_profile.RegisterRoutes(v1, spCtrl, cfg)

	// Jobs
	jobRepo := job.NewRepository(db)
	jobSvc := job.NewService(jobRepo, userRepo)
	jobCtrl := job.NewController(jobSvc)
	job.RegisterRoutes(v1, jobCtrl, cfg)

	// Applications
	appRepo := application.NewRepository(db)
	appSvc := application.NewService(appRepo)
	appCtrl := application.NewController(appSvc)
	application.RegisterRoutes(v1, appCtrl, cfg)

	// Interviews
	ivRepo := interview.NewRepository(db)
	ivSvc := interview.NewService(ivRepo)
	ivCtrl := interview.NewController(ivSvc)
	interview.RegisterRoutes(v1, ivCtrl, cfg)

	// Companies & Watchlists
	coRepo := company.NewRepository(db)
	coSvc := company.NewService(coRepo)
	coCtrl := company.NewController(coSvc)
	company.RegisterRoutes(v1, coCtrl, cfg)

	// Referrals
	refRepo := referral.NewRepository(db)
	refSvc := referral.NewService(refRepo)
	refCtrl := referral.NewController(refSvc)
	referral.RegisterRoutes(v1, refCtrl, cfg)

	// Bookmarks
	bmRepo := bookmark.NewRepository(db)
	bmSvc := bookmark.NewService(bmRepo)
	bmCtrl := bookmark.NewController(bmSvc)
	bookmark.RegisterRoutes(v1, bmCtrl, cfg)

	// Keyword Alerts
	kaRepo := keyword_alert.NewRepository(db)
	kaSvc := keyword_alert.NewService(kaRepo)
	kaCtrl := keyword_alert.NewController(kaSvc)
	keyword_alert.RegisterRoutes(v1, kaCtrl, cfg)

	// Notifications
	notifRepo := notification.NewRepository(db)
	notifSvc := notification.NewService(notifRepo)
	notifCtrl := notification.NewController(notifSvc)
	notification.RegisterRoutes(v1, notifCtrl, cfg)

	// Analytics
	analyticsSvc := analytics.NewService(db)
	analyticsCtrl := analytics.NewController(analyticsSvc)
	analytics.RegisterRoutes(v1, analyticsCtrl, cfg)

	// Elasticsearch Search
	if esClient != nil {
		searchSvc := search.NewSearchService(esClient)
		searchCtrl := search.NewController(searchSvc)
		search.RegisterRoutes(v1, searchCtrl, cfg)
		// Create ES index
		_ = searchSvc.CreateIndex(context.Background())
	}

	// Swagger docs
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ---- Start background scheduler ----
	providerRegistry := provider.NewDefaultRegistry()
	sched := scheduler.NewScheduler(db, cfg, providerRegistry)
	sched.Start()

	// ---- Start HTTP server with graceful shutdown ----
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("HTTP server listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down gracefully...")

	// Stop scheduler first
	sched.Stop()

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced shutdown", zap.Error(err))
	}

	logger.Info("CareerCopilot API stopped")
}

// autoMigrate runs GORM auto-migration for all entities
func autoMigrate(db interface{ AutoMigrate(...interface{}) error }) error {
	// All tables are created and managed via SQL migrations in backend/migrations.
	return nil
}
