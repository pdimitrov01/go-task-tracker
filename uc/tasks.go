package uc

import (
	"api/domain"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type TasksUC interface {
	GetTaskById(ctx context.Context, id uuid.UUID) (domain.Task, error)
	GetTasks(ctx context.Context) ([]domain.Task, error)
	CreateTask(ctx context.Context, data domain.Task) (domain.Task, error)
	//UpdateTask(ctx context.Context, data domain.Task) (domain.Task, error)
	//DeleteTask(ctx context.Context) error
}

type TasksService struct {
	tasksRepo domain.TasksRepo
}

func NewTasksService(tasksRepo domain.TasksRepo) *TasksService {
	return &TasksService{tasksRepo: tasksRepo}
}

func (ts TasksService) GetTaskById(ctx context.Context, id uuid.UUID) (domain.Task, error) {
	task, err := ts.tasksRepo.GetTaskById(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			return domain.Task{}, err
		}
		return domain.Task{}, fmt.Errorf("error fetching task: %v", err)
	}
	return task, nil
}

func (ts TasksService) GetTasks(ctx context.Context) ([]domain.Task, error) {
	task, err := ts.tasksRepo.GetTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching task: %v", err)
	}
	return task, nil
}

func (ts TasksService) CreateTask(ctx context.Context, data domain.Task) (domain.Task, error) {
	task, err := ts.tasksRepo.CreateTask(ctx, data)
	if err != nil {
		return domain.Task{}, fmt.Errorf("error creating task: %v", err)
	}
	return task, nil
}
