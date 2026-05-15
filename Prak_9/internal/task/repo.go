package task

import (
	"errors"
	"time"
)

var ErrTaskNotFound = errors.New("task not found")

type Repo struct {
	data map[int64]Task
}

func NewRepo() *Repo {
	return &Repo{
		data: map[int64]Task{
			1: {
				ID:          1,
				Title:       "Изучить Redis",
				Description: "Разобрать cache-aside",
				DueDate:     time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC),
			},
			2: {
				ID:          2,
				Title:       "Сделать ПЗ",
				Description: "Реализовать кэширование по id",
				DueDate:     time.Date(2026, 1, 21, 0, 0, 0, 0, time.UTC),
			},
		},
	}
}

func (r *Repo) GetAllTasks() ([]Task, error) {
	res := make([]Task, 0, len(r.data))
	for _, v := range r.data {
		res = append(res, v)
	}
	return res, nil
}

func (r *Repo) GetByID(id int64) (Task, error) {
	t, ok := r.data[id]
	if !ok {
		return Task{}, ErrTaskNotFound
	}
	return t, nil
}

func (r *Repo) Update(t Task) error {
	if _, ok := r.data[t.ID]; !ok {
		return ErrTaskNotFound
	}
	r.data[t.ID] = t
	return nil
}

func (r *Repo) Delete(id int64) error {
	if _, ok := r.data[id]; !ok {
		return ErrTaskNotFound
	}
	delete(r.data, id)
	return nil
}
