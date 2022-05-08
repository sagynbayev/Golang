package main

import (
	"context"
	lru "github.com/hashicorp/golang-lru"
	"hw9/internal/http"
	"hw9/internal/message_broker/kafka"
	"hw9/internal/store/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go CatchTermination(cancel)

	dbURL := "postgresql://postgres:postgres@localhost:5433/jobs"
	store := postgres.NewDB()
	if err := store.Connect(dbURL); err != nil {
		panic(err)
	}
	defer store.Close()

	cache, err := lru.New2Q(6)
	if err != nil {
		panic(err)
	}

	brokers := []string{"localhost:29092"}
	broker := kafka.NewBroker(brokers, cache, "peer1")
	if err := broker.Connect(ctx); err != nil {
		panic(err)
	}
	defer broker.Close()

	srv := http.NewServer(
		ctx,
		http.WithAddress(":8080"),
		http.WithStore(store),
		http.WithCache(cache),
		http.WithBroker(broker),
	)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}

func CatchTermination(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Print("[WARN] caught termination signal")
	cancel()
}