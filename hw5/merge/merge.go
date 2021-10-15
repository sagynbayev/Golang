package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func testChannels(a ...int) <-chan int {
	ch := make(chan int)
	go func() {
		for _, i := range a {
			ch <- i
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		}
		close(ch)
	}()
	return ch
}
func merge(ch ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(ch))
	for _, i := range ch {
		go func(ch1 <-chan int) {
			for v := range ch1 {
				out <- v
			}
			wg.Done()
		}(i)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	a := testChannels(79, 3412, 125, 45, 9, 12, 7, 4)
	b := testChannels(146541, 4567, 7845, 65, 24, 712, 5)
	for i := range merge(a, b) {
		fmt.Println(i)
	}
}
