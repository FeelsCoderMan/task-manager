package task

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/FeelsCoderMan/task-manager/internal/entity"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Task entity.Task

type Service interface {
	Get(id int) (Task, error)
	Create(task Task) (Task, error)
	Update(task Task) error
	Delete(id int) error
	GetLatest(count int) ([]Task, error)
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

func (s service) Create(task Task) (Task, error) {
	sqlStatement := `INSERT INTO tasks (name, created_at, updated_at, tags, priority, description, completed) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	row := s.db.QueryRow(sqlStatement, task.Name, task.CreatedAt, task.UpdatedAt, pq.Array(task.Tags), task.Priority, task.Description, task.Completed)
	err := row.Scan(&task.ID)

	if err != nil {
		errorMessage := fmt.Sprintf("Could not create the task")
		s.logger.error(errorMessage)
		return task, errors.New(errorMessage)
	}

	return task, nil
}

func (s service) Get(id int) (Task, error) {
	var task Task
	fmt.Printf("ID is %d\n", id)
	sqlStatement := `SELECT id, name FROM tasks WHERE id = $1`
	err := s.db.QueryRow(sqlStatement, id).Scan(&task.ID, &task.Name)

	if err != nil {
		errorMessage := fmt.Sprintf("Could not get the task %d", id)
		s.logger.error(errorMessage)
		return task, errors.New(errorMessage)
	}

	return task, nil
}

func (s service) GetLatest(count int) ([]Task, error) {
	result := make([]Task, 0)
	sqlStatement := `SELECT id, name, created_at, updated_at, tags, priority, description, completed FROM tasks ORDER BY created_at DESC LIMIT $1`
	rows, err := s.db.Query(sqlStatement, count)

	if err != nil {
		errorMessage := fmt.Sprintf("Could not retrieve latest tasks")
		s.logger.error(errorMessage)
		return result, errors.New(errorMessage)
	}

	defer rows.Close()

	for rows.Next() {
		var task Task
		var tagsBytes []byte
		err := rows.Scan(&task.ID, &task.Name, &task.CreatedAt, &task.UpdatedAt, &tagsBytes, &task.Priority, &task.Description, &task.Completed)

		if err != nil {
			errorMessage := fmt.Sprintf("Could not get the task from GetLatest")
			s.logger.error(errorMessage)
			return result, errors.New(errorMessage)
		}

		if tagsBytes != nil {
			task.Tags = parseTags(tagsBytes)
		}

		result = append(result, task)
	}

	return result, nil
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

func parseTags(data []byte) []string {
	var tags []string

	json.Unmarshal(data, &tags)
	return tags
}
