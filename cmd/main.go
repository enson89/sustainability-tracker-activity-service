package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/enson89/sustainability-tracker-activity-service/internal/cache"
	"github.com/enson89/sustainability-tracker-activity-service/internal/config"
	"github.com/enson89/sustainability-tracker-activity-service/internal/handler"
	"github.com/enson89/sustainability-tracker-activity-service/internal/middleware"
	"github.com/enson89/sustainability-tracker-activity-service/internal/repository"
	"github.com/enson89/sustainability-tracker-activity-service/internal/service"
)

func main() {
	// Initialize zap logger.
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	// Connect to PostgreSQL.
	db, err := repository.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	activityRepo := repository.NewActivityRepository(db)

	// Initialize Redis client from the cache package.
	redisClient, err := cache.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		logger.Fatal("failed to connect to redis", zap.Error(err))
	}

	activityService := service.NewActivityService(activityRepo, redisClient)
	activityHandler := handler.NewActivityHandler(activityService)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler(logger))
	router.Use(gin.Logger())

	api := router.Group("/api")
	{
		api.POST("/activities", activityHandler.CreateActivity)
		api.GET("/activities/:id", activityHandler.GetActivity)
		api.PUT("/activities/:id", activityHandler.UpdateActivity)
		api.DELETE("/activities/:id", activityHandler.DeleteActivity)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	router.Run(":" + port)
}
