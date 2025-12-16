package postgres

import (
	"container-manager/pkg/config"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewClient creates a new database connection pool.
func NewClient(cfg config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
