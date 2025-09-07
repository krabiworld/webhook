package db

import (
	"context"
	"errors"
	"webhook/internal/config"
	"webhook/internal/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Init() {
	var err error

	dialector, err := openDialector()
	if err != nil {
		log.Panic().Err(err).Msg("Failed to open dialector")
	}

	db, err = gorm.Open(dialector, &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect db")
	}

	log.Info().Msg("Database initialized")

	err = db.AutoMigrate(&models.Token{}, &models.Webhook{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate db")
	}

	ctx := context.Background()

	// create bootstrap token
	_, err = G[models.Token]().Where("id = ?", models.BootstrapToken).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		t, err := uuid.NewV7()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to generate uuid")
		}

		err = G[models.Token]().Create(ctx, &models.Token{
			ID:    models.BootstrapToken,
			Token: t.String(),
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create token")
		}

		log.Info().Str("token", t.String()).Msg("Bootstrap token created")
	} else if err != nil {
		log.Fatal().Err(err).Msg("Failed to get token")
	}
}

func G[T any]() gorm.Interface[T] {
	return gorm.G[T](db)
}

func openDialector() (gorm.Dialector, error) {
	switch config.Get().DatabaseType {
	case "postgres":
		return postgres.Open(config.Get().DatabaseUrl), nil
	case "sqlite":
		return sqlite.Open(config.Get().DatabaseUrl), nil
	default:
		return nil, errors.New("unsupported dialector type")
	}
}
