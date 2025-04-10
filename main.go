package main

import (
	"api/adapter/repo/postgres"
	"api/adapter/repo/postgres/gen"
	"api/config"
	"api/handler"
	"api/uc"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conf, err := config.FromEnv()
	if err != nil {
		log.Fatalf("failed to load env vars to config: %v", err)
	}

	db, err := repo.NewPostgresClient(ctx, *conf)
	if err != nil {
		log.Fatalf("error while creating postgres connection: %v", err)
	}

	if err := repo.RunMigrations("file://adapter/repo/postgres/migrations", *conf); err != nil {
		log.Fatalf("error while migrating postgres scripts: %v", err)
	}

	dbRepo := gen.New(db)

	r := createRouter(dbRepo)
	if err := http.ListenAndServe(":"+conf.ApiPort, r); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("error occured while listening port: %v", err)
	}
}

func createRouter(dbRepo *gen.Queries) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	tasksRepo := repo.NewTasksRepo(dbRepo)
	tasksService := uc.NewTasksService(tasksRepo)
	tasksHandler := handler.NewTasksHandler(tasksService)

	r.Group(func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Get("/task/{id}", tasksHandler.GetTaskById)
			r.Get("/tasks", tasksHandler.GetTasks)
			r.Post("/task", tasksHandler.CreateTask)
		})
	})

	return r
}
