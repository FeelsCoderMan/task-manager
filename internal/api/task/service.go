package task

import (
	"database/sql"
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
		formattedErr := fmt.Errorf("Could not execute insertion for the task: %w", err)
		s.logger.error(formattedErr)
		return formattedErr
	}

	return nil
}

func (s service) Get(id int) (Task, error) {
	sqlStatement := `SELECT id, name FROM tasks WHERE id = $1`
	rows, err := s.db.Query(sqlStatement, id)
	var task Task

	if err != nil {
		formattedErr := fmt.Errorf("Could not get the task %d for reason: %w", id, err)
		s.logger.error(formattedErr)
		return task, formattedErr
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&task.ID, &task.Name); err != nil {
			return task, err
		}
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
		formattedErr := fmt.Errorf("Could not update the task %d for reason: %w", task.ID, err)
		s.logger.error(formattedErr)
		return formattedErr
	}

	return nil
}

func (s service) Delete(id int) error {
	sqlStatement := `DELETE FROM tasks
		WHERE id = $1
	`
	_, err := s.db.Exec(sqlStatement, id)

	if err != nil {
		formattedErr := fmt.Errorf("Could not update the task %d for reason: %w", id, err)
		s.logger.error(formattedErr)
		return formattedErr
	}

	return nil
}
