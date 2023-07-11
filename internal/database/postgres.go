package database

import (
	"MiniProjRamadh/internal/models"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"sync"
)

func GetDBInstance(cfg *models.Config) (*sql.DB, error) {
	var err error
	var dbInstance *sql.DB
	var onceInstance sync.Once

	onceInstance.Do(func() {
		dbInstance, err = ConnectDB(cfg)
	})
	return dbInstance, err
}

func ConnectDB(cfg *models.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
	)
	db, err := sql.Open(cfg.DBDriver, dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Migrate migrates the database schema.
func Migrate(db *sql.DB) error {
	log.Info("running database migration")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/database/migrations",
		"mydatabase", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	log.Info("database migration completed successfully")
	return nil
}
