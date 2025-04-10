package handler

import (
	"api/domain"
	"api/uc"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"io"
	"net/http"
	"reflect"
)

type TasksHandler struct {
	tasksService uc.TasksUC
}

func NewTasksHandler(tasksService uc.TasksUC) *TasksHandler {
	return &TasksHandler{tasksService: tasksService}
}

func (th TasksHandler) GetTaskById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getContextFromRequest(r)
	defer cancel()

	taskId := chi.URLParam(r, "id")
	id, err := uuid.Parse(taskId)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	task, err := th.tasksService.GetTaskById(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{
				Code:    404,
				Message: "task not found",
			})
			return
		}

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, task)
}

func (th TasksHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getContextFromRequest(r)
	defer cancel()

	tasksList, err := th.tasksService.GetTasks(ctx)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, tasksList)
}

func (th TasksHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := getContextFromRequest(r)
	defer cancel()

	data, err := taskFromBody(r.Body)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	task, err := th.tasksService.CreateTask(ctx, *data)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("error creating new task: %v", err),
		})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, task)
}

func taskFromBody(in io.ReadCloser) (*domain.Task, error) {
	var payload domain.Task
	decoder := json.NewDecoder(in)
	if err := decoder.Decode(&payload); err != nil {
		return nil, err
	}
	if reflect.DeepEqual(payload, domain.Task{}) {
		return nil, fmt.Errorf("invalid reuqest body")
	}
	return &payload, nil
}
