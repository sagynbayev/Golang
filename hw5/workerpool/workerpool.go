package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "processing job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		if j%2 == 0 {
			results <- j * 2
		}
		wg.Done()
	}
}

func main() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		go worker(i, &wg, jobs, results)
	}

	for j := 1; j <= 5; j++ {
		jobs <- j
		wg.Add(1)
	}
	close(jobs)

	wg.Wait()
}
