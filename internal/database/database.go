package database

import (
	"flash-sale-coupon-system/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Initialize creates a new database connection
func Initialize(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Migrate runs database migrations
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Coupon{},
		&models.Claim{},
	)
}
