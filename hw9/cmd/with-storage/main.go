package main

import (
	"context"
	"hw9/internal/cache/redis_cache"
	"hw9/internal/http"
	"hw9/internal/store/postgres"
)

func main() {
	urlExample := "postgresql://postgres:postgres@localhost:5433/jobs"
	store := postgres.NewDB()
	if err := store.Connect(urlExample); err != nil {
		panic(err)
	}
	defer store.Close()

	cache := redis_cache.NewRedisCache("localhost:6379", 1, 3000)
	//defer cache.Close()

	srv := http.NewServer(
		context.Background(),
		http.WithAddress(":8080"),
		http.WithStore(store),
		http.WithCache(cache),
	)
	if err := srv.Run(); err != nil {
		panic(err)
	}

	srv.WaitForGracefulTermination()
}
