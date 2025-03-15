package db

import (
	"chach/massager/config"
	models "chach/massager/db/model"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storege struct {
	DB     *gorm.DB
	Config config.Config
}

func NewStorege(config config.Config) (*Storege, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.PgConfig.Host, config.PgConfig.User, config.PgConfig.Password, config.PgConfig.Database, config.PgConfig.Port)

	database, err := gorm.Open(postgres.Open(dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = database.AutoMigrate(&models.Message{},
		&models.User{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &Storege{
		DB:     database,
		Config: config,
	}, nil
}
