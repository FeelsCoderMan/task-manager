package task

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type handler struct {
	service Service
}

func newHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func customHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)

		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
}

func decodeFormToTask(form url.Values, task *Task) error {
	for k, v := range form {
		switch k {
		case "id":
			id, err := strconv.Atoi(v[0])

			if err != nil {
				return fmt.Errorf("task-decodeFormToTask(): Could convert id to int: %w", err)
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
	convertedId, err := strconv.Atoi(id)

	if err != nil {
		return -1, fmt.Errorf("Could not convert id to int: %w", err)
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
		return fmt.Errorf("Could not read request body: %w", err)
	}

	var task Task
	err = decodeFormToTask(r.Form, &task)

	if err != nil {
		return err
	}

	fmt.Println(task)
	err = h.service.Create(task)

	if err != nil {
		return err
	}

	return nil
}

func (h *handler) get(w http.ResponseWriter, r *http.Request) error {
	id, err := getTaskIdFromPath(r)

	if err != nil {
		return err
	}

	_, err = h.service.Get(id)

	if err != nil {
		return err
	}

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
		return err
	}

	err = h.service.Update(task)

	if err != nil {
		return err
	}

	return nil
}

func (h *handler) delete(w http.ResponseWriter, r *http.Request) error {
	id, err := getTaskIdFromPath(r)

	if err != nil {
		return err
	}

	err = h.service.Delete(id)

	if err != nil {
		return err
	}

	return nil
}
