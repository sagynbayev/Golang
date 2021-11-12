package store

import (
	"context"
	"hw8/internal/models"
)

type Store interface {
	Connect(url string) error
	Close() error
	Jobs() JobsRepository
	Categories() CategoriesRepository
}

type CategoriesRepository interface {
	Create(ctx context.Context, category *models.Category) error
	All(ctx context.Context) ([]*models.Category, error)
	ByID(ctx context.Context, id int) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id int) error
}

type JobsRepository interface {
	Create(ctx context.Context, job *models.Job) error
	All(ctx context.Context) ([]*models.Job, error)
	ByID(ctx context.Context, id int) (*models.Job, error)
	Update(ctx context.Context, job *models.Job) error
	Delete(ctx context.Context, id int) error
}
