package main

import (
	"log"
	
	"context"

	"crowdreview/config"
	"crowdreview/internal/handlers"
	"crowdreview/internal/models"
	"crowdreview/internal/repository"
	"crowdreview/internal/services"
	"crowdreview/internal/validation"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()

	db, err := connectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	rdb, err := connectRedis(cfg.RedisURL)
	if err != nil {
		log.Printf("warning: redis unavailable (%v), rate limiting disabled", err)
	}

	repos := repository.NewRepositories(db)
	engine := validation.NewFraudEngine()
	worker := validation.NewFraudWorker(engine, repos.Validation)
	worker.Start()

	svc := services.NewServices(cfg, repos, rdb, worker)
	router := handlers.SetupRouter(handlers.RouterDeps{
		Config:   cfg,
		Services: svc,
		Redis:    rdb,
	})

	log.Printf("CrowdReview API listening on :%s", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}

func connectDB(url string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Printf("could not ensure uuid extension: %v", err)
	}
	// Auto-migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.AdminUser{},
		&models.Company{},
		&models.Review{},
		&models.ReviewValidationResult{},
		&models.FraudSignal{},
		&models.Achievement{},
		&models.UserAchievement{},
	); err != nil {
		return nil, err
	}
	return db, nil
}

func connectRedis(url string) (*redis.Client, error) {
    opts, err := redis.ParseURL(url)
    if err != nil {
        return nil, err
    }
    client := redis.NewClient(opts)
    ctx := context.Background()
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, err
    }
    return client, nil
}
