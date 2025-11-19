package main

import (
	"context"
	"database/sql"
	"time"
)

// Task — модель для сканирования результатов SELECT
type Task struct {
	ID        int
	Title     string
	Done      bool
	CreatedAt time.Time
}

type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo { return &Repo{DB: db} }

// CreateTask — параметризованный INSERT с возвратом id
func (r *Repo) CreateTask(ctx context.Context, title string) (int, error) {
	var id int
	const q = `INSERT INTO tasks (title) VALUES ($1) RETURNING id;`
	err := r.DB.QueryRowContext(ctx, q, title).Scan(&id)
	return id, err
}

// ListTasks — базовый SELECT всех задач (демо для занятия)
func (r *Repo) ListTasks(ctx context.Context) ([]Task, error) {
	const q = `SELECT id, title, done, created_at FROM tasks ORDER BY id;`
	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repo) ListDone(ctx context.Context, done bool) ([]Task, error) {
	const query = "SELECT id, title, done, created_at FROM tasks WHERE done = $1;"
	rows, err := r.DB.QueryContext(ctx, query, done)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Repo) FindByID(ctx context.Context, id int) (*Task, error) {
	const query = "SELECT id, title, done, created_at FROM tasks WHERE id = $1;"
	rows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out *Task
	if err := rows.Scan(&out.ID, &out.Title, &out.Done, &out.CreatedAt); err != nil {
		return nil, err
	}
	return out, rows.Err()
}

func (r *Repo) CreateMany(ctx context.Context, titles []string) error {
	// 1. Начинаем транзакцию
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Eсли функция завершится с ошибкой (до Commit),
	// транзакция будет автоматически отменена.
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	// Выполняем INSERT для каждой задачи внутри транзакции
	const q = `INSERT INTO tasks (title) VALUES ($1);`
	for _, title := range titles {
		_, err = tx.ExecContext(ctx, q, title)
		if err != nil {
			return err
		}
	}

	// Фиксируем изменения (если все успешно)
	err = tx.Commit()
	return err
}
