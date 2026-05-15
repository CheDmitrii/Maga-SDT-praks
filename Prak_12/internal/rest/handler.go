package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"Prak_12/internal/store"
)

// Handler handles REST requests for /v1/tasks.
type Handler struct {
	store *store.Store
}

func NewHandler(s *store.Store) *Handler {
	return &Handler{store: s}
}

// errorResponse is the standard error body.
type errorResponse struct {
	Error string `json:"error"`
}

// createTaskRequest is the POST /v1/tasks body.
type createTaskRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

// updateTaskRequest is the PATCH /v1/tasks/{id} body.
type updateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

// GetTasks godoc
//
//	@Summary		Получить список задач
//	@Description	Возвращает все задачи
//	@Tags			tasks
//	@Produce		json
//	@Success		200	{array}		store.Task
//	@Router			/v1/tasks [get]
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.store.All())
}

// GetTaskByID godoc
//
//	@Summary		Получить задачу по ID
//	@Description	Возвращает одну задачу по её идентификатору
//	@Tags			tasks
//	@Produce		json
//	@Param			id	path		string	true	"ID задачи"
//	@Success		200	{object}	store.Task
//	@Failure		404	{object}	errorResponse
//	@Router			/v1/tasks/{id} [get]
func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
	t, err := h.store.GetByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	writeJSON(w, http.StatusOK, t)
}

// CreateTask godoc
//
//	@Summary		Создать задачу
//	@Description	Создаёт новую задачу
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			task	body		createTaskRequest	true	"Данные задачи"
//	@Success		201		{object}	store.Task
//	@Failure		400		{object}	errorResponse
//	@Router			/v1/tasks [post]
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	t := h.store.Create(req.Title, req.Description)
	writeJSON(w, http.StatusCreated, t)
}

// UpdateTask godoc
//
//	@Summary		Обновить задачу
//	@Description	Частичное обновление задачи по ID
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"ID задачи"
//	@Param			task	body		updateTaskRequest	true	"Обновляемые поля"
//	@Success		200		{object}	store.Task
//	@Failure		400		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Router			/v1/tasks/{id} [patch]
func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	t, err := h.store.Update(id, store.UpdateInput{
		Title:       req.Title,
		Description: req.Description,
		Done:        req.Done,
	})
	if errors.Is(err, store.ErrTaskNotFound) {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	writeJSON(w, http.StatusOK, t)
}

// DeleteTask godoc
//
//	@Summary		Удалить задачу
//	@Description	Удаляет задачу по ID
//	@Tags			tasks
//	@Param			id	path	string	true	"ID задачи"
//	@Success		204
//	@Failure		404	{object}	errorResponse
//	@Router			/v1/tasks/{id} [delete]
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
	if ok := h.store.Delete(id); !ok {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ServeHTTP dispatches /v1/tasks and /v1/tasks/{id} by method.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	isCollection := path == "/v1/tasks"

	if isCollection {
		switch r.Method {
		case http.MethodGet:
			h.GetTasks(w, r)
		case http.MethodPost:
			h.CreateTask(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	// /v1/tasks/{id}
	switch r.Method {
	case http.MethodGet:
		h.GetTaskByID(w, r)
	case http.MethodPatch:
		h.UpdateTask(w, r)
	case http.MethodDelete:
		h.DeleteTask(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
