package domain

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

var ErrTaskNotFound = errors.New("task not found")

type TasksRepo interface {
	GetTaskById(ctx context.Context, id uuid.UUID) (Task, error)
	GetTasks(ctx context.Context) ([]Task, error)
	CreateTask(ctx context.Context, data Task) (Task, error)
	//UpdateTask(ctx context.Context, data Task) (Task, error)
	//DeleteTask(ctx context.Context, id string) error
}

type Task struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DueDate     time.Time `json:"due_date"`
	CreatedAt   time.Time `json:"created_at"`
}
