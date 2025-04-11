package repo

import (
	"api/adapter/repo/postgres/gen"
	"api/domain"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetTaskById_Success(t *testing.T) {
	t.Parallel()
	id := uuid.MustParse("1461ec84-ccff-4f3c-af34-65d0856ac3ce")

	db := getIsolatedDatabase(t)
	assert.NotEqual(t, nil, db)

	repo := NewTasksRepo(gen.New(db))

	_, err := repo.CreateTask(context.Background(), domain.Task{
		ID:          id,
		Title:       "Do unit tests",
		Description: "Create extensive unit tests for all layers",
		Status:      "PENDING",
		DueDate:     time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	task, err := repo.GetTaskById(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, task.ID)
}

func TestGetTaskById_NotFound(t *testing.T) {
	t.Parallel()
	id := uuid.MustParse("1461ec84-ccff-4f3c-af34-65d0856ac3ce")

	db := getIsolatedDatabase(t)
	assert.NotEqual(t, nil, db)

	repo := NewTasksRepo(gen.New(db))

	_, err := repo.GetTaskById(context.Background(), id)
	require.Error(t, err, "not found")
}

func TestGetTasks_Success(t *testing.T) {
	t.Parallel()
	id1 := uuid.MustParse("1461ec84-ccff-4f3c-af34-65d0856ac3ce")
	id2 := uuid.MustParse("9d373cde-0ba2-45ba-b3dc-61bcfe2faadb")

	db := getIsolatedDatabase(t)
	assert.NotEqual(t, nil, db)

	repo := NewTasksRepo(gen.New(db))

	_, err := repo.CreateTask(context.Background(), domain.Task{
		ID:          id1,
		Title:       "Do unit tests",
		Description: "Create extensive unit tests for all layers",
		Status:      "PENDING",
		DueDate:     time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	_, err = repo.CreateTask(context.Background(), domain.Task{
		ID:          id2,
		Title:       "Do tests",
		Description: "Create tests for all layers",
		Status:      "PENDING",
		DueDate:     time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)

	tasks, err := repo.GetTasks(context.Background())
	require.NoError(t, err)
	require.Equal(t, 2, len(tasks))
}

func TestGetTasks_EmptyDB(t *testing.T) {
	t.Parallel()

	db := getIsolatedDatabase(t)
	assert.NotEqual(t, nil, db)

	repo := NewTasksRepo(gen.New(db))

	tasks, err := repo.GetTasks(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, len(tasks))
}

func TestCreateTask_Success(t *testing.T) {
	t.Parallel()
	id := uuid.MustParse("1461ec84-ccff-4f3c-af34-65d0856ac3ce")

	db := getIsolatedDatabase(t)
	assert.NotEqual(t, nil, db)

	repo := NewTasksRepo(gen.New(db))

	createdTask, err := repo.CreateTask(context.Background(), domain.Task{
		ID:          id,
		Title:       "Do unit tests",
		Description: "Create extensive unit tests for all layers",
		Status:      "PENDING",
		DueDate:     time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)
	require.Equal(t, id, createdTask.ID)

	task, err := repo.GetTaskById(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, createdTask.ID, task.ID)
}

func TestCreateTask_Conflict(t *testing.T) {
	t.Parallel()
	id := uuid.MustParse("1461ec84-ccff-4f3c-af34-65d0856ac3ce")

	db := getIsolatedDatabase(t)
	assert.NotEqual(t, nil, db)

	repo := NewTasksRepo(gen.New(db))

	createdTask, err := repo.CreateTask(context.Background(), domain.Task{
		ID:          id,
		Title:       "Do unit tests",
		Description: "Create extensive unit tests for all layers",
		Status:      "PENDING",
		DueDate:     time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
	})
	require.NoError(t, err)
	require.Equal(t, id, createdTask.ID)

	_, err = repo.CreateTask(context.Background(), domain.Task{
		ID:          id,
		Title:       "Do unit tests",
		Description: "Create extensive unit tests for all layers",
		Status:      "PENDING",
		DueDate:     time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
	})
	require.Error(t, err)
}
