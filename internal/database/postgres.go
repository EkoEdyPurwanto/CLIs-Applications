package database

import (
	"MiniProjRamadh/internal/models"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

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
