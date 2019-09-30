package postgres

import (
	"context"
	"database/sql"
	"github.com/ProtocolONE/go-core/config"
	"github.com/ProtocolONE/go-core/invoker"
	"github.com/ProtocolONE/go-core/logger"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	// Attach GORM postgres adapter
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// Attach PQ adapter
	_ "github.com/lib/pq"
	mocket "github.com/selvatico/go-mocket"
)

// ProviderCfg returns configuration for production GORM
func ProviderCfg(cfg config.Configurator) (*Config, func(), error) {
	c := &Config{
		invoker: invoker.NewInvoker(),
	}
	e := cfg.UnmarshalKeyOnReload(UnmarshalKey, c)
	return c, func() {}, e
}

// ProviderGORM returns GORM instance with resolved dependencies
func ProviderGORM(ctx context.Context, log logger.Logger, cfg *Config) (*gorm.DB, func(), error) {
	log = log.WithFields(logger.Fields{"service": Prefix})
	db, err := gorm.Open("postgres", cfg.Dsn)
	if err != nil {
		return nil, nil, err
	}
	db.DB().SetMaxOpenConns(cfg.MaxOpenConns)
	db.DB().SetMaxIdleConns(cfg.MaxIdleConns)
	db.DB().SetConnMaxLifetime(cfg.ConnMaxLifetime)
	if cfg.Debug {
		db.LogMode(true)
	}
	db.SetLogger(NewLoggerAdapter(log, logger.LevelDebug))
	cleanup := func() {
		_ = db.Close()
	}
	return db, cleanup, nil
}

// ProviderGORMTest returns stub/mock GORM instance with resolved dependencies
func ProviderGORMTest() (*gorm.DB, func(), error) {
	var db *gorm.DB
	mocket.Catcher.Register()
	sqlDB, err := sql.Open(mocket.DriverName, "gorm")
	if err != nil {
		return db, nil, err
	}
	db, err = gorm.Open("postgres", sqlDB)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		_ = db.Close()
	}
	db.LogMode(false)
	return db, cleanup, nil
}

var (
	WireSet     = wire.NewSet(ProviderGORM, ProviderCfg)
	WireTestSet = wire.NewSet(ProviderGORMTest)
)
