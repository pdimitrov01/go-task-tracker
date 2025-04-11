package repo

import (
	"api/config"
	"context"
	"database/sql"
	"errors"
	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
)

func NewPostgresClient(ctx context.Context, conf config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", conf.DbConnectionUrl)
	if err != nil {
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(migrationsPath string, conf config.Config) error {
	m, err := migrate.New(migrationsPath, conf.DbConnectionUrl)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("no new migrations")
		return nil
	}
	log.Println("migrations run successfully")
	return nil
}
