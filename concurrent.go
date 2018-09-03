package main

import (
	"sync"
	"time"
)

// JobFunc is job processing function
type JobFunc func(str string)

// thread function
func threadFunc(queue chan string, wg *sync.WaitGroup, f JobFunc) chan bool {
	quit := make(chan bool, 1)
	go func() {
		for {
			select {
			case str := <-queue:
				wg.Add(1)
				f(str)
				wg.Done()

			case <-quit:
				return
			}
		}
	}()
	return quit
}

// Concurrent is data parallel processing
func Concurrent(threadCount int, list []string, f JobFunc) {
	var wg sync.WaitGroup
	queue := make(chan string, jobs*10)

	quit := make([]chan bool, jobs)
	for i := 0; i < jobs; i++ {
		quit[i] = threadFunc(queue, &wg, f)
	}

	for _, str := range list {
		queue <- str
	}

	timer := time.Tick(10 * time.Millisecond)
	for range timer {
		if len(queue) == 0 {
			for _, q := range quit {
				q <- true
			}
			wg.Wait()
			break
		}
	}
}
