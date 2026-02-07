package database

import (
	"context"

	"github.com/D3rise/dchat/internal/interfaces"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

type Database struct {
	Gorm *gorm.DB
}

func NewDatabaseConn(lc fx.Lifecycle, envService interfaces.EnvService, log *zap.Logger) Database {
	logger := zapgorm2.New(log)

	dsn := envService.GetDatabaseDSN()
	db := Database{}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			gormDb, err := gorm.Open(
				postgres.Open(dsn),
				&gorm.Config{Logger: logger},
			)

			if err != nil {
				return err
			}

			db.Gorm = gormDb
			return nil
		},
	})

	return db
}
