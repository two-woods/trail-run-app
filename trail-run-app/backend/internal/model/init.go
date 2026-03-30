package model

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully")
	return nil
}

func AutoMigrate() error {
	err := DB.AutoMigrate(
		&User{},
		&Race{},
		&AidStation{},
		&Equipment{},
		&Result{},
		&EquipmentCheck{},
		&FavoriteRace{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	// Create indexes
	if err := CreateIndexes(); err != nil {
		return err
	}

	log.Println("Database migration completed")
	return nil
}

func CreateIndexes() error {
	// High frequency query indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_race_date ON races(date)",
		"CREATE INDEX IF NOT EXISTS idx_race_location ON races(province, city)",
		"CREATE INDEX IF NOT EXISTS idx_result_user ON results(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_result_race ON results(race_id)",
		"CREATE INDEX IF NOT EXISTS idx_favorite_user ON favorite_races(user_id)",
		"CREATE INDEX IF NOT EXISTS idx_aid_station_race ON aid_stations(race_id)",
	}

	for _, idx := range indexes {
		if err := DB.Exec(idx).Error; err != nil {
			log.Printf("Warning: failed to create index: %v", err)
		}
	}

	return nil
}
