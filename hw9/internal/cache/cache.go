package cache

import (
	"context"
	"hw9/internal/models"
)

type Cache interface {
	Close() error

	Categories() CategoriesCacheRepo
	Jobs() JobsCacheRepo
	DeleteAll(ctx context.Context) error
}

type CategoriesCacheRepo interface {
	Set(ctx context.Context, key string, value []*models.Category) error
	Get(ctx context.Context, key string) ([]*models.Category, error)
}

type JobsCacheRepo interface {
	Set(ctx context.Context, key string, value []*models.Job) error
	Get(ctx context.Context, key string) ([]*models.Job, error)
}