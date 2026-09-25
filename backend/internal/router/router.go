package router

import (
	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/campusbooks/campusbooks/internal/config"
	"github.com/campusbooks/campusbooks/internal/dto"
	"github.com/campusbooks/campusbooks/internal/handler"
	"github.com/campusbooks/campusbooks/internal/middleware"
	"github.com/campusbooks/campusbooks/internal/repository"
	"github.com/campusbooks/campusbooks/internal/service"
	"github.com/campusbooks/campusbooks/internal/util"
)

// Setup builds the gin engine and wires all dependencies.
func Setup(cfg *config.Config, db *gorm.DB, codeStore util.CodeStore, minio *util.MinIOClient, logger *slog.Logger) *gin.Engine {
	userRepo := repository.NewUserRepository(db)
	bookRepo := repository.NewBookRepository(db)
	wishRepo := repository.NewWishRepository(db)
	convRepo := repository.NewConversationRepository(db)
	msgRepo := repository.NewMessageRepository(db)
	evalRepo := repository.NewEvaluationRepository(db)
	favRepo := repository.NewFavoriteRepository(db)
	historyRepo := repository.NewBrowseHistoryRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	borrowRepo := repository.NewBorrowRepository(db)

	authService := service.NewAuthService(userRepo, codeStore, logger, cfg)
	userService := service.NewUserService(userRepo, evalRepo, logger)
	bookService := service.NewBookService(db, bookRepo, favRepo, historyRepo, userRepo, borrowRepo, logger)
	wishService := service.NewWishService(wishRepo, logger)
	convService := service.NewConversationService(convRepo, msgRepo, bookRepo, wishRepo, logger)
	evalService := service.NewEvaluationService(evalRepo, bookRepo, userRepo, logger)
	auditService := service.NewAuditService(auditRepo, logger)
	uploadService := service.NewUploadService(minio, logger)
	borrowService := service.NewBorrowService(db, borrowRepo, bookRepo, convService, logger)

	authHandler := handler.NewAuthHandler(authService, logger)
	userHandler := handler.NewUserHandler(userService, logger)
	bookHandler := handler.NewBookHandler(bookService, logger)
	wishHandler := handler.NewWishHandler(wishService, logger)
	convHandler := handler.NewConversationHandler(convService, logger)
	evalHandler := handler.NewEvaluationHandler(evalService, logger)
	auditHandler := handler.NewAuditHandler(auditService, logger)
	uploadHandler := handler.NewUploadHandler(uploadService, minio, logger)
	borrowHandler := handler.NewBorrowHandler(borrowService, logger)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.RequestID())
	r.Use(middleware.RequestLogger(logger))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     splitOrigins(cfg.CORSOrigins),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: true,
	}))
	r.Use(middleware.ErrorHandler(logger))
	r.Use(middleware.AuditMiddleware(auditService, logger))

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, dto.OK(gin.H{"status": "ok"})) })

	limiter := middleware.NewRateLimiter(cfg.RateLimitReq, cfg.RateLimitWin, logger)
	v1 := r.Group("/api/v1")
	{
		registerAuthRoutes(v1, cfg, authHandler, limiter)
		registerUserRoutes(v1, cfg, userHandler, evalHandler, limiter)
		registerBookRoutes(v1, cfg, bookHandler, limiter)
		registerBorrowRoutes(v1, cfg, borrowHandler, limiter)
		registerWishRoutes(v1, cfg, wishHandler, convHandler, limiter)
		registerConversationRoutes(v1, cfg, convHandler, limiter)
		registerEvaluationRoutes(v1, cfg, evalHandler, limiter)
		registerAuditRoutes(v1, cfg, auditHandler, limiter)
		registerUploadRoutes(v1, cfg, uploadHandler, limiter)
	}
	return r
}

func splitOrigins(s string) []string {
	if s == "" {
		return []string{"http://localhost:8011"}
	}
	origins := []string{}
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			if seg := s[start:i]; seg != "" {
				origins = append(origins, seg)
			}
			start = i + 1
		}
	}
	return origins
}
