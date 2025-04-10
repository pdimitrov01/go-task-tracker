package main

import (
	"api/adapter/repo/postgres"
	"api/config"
	"context"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conf, err := config.FromEnv()
	if err != nil {
		log.Fatalf("failed to load env vars to config: %v", err)
	}

	_, err = repo.NewPostgresClient(ctx, *conf)
	if err != nil {
		log.Fatalf("error while creating postgres connection: %v", err)
	}

	if err := repo.RunMigrations("file://adapter/repo/postgres/migrations", *conf); err != nil {
		log.Fatalf("error while migrating postgres scripts: %v", err)
	}

}
