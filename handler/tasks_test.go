package handler

import (
	"api/domain"
	mock "api/mocks/mock_uc"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func getExpectedBody() domain.Task {
	return domain.Task{
		ID:          uuid.MustParse("1461ec84-ccff-4f3c-af34-65d0856ac3ce"),
		Title:       "Do unit tests",
		Description: "Create extensive unit tests for all layers",
		Status:      "PENDING",
		DueDate:     time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
		CreatedAt:   time.Date(2025, 4, 3, 0, 0, 0, 0, time.UTC),
	}
}

func TestGetTaskById(t *testing.T) {
	tests := []struct {
		name               string
		id                 string
		ucMock             func(ucMock mock.MockTasksUC)
		expectedStatusCode int
		expectedBody       domain.Task
	}{
		{
			name: "happy path - OK",
			id:   "1461ec84-ccff-4f3c-af34-65d0856ac3ce",
			ucMock: func(ucMock mock.MockTasksUC) {
				ucMock.EXPECT().GetTaskById(gomock.Any(), gomock.Eq(uuid.MustParse("1461ec84-ccff-4f3c-af34-65d0856ac3ce"))).Return(getExpectedBody(), nil)
			},
			expectedStatusCode: 200,
			expectedBody:       getExpectedBody(),
		},
		{
			name:               "wrong id type",
			id:                 "invalid id",
			expectedStatusCode: 400,
			expectedBody:       domain.Task{},
		},
		{
			name: "no task found",
			id:   "1461ec84-ccff-4f3c-af34-65d0856ac3ce",
			ucMock: func(ucMock mock.MockTasksUC) {
				ucMock.EXPECT().GetTaskById(gomock.Any(), gomock.Eq(uuid.MustParse("1461ec84-ccff-4f3c-af34-65d0856ac3ce"))).Return(domain.Task{}, domain.ErrTaskNotFound)
			},
			expectedStatusCode: 404,
			expectedBody:       domain.Task{},
		},
		{
			name: "internal server error",
			id:   "1461ec84-ccff-4f3c-af34-65d0856ac3ce",
			ucMock: func(ucMock mock.MockTasksUC) {
				ucMock.EXPECT().GetTaskById(gomock.Any(), gomock.Eq(uuid.MustParse("1461ec84-ccff-4f3c-af34-65d0856ac3ce"))).Return(domain.Task{}, errors.New("not found"))
			},
			expectedStatusCode: 500,
			expectedBody:       domain.Task{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := chi.NewRouter()
			recorder := httptest.NewRecorder()
			ucMock := mock.NewMockTasksUC(ctrl)
			handler := NewTasksHandler(ucMock)

			if tt.ucMock != nil {
				tt.ucMock(*ucMock)
			}

			r.Get("/api/task/{id}", handler.GetTaskById)
			req, err := http.NewRequest(http.MethodGet, "/api/task/"+tt.id, nil)
			require.NoError(t, err)

			r.ServeHTTP(recorder, req)
			require.Equal(t, tt.expectedStatusCode, recorder.Code)

			var body domain.Task
			err = json.Unmarshal(recorder.Body.Bytes(), &body)
			require.NoError(t, err)
			require.Equal(t, tt.expectedBody, body)
		})
	}
}

func TestGetTasks(t *testing.T) {
	tests := []struct {
		name               string
		ucMock             func(ucMock mock.MockTasksUC)
		expectedStatusCode int
		expectedCount      int
	}{
		{
			name: "happy path - OK",
			ucMock: func(ucMock mock.MockTasksUC) {
				ucMock.EXPECT().GetTasks(gomock.Any()).Return([]domain.Task{getExpectedBody()}, nil)
			},
			expectedStatusCode: 200,
			expectedCount:      1,
		},
		{
			name: "no tasks found",
			ucMock: func(ucMock mock.MockTasksUC) {
				ucMock.EXPECT().GetTasks(gomock.Any()).Return([]domain.Task{}, nil)
			},
			expectedStatusCode: 200,
			expectedCount:      0,
		},
		{
			name: "internal server error",
			ucMock: func(ucMock mock.MockTasksUC) {
				ucMock.EXPECT().GetTasks(gomock.Any()).Return(nil, errors.New("error occurred"))
			},
			expectedStatusCode: 500,
			expectedCount:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := chi.NewRouter()
			recorder := httptest.NewRecorder()
			ucMock := mock.NewMockTasksUC(ctrl)
			handler := NewTasksHandler(ucMock)

			if tt.ucMock != nil {
				tt.ucMock(*ucMock)
			}

			r.Get("/api/tasks", handler.GetTasks)
			req, err := http.NewRequest(http.MethodGet, "/api/tasks", nil)
			require.NoError(t, err)

			r.ServeHTTP(recorder, req)
			require.Equal(t, tt.expectedStatusCode, recorder.Code)

			if tt.expectedStatusCode == http.StatusOK {
				var body []domain.Task
				err = json.Unmarshal(recorder.Body.Bytes(), &body)
				require.NoError(t, err)
				require.Equal(t, tt.expectedCount, len(body))
			} else {
				var errResp ErrorResponse
				err = json.Unmarshal(recorder.Body.Bytes(), &errResp)
				require.NoError(t, err)
				require.Equal(t, tt.expectedStatusCode, errResp.Code)
			}
		})
	}
}

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name               string
		body               string
		ucMock             func(ucMock mock.MockTasksUC)
		expectedStatusCode int
		expectedTask       string
	}{
		{
			name: "happy path - OK",
			body: `{"id": "1461ec84-ccff-4f3c-af34-65d0856ac3ce", "title": "Do unit tests", "description": "Create extensive unit tests for all layers", "status": "PENDING", "due_date": "2025-05-12T00:00:00Z"}`,
			ucMock: func(ucMock mock.MockTasksUC) {
				ucMock.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(getExpectedBody(), nil)
			},
			expectedStatusCode: 200,
			expectedTask:       "1461ec84-ccff-4f3c-af34-65d0856ac3ce",
		},
		{
			name:               "bad request",
			body:               "invalid body",
			expectedStatusCode: 400,
		},
		{
			name: "internal server error",
			body: `{"id": "1461ec84-ccff-4f3c-af34-65d0856ac3ce", "title": "Do unit tests", "description": "Create extensive unit tests for all layers", "status": "PENDING", "due_date": "2025-05-12T00:00:00Z"}`,
			ucMock: func(ucMock mock.MockTasksUC) {
				ucMock.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(domain.Task{}, errors.New("not found"))
			},
			expectedStatusCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := chi.NewRouter()
			recorder := httptest.NewRecorder()
			ucMock := mock.NewMockTasksUC(ctrl)
			handler := NewTasksHandler(ucMock)

			if tt.ucMock != nil {
				tt.ucMock(*ucMock)
			}

			body := []byte(tt.body)

			r.Post("/api/task", handler.CreateTask)
			req, err := http.NewRequest(http.MethodPost, "/api/task", bytes.NewReader(body))
			require.NoError(t, err)

			r.ServeHTTP(recorder, req)
			require.Equal(t, tt.expectedStatusCode, recorder.Code)

			if tt.expectedStatusCode == http.StatusOK {
				var body domain.Task
				err = json.Unmarshal(recorder.Body.Bytes(), &body)
				require.NoError(t, err)
				require.Equal(t, tt.expectedTask, body.ID.String())
			} else {
				var errResp ErrorResponse
				err = json.Unmarshal(recorder.Body.Bytes(), &errResp)
				require.NoError(t, err)
				require.Equal(t, tt.expectedStatusCode, errResp.Code)
			}
		})
	}
}
