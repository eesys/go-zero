package gormx

import (
    "time"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "gorm.io/gorm/schema"
)

type Config struct {
    Gorm *gorm.Config
    MaxOpenConns int
    MaxIdleConns int
    ConnMaxLifetime time.Duration
}

func defaultConfig() *gorm.Config {
    return &gorm.Config{
        NamingStrategy: schema.NamingStrategy{
            SingularTable: true,
        },
        Logger: logger.Default.LogMode(logger.Warn),
        DisableNestedTransaction: true,
    }
}

