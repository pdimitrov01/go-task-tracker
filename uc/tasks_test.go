package uc

import (
	"api/domain"
	mock "api/mocks/mock_domain"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func getTask() domain.Task {
	return domain.Task{
		ID:          uuid.MustParse("1461ec84-ccff-4f3c-af34-65d0856ac3ce"),
		Title:       "Do unit tests",
		Description: "Create extensive unit tests for all layers",
		Status:      "PENDING",
		DueDate:     time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
	}
}

func getTasksList() []domain.Task {
	return []domain.Task{
		{
			ID:          uuid.MustParse("1461ec84-ccff-4f3c-af34-65d0856ac3ce"),
			Title:       "Do unit tests",
			Description: "Create extensive unit tests for all layers",
			Status:      "PENDING",
			DueDate:     time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
			CreatedAt:   time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:          uuid.MustParse("1461ec84-ccff-4f3c-af34-65d0856ac3dd"),
			Title:       "Do other tests",
			Description: "Create tests for all layers",
			Status:      "PENDING",
			DueDate:     time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
			CreatedAt:   time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
		},
	}
}

func TestGetTaskById(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		repoMock       func(repoMock mock.MockTasksRepo)
		expectedResult domain.Task
		checks         func(t *testing.T, expected, result domain.Task, err error)
	}{
		{
			name: "happy path - OK",
			id:   "1461ec84-ccff-4f3c-af34-65d0856ac3ce",
			repoMock: func(repoMock mock.MockTasksRepo) {
				repoMock.EXPECT().GetTaskById(gomock.Any(), gomock.Eq(uuid.MustParse("1461ec84-ccff-4f3c-af34-65d0856ac3ce"))).Return(getTask(), nil)
			},
			expectedResult: getTask(),
			checks: func(t *testing.T, expected, result domain.Task, err error) {
				require.NoError(t, err)
				require.Equal(t, expected.ID, result.ID)
			},
		},
		{
			name: "no task found",
			id:   "1461ec84-ccff-4f3c-af34-65d0856ac3ce",
			repoMock: func(repoMock mock.MockTasksRepo) {
				repoMock.EXPECT().GetTaskById(gomock.Any(), gomock.Eq(uuid.MustParse("1461ec84-ccff-4f3c-af34-65d0856ac3ce"))).Return(domain.Task{}, domain.ErrTaskNotFound)
			},
			checks: func(t *testing.T, expected, result domain.Task, err error) {
				require.EqualError(t, err, domain.ErrTaskNotFound.Error())
			},
		},
		{
			name: "error",
			id:   "1461ec84-ccff-4f3c-af34-65d0856ac3ce",
			repoMock: func(repoMock mock.MockTasksRepo) {
				repoMock.EXPECT().GetTaskById(gomock.Any(), gomock.Eq(uuid.MustParse("1461ec84-ccff-4f3c-af34-65d0856ac3ce"))).Return(domain.Task{}, errors.New("not found"))
			},
			checks: func(t *testing.T, expected, result domain.Task, err error) {
				require.EqualError(t, err, "error fetching task: not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock.NewMockTasksRepo(ctrl)
			service := NewTasksService(repo)

			if tt.repoMock != nil {
				tt.repoMock(*repo)
			}

			result, err := service.GetTaskById(context.Background(), uuid.MustParse(tt.id))
			tt.checks(t, tt.expectedResult, result, err)
		})
	}
}

func TestGetTasks(t *testing.T) {
	tests := []struct {
		name           string
		repoMock       func(repoMock mock.MockTasksRepo)
		expectedResult []domain.Task
		checks         func(t *testing.T, expected, result []domain.Task, err error)
	}{
		{
			name: "happy path - OK",
			repoMock: func(repoMock mock.MockTasksRepo) {
				repoMock.EXPECT().GetTasks(gomock.Any()).Return(getTasksList(), nil)
			},
			expectedResult: getTasksList(),
			checks: func(t *testing.T, expected, result []domain.Task, err error) {
				require.NoError(t, err)
				require.Equal(t, 2, len(result))
			},
		},
		{
			name: "error",
			repoMock: func(repoMock mock.MockTasksRepo) {
				repoMock.EXPECT().GetTasks(gomock.Any()).Return(nil, errors.New("not found"))
			},
			checks: func(t *testing.T, expected, result []domain.Task, err error) {
				require.EqualError(t, err, "error fetching task: not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock.NewMockTasksRepo(ctrl)
			service := NewTasksService(repo)

			if tt.repoMock != nil {
				tt.repoMock(*repo)
			}

			result, err := service.GetTasks(context.Background())
			tt.checks(t, tt.expectedResult, result, err)
		})
	}
}

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name           string
		data           domain.Task
		repoMock       func(repoMock mock.MockTasksRepo)
		expectedResult domain.Task
		checks         func(t *testing.T, expected, result domain.Task, err error)
	}{
		{
			name: "happy path - OK",
			data: getTask(),
			repoMock: func(repoMock mock.MockTasksRepo) {
				repoMock.EXPECT().CreateTask(gomock.Any(), gomock.Eq(getTask())).Return(getTask(), nil)
			},
			expectedResult: getTask(),
			checks: func(t *testing.T, expected, result domain.Task, err error) {
				require.NoError(t, err)
				require.Equal(t, expected.ID, result.ID)
			},
		},
		{
			name: "error",
			data: getTask(),
			repoMock: func(repoMock mock.MockTasksRepo) {
				repoMock.EXPECT().CreateTask(gomock.Any(), gomock.Eq(getTask())).Return(domain.Task{}, errors.New("not found"))
			},
			checks: func(t *testing.T, expected, result domain.Task, err error) {
				require.EqualError(t, err, "error creating task: not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock.NewMockTasksRepo(ctrl)
			service := NewTasksService(repo)

			if tt.repoMock != nil {
				tt.repoMock(*repo)
			}

			result, err := service.CreateTask(context.Background(), tt.data)
			tt.checks(t, tt.expectedResult, result, err)
		})
	}
}
