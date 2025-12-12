package task

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/FeelsCoderMan/task-manager/internal/entity"
)

type Task entity.Task

type Service interface {
	Get(id int) (Task, error)
	Create(task Task) error
	Update(task Task) error
	Delete(id int) error
}

type service struct {
	db     *sql.DB
	logger TaskLogger
}

func NewService(db *sql.DB) Service {
	return service{
		db:     db,
		logger: NewTaskLogger(),
	}
}

func (s service) Create(task Task) error {
	sqlStatement := `INSERT INTO tasks (id, name) VALUES ($1, $2)`

	_, err := s.db.Exec(sqlStatement, task.ID, task.Name)

	if err != nil {
		errorMessage := fmt.Sprintf("Could not create a task %d", task.ID)
		s.logger.error(errorMessage)
		return errors.New(errorMessage)
	}

	return nil
}

func (s service) Get(id int) (Task, error) {
	var task Task
	sqlStatement := `SELECT id, name FROM tasks WHERE id = $1`
	err := s.db.QueryRow(sqlStatement, id).Scan(&task.ID, &task.Name)

	if err != nil {
		errorMessage := fmt.Sprintf("Could not get the task %d", task.ID)
		s.logger.error(errorMessage)
		return task, errors.New(errorMessage)
	}

	return task, nil
}

func (s service) Update(task Task) error {
	sqlStatement := `UPDATE tasks SET
		name = $1
		WHERE id = $2
	`
	_, err := s.db.Exec(sqlStatement, task.Name, task.ID)

	if err != nil {
		errorMessage := fmt.Sprintf("Could not update the task %d", task.ID)
		s.logger.error(errorMessage)
		return errors.New(errorMessage)
	}

	return nil
}

func (s service) Delete(id int) error {
	sqlStatement := `DELETE FROM tasks
		WHERE id = $1
	`
	_, err := s.db.Exec(sqlStatement, id)

	if err != nil {
		errorMessage := fmt.Sprintf("Could not delete the task %d", id)
		s.logger.error(errorMessage)
		return errors.New(errorMessage)
	}

	return nil
}
