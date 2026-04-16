package database

import (
	"log"

	"github.com/xinewang/oen/internal/config"
	"github.com/xinewang/oen/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.DatabaseConfig) *gorm.DB {
	dsn := cfg.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	if err := AutoMigrate(db); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}

	log.Println("Database connected and migrated successfully")
	return db
}

// AutoMigrate runs GORM auto migration for all models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Agent{},
		&model.AgentHeartbeat{},
		&model.Artifact{},
		&model.ArtifactVersion{},
		&model.ArtifactView{},
		&model.CandidateResource{},
		&model.Recommendation{},
		&model.RecommendationDecision{},
		&model.ConsentRecord{},
		&model.AuditLog{},
	)
}
