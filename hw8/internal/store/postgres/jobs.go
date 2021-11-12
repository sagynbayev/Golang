package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"hw8/internal/models"
	"hw8/internal/store"
)

func (db *DB) Jobs() store.JobsRepository {
	if db.jobs == nil {
		db.jobs = NewJobsRepository(db.conn)
	}
	return db.jobs
}

type JobsRepository struct {
	conn *sqlx.DB
}

func NewJobsRepository(conn *sqlx.DB) store.JobsRepository {
	return &JobsRepository{conn: conn}
}

func (c JobsRepository) Create(ctx context.Context, job *models.Job) error {
	_, err := c.conn.Exec("INSERT INTO jobs(name, category_id, price, description) VALUES ($1, $2, $3, $4)", job.Name, job.CategoryID, job.Price, job.Description)
	if err != nil {
		return err
	}

	return nil
}

func (c JobsRepository) All(ctx context.Context) ([]*models.Job, error) {
	jobs := make([]*models.Job, 0)
	if err := c.conn.Select(&jobs, "SELECT * FROM jobs"); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (c JobsRepository) ByID(ctx context.Context, id int) (*models.Job, error) {
	job := new(models.Job)
	if err := c.conn.Get(job, "SELECT * FROM categories WHERE id=$1", id); err != nil {
		return nil, err
	}

	return job, nil
}

func (c JobsRepository) Update(ctx context.Context, job *models.Job) error {
	_, err := c.conn.Exec("UPDATE jobs SET name = $1, price = $2, description = $3  WHERE id = $4", job.Name, job.Price, job.Description, job.ID)
	if err != nil {
		return err
	}

	return nil
}

func (c JobsRepository) Delete(ctx context.Context, id int) error {
	_, err := c.conn.Exec("DELETE FROM jobs WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}