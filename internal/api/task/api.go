package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/FeelsCoderMan/task-manager/internal/api/httpError"
)

type handler struct {
	service Service
}

type HttpResponseSuccess interface {
	IsSuccess() bool
}

func (hr *httpResponseSuccess) IsSuccess() bool {
	return hr.Success
}

func (hr *httpResponseSuccessMultiple) IsSuccess() bool {
	return hr.Success
}

type httpResponseSuccess struct {
	Success bool  `json:"success"`
	Task    *Task `json:"task,omitempty"`
}

type httpResponseSuccessMultiple struct {
	Success bool   `json:"success"`
	Tasks   []Task `json:"tasks,omitempty"`
}

type httpResponseFailed struct {
	Error        bool   `json:"error"`
	ErrorMessage string `json:"errorMessage"`
}

func newHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func newHttpResponseFailed(errorMessage string) *httpResponseFailed {
	return &httpResponseFailed{
		Error:        true,
		ErrorMessage: errorMessage,
	}
}

func newHttpResponseSuccess(task *Task) *httpResponseSuccess {
	return &httpResponseSuccess{
		Success: true,
		Task:    task,
	}
}

func newHttpResponseSuccessMultiple(tasks []Task) *httpResponseSuccessMultiple {
	return &httpResponseSuccessMultiple{
		Success: true,
		Tasks:   tasks,
	}
}

func setJsonHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func jsonError(w http.ResponseWriter, error string, code int) {
	setJsonHeaders(w)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(newHttpResponseFailed(error))
}

func jsonSuccess(w http.ResponseWriter, hr HttpResponseSuccess) {
	setJsonHeaders(w)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(hr)
}

func customHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)

		if err != nil {
			var httpError *httpError.HttpError
			if errors.As(err, &httpError) {
				jsonError(w, err.Error(), httpError.Code)
			} else {
				jsonError(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func decodeFormToTask(form url.Values, task *Task) error {
	for k, v := range form {
		switch k {
		case "name":
			task.Name = v[0]
		case "description":
			task.Description = v[0]
		case "priority":
			priority, err := strconv.Atoi(v[0])

			if err != nil {
				return fmt.Errorf("task-decodeFormToTask(): Could convert id %s to int", v[0])
			}

			task.Priority = uint8(priority)
		case "tags":
			task.Tags = []string{v[0]}
		}
	}
	return nil
}

func updateDate(task *Task) {
	task.UpdatedAt = time.Now()
	task.CreatedAt = time.Now()
}

func getTaskIdFromPath(r *http.Request) (int, error) {
	id := r.PathValue("id")

	if id == "" {
		return -1, fmt.Errorf("id parameter is missing")
	}

	convertedId, err := strconv.Atoi(id)

	if err != nil {
		return -1, fmt.Errorf("id parameter %s is not a valid attribute", id)
	}

	return convertedId, err
}

func RegisterHandlers(r *http.ServeMux, service Service) {
	handler := newHandler(service)
	r.HandleFunc("GET /task/latest/", customHandler(handler.getLatest))
	r.HandleFunc("GET /task/{id}", customHandler(handler.get))
	r.HandleFunc("POST /task/", customHandler(handler.post))
	r.HandleFunc("PUT /task/{id}", customHandler(handler.put))
	r.HandleFunc("DELETE /task/{id}", customHandler(handler.delete))
}

func (h *handler) post(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseMultipartForm(1024)

	if err != nil {
		return httpError.BadRequest("Unable to parse request body. Please check your body.")
	}

	var task Task
	err = decodeFormToTask(r.Form, &task)

	if err != nil {
		return httpError.BadRequest("Body contains invalid attributes for the task object.")
	}

	updateDate(&task)
	task, err = h.service.Create(task)

	if err != nil {
		return httpError.InternalServerError("An unexpected error occurred while creating the task")
	}

	jsonSuccess(w, newHttpResponseSuccess(&task))
	return nil
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) error {
	id, err := getTaskIdFromPath(r)

	if err != nil {
		return httpError.BadRequest("id parameter is missing or not valid")
	}

	task, err := h.service.Get(id)

	if err != nil {
		return httpError.InternalServerError("An unexpected error occurred while retrieving the task")
	}

	jsonSuccess(w, newHttpResponseSuccess(&task))
	return nil
}

func (h *handler) getLatest(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()

	countStr := query.Get("count")

	if countStr == "" {
		return httpError.BadRequest("count parameter is missing")
	}

	count, err := strconv.Atoi(countStr)

	if err != nil {
		return httpError.BadRequest("count parameter is not valid")
	}

	tasks, err := h.service.GetLatest(count)

	if err != nil {
		return httpError.InternalServerError("An expected error occurred while retriving latest tasks")
	}

	jsonSuccess(w, newHttpResponseSuccessMultiple(tasks))
	return nil
}

func (h *handler) put(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()

	if err != nil {
		return fmt.Errorf("Could not read request body: %w", err)
	}

	var task Task
	err = decodeFormToTask(r.Form, &task)

	if err != nil {
		return httpError.BadRequest("Body contains invalid attributes for the task object.")
	}

	err = h.service.Update(task)

	if err != nil {
		return httpError.InternalServerError("An unexpected error occurred while updating the task")
	}

	jsonSuccess(w, newHttpResponseSuccess(&task))
	return nil
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) error {
	id, err := getTaskIdFromPath(r)

	if err != nil {
		return httpError.BadRequest("id parameter is missing or not valid")
	}

	err = h.service.Delete(id)

	if err != nil {
		return httpError.InternalServerError("An unexpected error occurred while deleting the task")
	}

	var task Task
	jsonSuccess(w, newHttpResponseSuccess(&task))
	return nil
}
