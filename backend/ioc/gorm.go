package ioc

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/mymikasa/prompthub/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	gormCfg := &gorm.Config{}
	if cfg.IsDev() {
		gormCfg.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormCfg.Logger = logger.Default.LogMode(logger.Warn)
	}

	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get underlying db: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	slog.Info("mysql connected")
	return db, nil
}
