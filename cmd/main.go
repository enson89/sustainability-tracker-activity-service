package main

import (
	"log"
	"os"

	"github.com/enson89/sustainability-tracker-activity-service/internal/cache"
	"github.com/enson89/sustainability-tracker-activity-service/internal/config"
	"github.com/enson89/sustainability-tracker-activity-service/internal/handler"
	"github.com/enson89/sustainability-tracker-activity-service/internal/repository"
	"github.com/enson89/sustainability-tracker-activity-service/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to PostgreSQL
	db, err := repository.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	activityRepo := repository.NewActivityRepository(db)

	// Initialize Redis client from the cache package
	redisClient, err := cache.NewRedisClient(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	activityService := service.NewActivityService(activityRepo, redisClient)
	activityHandler := handler.NewActivityHandler(activityService)

	router := gin.Default()
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
