package inmemory

import (
	"context"
	"fmt"
	"hw6/internal/models"
	"hw6/internal/store"
	"sync"
)

type DB struct {
	data map[int]*models.Job

	mu *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		data: make(map[int]*models.Job),
		mu:   new(sync.RWMutex),
	}
}

func (db *DB) Create(ctx context.Context, job *models.Job) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[job.ID] = job
	return nil
}

func (db *DB) All(ctx context.Context) ([]*models.Job, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	jobs := make([]*models.Job, 0, len(db.data))
	for _, job := range db.data {
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (db *DB) ByID(ctx context.Context, id int) (*models.Job, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	job, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("No job with id %d", id)
	}

	return job, nil
}

func (db *DB) Update(ctx context.Context, job *models.Job) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[job.ID] = job
	return nil
}

func (db *DB) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}
