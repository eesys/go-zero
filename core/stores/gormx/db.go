package gormx

import (
    "database/sql"
    "fmt"
    "net/url"
    "gorm.io/driver/mysql"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Dialect string

const (
    MySQL    Dialect = "mysql"
    Postgres Dialect = "postgres"
)

type DSN struct {
    Driver Dialect
    DSN    string
}

func Open(d DSN, cfg *Config) (*gorm.DB, *sql.DB, error) {
    if cfg == nil || cfg.Gorm == nil {
        cfg = &Config{Gorm: defaultConfig()}
    }
    var g *gorm.DB
    var err error
    switch d.Driver {
    case MySQL:
        g, err = gorm.Open(mysql.Open(d.DSN), cfg.Gorm)
    case Postgres:
        g, err = gorm.Open(postgres.Open(d.DSN), cfg.Gorm)
    default:
        err = fmt.Errorf("unsupported driver: %s", d.Driver)
    }
    if err != nil {
        return nil, nil, err
    }
    sqlDB, err := g.DB()
    if err != nil {
        return nil, nil, err
    }
    if cfg.MaxOpenConns > 0 {
        sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
    }
    if cfg.MaxIdleConns > 0 {
        sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
    }
    // caller may configure ConnMaxLifetime via sql.DB returned
    return g, sqlDB, nil
}

func MaskDSN(dsn string) string {
    u, err := url.Parse(dsn)
    if err != nil {
        return ""
    }
    if u.User != nil {
        name := u.User.Username()
        u.User = url.UserPassword(name, "***")
    }
    return u.String()
}
