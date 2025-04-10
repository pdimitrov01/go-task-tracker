package repo

import (
	"api/adapter/repo/postgres/gen"
	"api/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

const traceNameTasksRepo = "TasksRepo"

type TasksRepo struct {
	querier gen.Querier
}

func NewTasksRepo(querier gen.Querier) *TasksRepo {
	return &TasksRepo{querier: querier}
}

func (tr TasksRepo) GetTaskById(ctx context.Context, id uuid.UUID) (domain.Task, error) {
	ctx, span := otel.GetTracerProvider().Tracer(traceNameTasksRepo).Start(ctx, traceNameTasksRepo+".GetTaskById")
	span.SetAttributes(attribute.String("task_id", id.String()))
	defer span.End()

	task, err := tr.querier.GetTaskById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Task{}, fmt.Errorf("task not found in db %s: %v", id, err)
		}
		return domain.Task{}, fmt.Errorf("failed to get task %s: %v", id, err)
	}

	return task.ToDomain(), nil
}

func (tr TasksRepo) GetTasks(ctx context.Context) ([]domain.Task, error) {
	ctx, span := otel.GetTracerProvider().Tracer(traceNameTasksRepo).Start(ctx, traceNameTasksRepo+".GetTasks")
	defer span.End()

	data, err := tr.querier.GetTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find all tasks: %v", err)
	}

	var tasks []domain.Task
	for _, currentTask := range data {
		tasks = append(tasks, currentTask.ToDomain())
	}

	return tasks, nil
}

func (tr TasksRepo) CreateTask(ctx context.Context, data domain.Task) (domain.Task, error) {
	ctx, span := otel.GetTracerProvider().Tracer(traceNameTasksRepo).Start(ctx, traceNameTasksRepo+".GetTasks")
	defer span.End()

	task, err := tr.querier.SaveTask(ctx, gen.SaveTaskParams{
		ID:          data.ID,
		Title:       data.Title,
		Description: data.Description,
		Status:      data.Status,
		DueDate:     data.DueDate,
	})
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to save task: %v", err)
	}

	return task.ToDomain(), nil
}
