package handlers

import (
	"crowdreview/config"
	"crowdreview/internal/services"
	"crowdreview/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RouterDeps bundles external dependencies for handlers.
type RouterDeps struct {
	Config   config.Config
	Services services.Services
	Redis    *redis.Client
}

// SetupRouter initializes Gin with routes and middleware.
func SetupRouter(deps RouterDeps) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())
	if deps.Redis != nil {
		r.Use(middleware.RateLimitMiddleware(deps.Redis, deps.Config))
	}

	authHandler := NewAuthHandler(deps.Services.Auth, deps.Config)
	companyHandler := NewCompanyHandler(deps.Services.Company)
	reviewHandler := NewReviewHandler(deps.Services.Review)
	adminHandler := NewAdminHandler(deps.Services.Admin)

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}

	companies := r.Group("/companies")
	{
		companies.GET("", companyHandler.List)
		companies.GET("/:id", companyHandler.Get)
		companies.POST("", middleware.AuthRequired(deps.Config), middleware.AdminRequired(), companyHandler.Create)
		companies.PATCH("/:id", middleware.AuthRequired(deps.Config), middleware.AdminRequired(), companyHandler.Update)
		companies.GET("/:id/reviews", reviewHandler.ListByCompany)
	}

	reviews := r.Group("/reviews")
	reviews.Use(middleware.AuthRequired(deps.Config))
	{
		reviews.POST("/create", reviewHandler.Create)
	}

	admin := r.Group("/admin")
	admin.Use(middleware.AuthRequired(deps.Config), middleware.AdminRequired())
	{
		admin.GET("/dashboard/insights", adminHandler.Insights)
		admin.GET("/reviews/suspicious", adminHandler.Suspicious)
		admin.POST("/reviews/:id/respond", adminHandler.Respond)
	}

	// Swagger placeholder - requires docs generation (swag init)
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
