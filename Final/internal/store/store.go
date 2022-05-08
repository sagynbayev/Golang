package store

import (
	"context"
	"hw9/internal/models"
)

type Store interface {
	Connect(url string) error
	Close() error
	Jobs() JobsRepository
	Categories() CategoriesRepository
	Users() UsersRepository
}

type CategoriesRepository interface {
	Create(ctx context.Context, category *models.Category) error
	All(ctx context.Context, filter *models.CategoriesFilter) ([]*models.Category, error)
	ByID(ctx context.Context, id int) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id int) error
}

type JobsRepository interface {
	Create(ctx context.Context, job *models.Job) error
	All(ctx context.Context, filter *models.JobsFilter) ([]*models.Job, error)
	ByID(ctx context.Context, id int) (*models.Job, error)
	Update(ctx context.Context, job *models.Job) error
	Delete(ctx context.Context, id int) error
}
type UsersRepository interface {
	Create(ctx context.Context, user *models.User) error
	All(ctx context.Context, filter *models.UsersFilter) ([]*models.User, error)
	ByID(ctx context.Context, id int) (*models.User, error)
	Update(ctx context.Context, job *models.User) error
	Delete(ctx context.Context, id int) error
}
