package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/FeelsCoderMan/task-manager/internal/api/httpError"
)

type handler struct {
	service Service
}

type httpResponseSuccess struct {
	Success bool  `json:"success"`
	Task    *Task `json:"task,omitempty"`
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

func setJsonHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func jsonError(w http.ResponseWriter, error string, code int) {
	setJsonHeaders(w)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(newHttpResponseFailed(error))
}

func jsonSuccess(w http.ResponseWriter, hr *httpResponseSuccess) {
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
		case "id":
			id, err := strconv.Atoi(v[0])

			if err != nil {
				return fmt.Errorf("task-decodeFormToTask(): Could convert id %s to int", v[0])
			}

			task.ID = id
		case "name":
			task.Name = v[0]
		}
	}
	return nil
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
	r.HandleFunc("GET /tasks/{id}", customHandler(handler.get))
	r.HandleFunc("POST /tasks", customHandler(handler.post))
	r.HandleFunc("PUT /tasks/{id}", customHandler(handler.put))
	r.HandleFunc("DELETE /tasks/{id}", customHandler(handler.delete))
}

func (h *handler) post(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()

	if err != nil {
		return httpError.BadRequest("Unable to parse request body. Please check your body.")
	}

	var task Task
	err = decodeFormToTask(r.Form, &task)

	if err != nil {
		return httpError.BadRequest("Body contains invalid attributes for the task object.")
	}

	err = h.service.Create(task)

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
