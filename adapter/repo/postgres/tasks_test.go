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
