package pgstore

import (
	"context"

	"github.com/christianferraz/goexpert/40-RocketSeat/4-Taskify/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGTaskStore struct {
	Queries *Queries
	Pool    *pgxpool.Pool
}

func NewPGTaskStore(pool *pgxpool.Pool) PGTaskStore {
	return PGTaskStore{
		Queries: New(pool),
		Pool:    pool,
	}
}

func (pgs *PGTaskStore) CreateTask(ctx context.Context, title string, description string, priority int32) (store.Task, error) {
	task, err := pgs.Queries.CreateTask(ctx, CreateTaskParams{
		Title:       title,
		Description: description,
		Priority:    priority,
	})
	if err != nil {
		return store.Task{}, err
	}
	return store.Task{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Priority:    task.Priority,
		CreatedAt:   task.CreatedAt.Time,
		UpdatedAt:   task.UpdatedAt.Time,
	}, nil
}
func (pgtaskstore *PGTaskStore) GetTaskById(ctx context.Context, id int32) (store.Task, error) {
	panic("not implemented") // TODO: Implement
}
func (pgtaskstore *PGTaskStore) ListTasks(ctx context.Context) ([]store.Task, error) {
	panic("not implemented") // TODO: Implement
}
func (pgtaskstore *PGTaskStore) UpdateTask(ctx context.Context, id int32, title string, description string, priority int32) (store.Task, error) {
	panic("not implemented") // TODO: Implement
}
func (pgtaskstore *PGTaskStore) DeleteTask(ctx context.Context, id int32) error {
	panic("not implemented") // TODO: Implement
}
