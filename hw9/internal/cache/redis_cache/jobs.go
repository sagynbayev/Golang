package redis_cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"hw9/internal/cache"
	"hw9/internal/models"
	"time"
)

func (rc RedisCache) Jobs() cache.JobsCacheRepo {
	if rc.jobs == nil {
		rc.jobs = newJobsRepo(rc.client, rc.expires)
	}

	return rc.jobs
}

type JobsRepo struct {
	client  *redis.Client
	expires time.Duration
}

func newJobsRepo(client *redis.Client, exp time.Duration) cache.JobsCacheRepo {
	return &JobsRepo{
		client:  client,
		expires: exp,
	}
}

func (c JobsRepo) Set(ctx context.Context, key string, value []*models.Job) error {
	jobBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err = c.client.Set(ctx, key, jobBytes, c.expires*time.Second).Err(); err != nil {
		return err
	}

	return nil
}

func (c JobsRepo) Get(ctx context.Context, key string) ([]*models.Job, error) {
	result, err := c.client.Get(ctx, key).Result()
	switch err {
	case nil:
		break
	case redis.Nil:
		return nil, nil
	default:
		return nil, err
	}

	jobs := make([]*models.Job, 0)
	if err = json.Unmarshal([]byte(result), &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}